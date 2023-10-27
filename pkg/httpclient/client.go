package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Middleware func(r *http.Request)

type HttpClient struct {
	baseURL string
	*http.Client
	middlewares []Middleware
}

func NewHttpClient(client *http.Client, baseURL string, middlewares ...Middleware) *HttpClient {
	return &HttpClient{Client: client, baseURL: baseURL, middlewares: middlewares}
}

func (c *HttpClient) AddRequestMiddleware(f func(r *http.Request)) {
	c.middlewares = append(c.middlewares, f)
}

func (c *HttpClient) MakeRequest(ctx context.Context, method string, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed create new request, %w", err)
	}
	for _, m := range c.middlewares {
		m(req)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.Do(req)
}
