package client

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ─────────────────────────────────────────────────────────────────────────────
// doRequest & Example API
// ─────────────────────────────────────────────────────────────────────────────

func (c *thalassaCloudClient) Do(ctx context.Context, req *resty.Request, method httpMethod, url string) (*resty.Response, error) {
	// If we have a circuit breaker, wrap the request call in breaker.Execute.
	if c.breaker != nil {
		result, err := c.breaker.Execute(func() (interface{}, error) {
			return c.executeRequest(ctx, req, method, url)
		})
		if err != nil {
			return nil, fmt.Errorf("circuit breaker error for %s %s: %w", method, url, err)
		}
		return result.(*resty.Response), nil
	}

	// If no circuit breaker, just do the request directly.
	return c.executeRequest(ctx, req, method, url)
}

// executeRequest does the actual rate-limit & resty request call.
func (c *thalassaCloudClient) executeRequest(ctx context.Context, req *resty.Request, method httpMethod, url string) (*resty.Response, error) {
	// Enforce rate limiting if configured.
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait error: %w", err)
		}
	}
	req.SetContext(ctx)
	if c.organisationIdentity != nil {
		req.SetHeader("X-Organisation-Identity", *c.organisationIdentity)
	}
	if c.projectIdentity != nil {
		req.SetHeader("X-Project-Identity", *c.projectIdentity)
	}
	// All API calls are JSON.
	req.SetHeader("Accept", "application/json")

	var (
		resp *resty.Response
		err  error
	)
	switch method {
	case GET:
		resp, err = req.Get(url)
	case POST:
		resp, err = req.Post(url)
	case PUT:
		resp, err = req.Put(url)
	case PATCH:
		resp, err = req.Patch(url)
	case DELETE:
		resp, err = req.Delete(url)
	default:
		return nil, ErrUnsupportedHTTPMethod
	}

	if err != nil {
		return nil, fmt.Errorf("request to %s %s failed: %w", method, url, err)
	}
	return resp, nil
}

type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
	PATCH  httpMethod = "PATCH"
	DELETE httpMethod = "DELETE"
)

// ExampleAPI is a placeholder for a Thalassa Cloud endpoint call.
func (c *thalassaCloudClient) ExampleAPI(ctx context.Context, param string) (string, error) {
	req := c.resty.R().
		SetQueryParam("some-param", param).
		SetHeader("Accept", "application/json")

	resp, err := c.Do(ctx, req, GET, "/v1/example")
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode(), resp.String())
	}
	return resp.String(), nil
}

func (c *thalassaCloudClient) Check(resp *resty.Response) error {
	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return ErrNotFound
		}
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode(), resp.String())
	}
	return nil
}
