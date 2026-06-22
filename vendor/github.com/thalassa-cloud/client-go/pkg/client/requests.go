package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func parseHTTPMethod(m string) (httpMethod, error) {
	switch strings.ToUpper(strings.TrimSpace(m)) {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "PATCH":
		return PATCH, nil
	case "DELETE":
		return DELETE, nil
	default:
		return "", ErrUnsupportedHTTPMethod
	}
}

// RawRequest performs an HTTP request using the client's base URL, authentication,
// and configuration. path is relative to the client's base URL. body is optional;
// when non-nil it is sent with Content-Type: application/json.
func (c *thalassaCloudClient) RawRequest(ctx context.Context, method, path string, body []byte) (*resty.Response, error) {
	m, err := parseHTTPMethod(method)
	if err != nil {
		return nil, err
	}
	req := c.R()
	if len(body) > 0 {
		req.SetBody(body).SetHeader("Content-Type", "application/json")
	}
	return c.Do(ctx, req, m, path)
}

func (c *thalassaCloudClient) Check(resp *resty.Response) error {
	if resp.IsError() {
		switch resp.StatusCode() {
		case http.StatusNotFound:
			var errorMessage ServerErrorMessage
			if err := json.Unmarshal(resp.Body(), &errorMessage); err != nil {
				return ErrNotFound
			}
			return errors.Join(ErrNotFound, errors.New(errorMessage.Message))
		case http.StatusBadRequest:
			var errorMessage ServerErrorMessage
			if err := json.Unmarshal(resp.Body(), &errorMessage); err != nil {
				return ErrBadRequest
			}
			return errors.Join(ErrBadRequest, errors.New(errorMessage.Message))
		default:
			return fmt.Errorf("server returned status %d: %s", resp.StatusCode(), resp.String())
		}
	}
	return nil
}

type ServerErrorMessage struct {
	Message string `json:"message"`
}
