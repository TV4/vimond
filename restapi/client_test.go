package restapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
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

func TestNewRequest(t *testing.T) {
	c := testClient()

	for _, tt := range []struct {
		method string
		path   string
		query  url.Values
		body   io.Reader

		rawurl    string
		bodyBytes []byte
		err       error
	}{
		{http.MethodGet, "/Foo", url.Values{"bar": {"hey"}}, nil, "http://example.com/Foo?bar=hey", nil, nil},
		{http.MethodGet, "/Bar", url.Values{"baz": {"123"}}, nil, "http://example.com/Bar?baz=123", nil, nil},
		{http.MethodGet, "::/foo", url.Values{"qux": {"456"}}, nil, "", nil, errors.New("parse \"::/foo?qux=456\": missing protocol scheme")},
		{http.MethodPost, "/Foo", url.Values{"bar": {"hey"}}, bytes.NewReader([]byte("foo-body")), "http://example.com/Foo?bar=hey", []byte("foo-body"), nil},
	} {
		t.Run(tt.path, func(t *testing.T) {
			req, err := c.newRequest(context.Background(), tt.method, tt.path, tt.query, tt.body)
			if err != nil {
				if tt.err == nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got, want := err.Error(), tt.err.Error(); got != want {
					t.Fatalf("err.Error() = %q, want %q", got, want)
				}

				return
			}

			if got, want := req.Method, tt.method; got != want {
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

			switch {
			case req.Body == nil && tt.bodyBytes == nil:
				// All good
			case req.Body == nil && tt.bodyBytes != nil:
				t.Fatalf("req.Body is nil, want %q", tt.bodyBytes)
			case req.Body != nil && tt.bodyBytes == nil:
				gotBody, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				t.Fatalf("req.Body = %q, want nil", gotBody)
			case req.Body != nil && tt.bodyBytes != nil:
				gotBody, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got, want := gotBody, tt.bodyBytes; !bytes.Equal(got, want) {
					t.Fatalf("req.Body = %q, want %q", got, want)
				}
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
