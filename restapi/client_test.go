package restapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		c := NewClient()

		if got, want := c.httpClient.Timeout, defaultTimeout; got != want {
			t.Fatalf("c.httpClient.Timeout = %v, want %v", got, want)
		}

		if got, want := c.baseURL.Scheme, defaultScheme; got != want {
			t.Fatalf("c.baseURL.Scheme = %q, want %q", got, want)
		}

		if got, want := c.baseURL.Host, defaultHost; got != want {
			t.Fatalf("c.baseURL.Host = %q, want %q", got, want)
		}

		if got, want := c.userAgent, defaultUserAgent; got != want {
			t.Fatalf("c.userAgent = %q, want %q", got, want)
		}

		if got, want := c.apiKey, ""; got != want {
			t.Fatalf("c.apiKey = %q, want %q", got, want)
		}

		if got, want := c.secret, ""; got != want {
			t.Fatalf("c.secret = %q, want %q", got, want)
		}
	})

	t.Run("HTTPClient", func(t *testing.T) {
		timeout := 123 * time.Second

		c := NewClient(HTTPClient(&http.Client{Timeout: timeout}))

		if got, want := c.httpClient.Timeout, timeout; got != want {
			t.Fatalf("c.httpClient.Timeout = %v, want %v", got, want)
		}
	})

	t.Run("BaseURL", func(t *testing.T) {
		rawurl := "http://example.com"

		c := NewClient(BaseURL(rawurl))

		if got, want := c.baseURL.String(), rawurl; got != want {
			t.Fatalf("c.baseURL.String() = %q, want %q", got, want)
		}
	})

	t.Run("UserAgent", func(t *testing.T) {
		ua := "user-agent-test"

		c := NewClient(UserAgent(ua))

		if got, want := c.userAgent, ua; got != want {
			t.Fatalf("c.userAgent = %q, want %q", got, want)
		}
	})

	t.Run("Credentials", func(t *testing.T) {
		apiKey := "test apiKey"
		secret := "test secret"

		c := NewClient(Credentials(apiKey, secret))

		if got, want := c.apiKey, apiKey; got != want {
			t.Fatalf("c.apiKey = %q, want %q", got, want)
		}

		if got, want := c.secret, secret; got != want {
			t.Fatalf("c.apiKey = %q, want %q", got, want)
		}
	})
}

func TestGetRequest(t *testing.T) {
	c := testClient()

	for _, tt := range []struct {
		path   string
		query  url.Values
		rawurl string
		err    error
	}{
		{"/Foo", url.Values{"bar": {"hey"}}, "http://example.com/Foo?bar=hey", nil},
		{"/Bar", url.Values{"baz": {"123"}}, "http://example.com/Bar?baz=123", nil},
		{"::/foo", url.Values{"qux": {"456"}}, "", errors.New("parse ::/foo?qux=456: missing protocol scheme")},
	} {
		t.Run(tt.path, func(t *testing.T) {
			req, err := c.getRequest(context.Background(), tt.path, tt.query)
			if err != nil {
				if tt.err == nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got, want := err.Error(), tt.err.Error(); got != want {
					t.Fatalf("err.Error() = %q, want %q", got, want)
				}

				return
			}

			if got, want := req.Method, http.MethodGet; got != want {
				t.Fatalf("req.Method = %q, want %q", got, want)
			}

			if got, want := req.URL.Path, tt.path; got != want {
				t.Fatalf("req.URL.Path = %q, want %q", got, want)
			}

			if got, want := req.UserAgent(), defaultUserAgent; got != want {
				t.Fatalf("req.UserAgent() = %q, want %q", got, want)
			}

			if got, want := req.URL.String(), tt.rawurl; got != want {
				t.Fatalf("req.URL.String() = %q, want %q", got, want)
			}
		})
	}
}

const (
	testBaseURL = "http://example.com"
)

func testClient(options ...func(*Client)) *Client {
	return NewClient(append(options, BaseURL(testBaseURL))...)
}

func testServerAndClient(hf http.HandlerFunc, options ...func(*Client)) (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(hf))

	return ts, NewClient(append(options, BaseURL(ts.URL))...)
}
