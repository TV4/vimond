/*

Package restapi contains a client for the Vimond REST API

*/
package restapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Errors
var (
	ErrInvalidAssetID = errors.New("vimond/restapi: invalid asset id")
	ErrNotFound       = errors.New("vimond/restapi: not found")
	ErrUnknown        = errors.New("vimond/restapi: unknown")
)

const (
	defaultScheme    = "https"
	defaultHost      = "restapi-vimond-prod.b17g.net"
	defaultUserAgent = "vimond/restapi/client.go"
	defaultTimeout   = 20 * time.Second
)

// Client for the Vimond REST API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	apiKey     string
	secret     string
	userAgent  string
}

// NewClient creates a new Vimond REST API Client
func NewClient(options ...func(*Client)) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: &url.URL{
			Scheme: defaultScheme,
			Host:   defaultHost,
		},
		userAgent: defaultUserAgent,
	}

	for _, f := range options {
		f(c)
	}

	return c
}

// HTTPClient changes the *client HTTP client to the provided *http.Client
func HTTPClient(hc *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// BaseURL changes the *client base URL based on the provided rawurl
func BaseURL(rawurl string) func(*Client) {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// UserAgent changes the User-Agent used in requests sent by the *client
func UserAgent(ua string) func(*Client) {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// Credentials changes the apiKey and secret used to sign API requests
func Credentials(apiKey, secret string) func(*Client) {
	return func(c *Client) {
		c.apiKey = apiKey
		c.secret = secret
	}
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil, c.setAuthorizationHeader())
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) getJSON(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil, c.setAuthorizationHeader(), func(req *http.Request) {
		req.Header.Add("Accept", "application/json; v=3; charset=utf-8")
		req.Header.Add("Content-Type", "application/json; v=3; charset=utf-8")
	})
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) post(ctx context.Context, path string, query url.Values, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, query, body, c.setAuthorizationHeader())
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) postJSON(ctx context.Context, path string, query url.Values, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, query, body, c.setAuthorizationHeader(), func(req *http.Request) {
		req.Header.Add("Accept", "application/json; v=3; charset=utf-8")
		req.Header.Add("Content-Type", "application/json; v=3; charset=utf-8")
	})
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) put(ctx context.Context, path string, query url.Values, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, query, body, c.setAuthorizationHeader())
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) putJSON(ctx context.Context, path string, query url.Values, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, query, body, c.setAuthorizationHeader(), func(req *http.Request) {
		req.Header.Add("Accept", "application/json; v=3; charset=utf-8")
		req.Header.Add("Content-Type", "application/json; v=3; charset=utf-8")
	})
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body io.Reader, options ...func(*http.Request)) (*http.Request, error) {
	rawurl := path

	if len(query) > 0 {
		rawurl += "?" + query.Encode()
	}

	rel, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.baseURL.ResolveReference(rel).String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	for _, option := range options {
		option(req)
	}

	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) setAuthorizationHeader() func(*http.Request) {
	return func(req *http.Request) {
		if c.apiKey != "" && c.secret != "" {
			req.Header = authorizationHeader(req.Method, req.URL.Path, time.Now(), c.apiKey, c.secret)
		}
	}
}
