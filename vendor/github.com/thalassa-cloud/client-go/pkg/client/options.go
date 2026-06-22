package client

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

func WithBaseURL(url string) Option {
	return func(c *thalassaCloudClient) error {
		c.resty.SetBaseURL(url)
		c.baseURL = url
		return nil
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *thalassaCloudClient) error {
		c.resty.SetTimeout(d)
		return nil
	}
}

func WithOrganisation(organisationIdentity string) Option {
	return func(c *thalassaCloudClient) error {
		c.organisationIdentity = &organisationIdentity
		return nil
	}
}

func WithProject(projectIdentity string) Option {
	return func(c *thalassaCloudClient) error {
		c.projectIdentity = &projectIdentity
		return nil
	}
}

func WithRetries(count int, waitTime, maxWaitTime time.Duration) Option {
	return func(c *thalassaCloudClient) error {
		if count > 0 {
			c.resty.
				SetRetryCount(count).
				SetRetryWaitTime(waitTime).
				SetRetryMaxWaitTime(maxWaitTime)
		}
		return nil
	}
}

func WithAuthCustom() Option {
	return func(c *thalassaCloudClient) error {
		c.authType = AuthCustom
		return nil
	}
}

func WithRateLimit(rps float64, burst int) Option {
	return func(c *thalassaCloudClient) error {
		c.limiter = rate.NewLimiter(rate.Limit(rps), burst)
		return nil
	}
}

// WithCircuitBreaker configures a circuit breaker using sony/gobreaker.
func WithCircuitBreaker(name string, st gobreaker.Settings) Option {
	return func(c *thalassaCloudClient) error {
		st.Name = name
		c.breaker = gobreaker.NewCircuitBreaker(st)
		return nil
	}
}

func WithMiddleware(mw func(*resty.Client, *resty.Request) error) Option {
	return func(c *thalassaCloudClient) error {
		c.resty.OnBeforeRequest(mw)
		return nil
	}
}

// AddMiddleware is a convenience to add further request hooks after creation.
func (c *thalassaCloudClient) AddMiddleware(mw func(*resty.Client, *resty.Request) error) {
	c.resty.OnBeforeRequest(mw)
}

func WithUserAgent(ua string) Option {
	return func(c *thalassaCloudClient) error {
		c.userAgent = ua
		return nil
	}
}

// WithInsecure disables SSL certificate verification.
// This should only be used for development/testing purposes.
func WithInsecure() Option {
	return func(c *thalassaCloudClient) error {
		c.insecure = true
		c.resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		return nil
	}
}
