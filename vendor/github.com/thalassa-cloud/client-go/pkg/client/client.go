package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/sony/gobreaker"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
)

const (
	DefaultUserAgent = "thalassa-cloud-client-go (https://github.com/thalassa-cloud/client-go)"
)

// Option is a function that modifies the Client.
type Option func(*thalassaCloudClient) error

var (
	ErrMissingBaseURL          = errors.New("missing base URL; use WithBaseURL(...)")
	ErrMissingOIDCConfig       = errors.New("OIDC configuration is missing")
	ErrEmptyPersonalToken      = errors.New("personal access token cannot be empty")
	ErrMissingBasicCredentials = errors.New("basic auth requires username/password")
	ErrUnsupportedHTTPMethod   = errors.New("unsupported HTTP method")
	ErrNotFound                = errors.New("not found")
)

type AuthenticationType int

const (
	AuthNone AuthenticationType = iota
	AuthOIDC
	AuthPersonalAccessToken
	AuthBasic
	AuthCustom
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

type Client interface {
	Do(ctx context.Context, req *resty.Request, method httpMethod, url string) (*resty.Response, error)
	Check(resp *resty.Response) error

	R() *resty.Request

	WithOptions(opts ...Option) Client

	// GetOrganisationIdentity returns the organisation identity for the client, if set
	GetOrganisationIdentity() string

	// SetOrganisation sets the organisation for the client
	SetOrganisation(organisation string)

	// GetAuthToken returns the authentication token for the client, if set
	GetAuthToken() string

	// GetBaseURL returns the base URL for the client
	GetBaseURL() string

	// DialWebsocket creates a websocket connection to the specified URL
	DialWebsocket(ctx context.Context, wsURL string) (*websocket.Conn, error)
}

// NewClient applies all options, configures authentication, and returns the client.
func NewClient(opts ...Option) (Client, error) {
	c := &thalassaCloudClient{
		resty:     resty.New(),
		userAgent: DefaultUserAgent,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.resty.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	// Configure built-in authentication once we have all fields set.
	if err := c.configureAuth(); err != nil {
		return nil, err
	}

	return c, nil
}

type thalassaCloudClient struct {
	// Underlying resty client.
	resty *resty.Client

	baseURL   string
	userAgent string

	organisationIdentity *string
	projectIdentity      *string

	// Authentication fields.
	authType AuthenticationType

	// OIDC (client credentials).
	oidcConfig        *clientcredentials.Config
	oidcToken         *oauth2.Token // cached token
	allowInsecureOIDC bool

	// Personal Access Token.
	personalToken string

	// Basic Auth.
	basicUsername string
	basicPassword string

	// Rate limiting.
	limiter *rate.Limiter

	// Optional circuit breaker
	breaker *gobreaker.CircuitBreaker

	insecure bool
}

func (c *thalassaCloudClient) WithOptions(opts ...Option) Client {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *thalassaCloudClient) R() *resty.Request {
	return c.resty.R().SetHeader("User-Agent", c.userAgent)
}

func (c *thalassaCloudClient) GetOrganisationIdentity() string {
	if c.organisationIdentity != nil {
		return *c.organisationIdentity
	}
	return ""
}

func (c *thalassaCloudClient) SetOrganisation(organisation string) {
	c.organisationIdentity = &organisation
}

func (c *thalassaCloudClient) GetAuthToken() string {
	switch c.authType {
	case AuthOIDC:
		if c.oidcToken != nil && c.oidcToken.Valid() {
			return c.oidcToken.AccessToken
		}
	case AuthPersonalAccessToken:
		return c.personalToken
	}
	return ""
}

// DialWebsocket creates a websocket connection to the specified URL, with authentication
// and organization headers from the client.
func (c *thalassaCloudClient) DialWebsocket(ctx context.Context, wsURL string) (*websocket.Conn, error) {

	wsUrlWithToken := wsURL + "?token=" + c.GetAuthToken()

	// Parse the WebSocket URL
	parsedURL, err := url.Parse(wsUrlWithToken)
	if err != nil {
		return nil, fmt.Errorf("invalid websocket URL: %w", err)
	}

	// Create dialer with any needed options
	dialer := websocket.DefaultDialer

	// Prepare headers
	header := http.Header{}

	// Apply authentication
	if token := c.GetAuthToken(); token != "" {
		header.Add("Authorization", "Token "+token)
	}

	// Apply organization identity
	if orgIdentity := c.GetOrganisationIdentity(); orgIdentity != "" {
		header.Add("X-Organisation-Identity", orgIdentity)
	}

	// Apply project identity if available
	if c.projectIdentity != nil {
		header.Add("X-Project-Identity", *c.projectIdentity)
	}

	// Connect to WebSocket
	conn, _, err := dialer.DialContext(ctx, parsedURL.String(), header)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %w", err)
	}

	return conn, nil
}

func (c *thalassaCloudClient) GetBaseURL() string {
	return c.baseURL
}
