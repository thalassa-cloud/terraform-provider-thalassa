package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	oidcGrantTypeTokenExchange = "urn:ietf:params:oauth:grant-type:token-exchange"
	oidcTokenTypeJWT           = "urn:ietf:params:oauth:token-type:jwt"
)

// OIDCTokenExchangeConfig configures exchanging an external OIDC JWT (e.g. GitLab ID token)
// for a Thalassa Cloud API bearer token via POST application/x-www-form-urlencoded
// to TokenURL (typically {api}/oidc/token).
type OIDCTokenExchangeConfig struct {
	TokenURL         string
	SubjectToken     string
	SubjectTokenFile string
	OrganisationID   string
	ServiceAccountID string
	// AccessTokenLifetime is optional (e.g. "39600s"); sent as access_token_lifetime when non-empty.
	AccessTokenLifetime string
}

// WithAuthOIDCTokenExchange exchanges a subject JWT for an API access token before each request
// when the cached token is missing or expired. The resulting bearer token is used like AuthOIDC.
//
// Provide the subject JWT via SubjectToken, or SubjectTokenFile (contents read at each exchange;
// use for mounted tokens in Kubernetes so rotation is picked up on refresh).
func WithAuthOIDCTokenExchange(cfg OIDCTokenExchangeConfig) Option {
	return func(c *thalassaCloudClient) error {
		hasSubject := strings.TrimSpace(cfg.SubjectToken) != ""
		hasFile := strings.TrimSpace(cfg.SubjectTokenFile) != ""
		if strings.TrimSpace(cfg.TokenURL) == "" ||
			strings.TrimSpace(cfg.OrganisationID) == "" ||
			strings.TrimSpace(cfg.ServiceAccountID) == "" ||
			(!hasSubject && !hasFile) {
			return ErrOIDCTokenExchangeConfig
		}
		c.authType = AuthOIDCTokenExchange
		cfgCopy := cfg
		c.oidcTokenExchange = &cfgCopy
		return nil
	}
}

// WithAuthOIDC uses client credentials for service-to-service flows.
func WithAuthOIDC(clientID, clientSecret, tokenURL string, scopes ...string) Option {
	return withAuthOIDC(clientID, clientSecret, tokenURL, false, scopes...)
}

func WithAuthOIDCInsecure(clientID, clientSecret, tokenURL string, allowInsecure bool, scopes ...string) Option {
	return withAuthOIDC(clientID, clientSecret, tokenURL, allowInsecure, scopes...)
}

func withAuthOIDC(clientID, clientSecret, tokenURL string, insecure bool, scopes ...string) Option {
	return func(c *thalassaCloudClient) error {
		c.authType = AuthOIDC
		c.oidcConfig = &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     tokenURL,
			Scopes:       scopes,
		}
		c.allowInsecureOIDC = insecure
		return nil
	}
}

func WithToken(token string) Option {
	return func(c *thalassaCloudClient) error {
		c.authType = AuthToken
		c.oidcToken = &oauth2.Token{AccessToken: token}
		return nil
	}
}

func WithAuthPersonalToken(token string) Option {
	return func(c *thalassaCloudClient) error {
		c.authType = AuthPersonalAccessToken
		c.personalToken = token
		return nil
	}
}

func WithAuthBasic(username, password string) Option {
	return func(c *thalassaCloudClient) error {
		c.authType = AuthBasic
		c.basicUsername = username
		c.basicPassword = password
		return nil
	}
}

func WithAuthNone() Option {
	return func(c *thalassaCloudClient) error {
		c.authType = AuthNone
		return nil
	}
}

func (c *thalassaCloudClient) configureAuth() error {
	switch c.authType {
	case AuthToken:
		if c.oidcToken == nil {
			return ErrMissingToken
		}
		c.resty.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
			if c.oidcToken == nil || !c.oidcToken.Valid() {
				return fmt.Errorf("token is not valid")
			}
			req.SetAuthToken(c.oidcToken.AccessToken)
			return nil
		})
	case AuthOIDC:
		if c.oidcConfig == nil {
			return ErrMissingOIDCConfig
		}
		// For each request, ensure token is valid or refresh it.
		c.resty.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
			if c.oidcToken == nil || !c.oidcToken.Valid() {
				ctx := req.Context()
				if c.allowInsecureOIDC {
					ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
						Transport: &http.Transport{
							TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						},
					})
				}
				tok, err := c.oidcConfig.Token(ctx)
				if err != nil {
					return fmt.Errorf("failed to fetch OIDC token: %w", err)
				}
				c.oidcToken = tok
			}
			req.SetAuthToken(c.oidcToken.AccessToken)
			return nil
		})

	case AuthOIDCTokenExchange:
		if c.oidcTokenExchange == nil {
			return ErrOIDCTokenExchangeConfig
		}
		c.resty.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
			if err := c.ensureOIDCTokenExchange(req.Context()); err != nil {
				return err
			}
			req.SetAuthToken(c.oidcToken.AccessToken)
			return nil
		})

	case AuthPersonalAccessToken:
		if c.personalToken == "" {
			return ErrEmptyPersonalToken
		}
		// c.resty.SetAuthToken(c.personalToken)
		c.resty.SetHeader("Authorization", "Token "+c.personalToken)
	case AuthBasic:
		if c.basicUsername == "" || c.basicPassword == "" {
			return ErrMissingBasicCredentials
		}
		c.resty.SetBasicAuth(c.basicUsername, c.basicPassword)

	case AuthCustom:
		// Let the user attach custom OnBeforeRequest callbacks.

	case AuthNone:
		// No authentication.

	default:
		// Should not occur. No special action.
	}
	return nil
}

func (c *thalassaCloudClient) ensureOIDCTokenExchange(ctx context.Context) error {
	c.oidcTokenExchangeMu.Lock()
	defer c.oidcTokenExchangeMu.Unlock()
	if c.oidcToken != nil && c.oidcToken.Valid() {
		return nil
	}
	tok, err := c.fetchOIDCTokenExchange(ctx)
	if err != nil {
		return err
	}
	c.oidcToken = tok
	return nil
}

func (c *thalassaCloudClient) tokenExchangeHTTPClient() *http.Client {
	if c.insecure {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // dev-only, matches WithInsecure
			},
		}
	}
	return http.DefaultClient
}

func resolveOIDCSubjectToken(cfg *OIDCTokenExchangeConfig) (string, error) {
	if fp := strings.TrimSpace(cfg.SubjectTokenFile); fp != "" {
		b, err := os.ReadFile(fp)
		if err != nil {
			return "", fmt.Errorf("OIDC token exchange: read subject token file: %w", err)
		}
		t := strings.TrimSpace(string(b))
		if t == "" {
			return "", fmt.Errorf("OIDC token exchange: subject token file %q is empty", fp)
		}
		return t, nil
	}
	if t := strings.TrimSpace(cfg.SubjectToken); t != "" {
		return t, nil
	}
	return "", fmt.Errorf("OIDC token exchange: no subject token configured")
}

func (c *thalassaCloudClient) fetchOIDCTokenExchange(ctx context.Context) (*oauth2.Token, error) {
	cfg := c.oidcTokenExchange
	subjectToken, err := resolveOIDCSubjectToken(cfg)
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Set("grant_type", oidcGrantTypeTokenExchange)
	form.Set("subject_token", subjectToken)
	form.Set("subject_token_type", oidcTokenTypeJWT)
	form.Set("organisation_id", cfg.OrganisationID)
	form.Set("service_account_id", cfg.ServiceAccountID)
	if strings.TrimSpace(cfg.AccessTokenLifetime) != "" {
		form.Set("access_token_lifetime", cfg.AccessTokenLifetime)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("OIDC token exchange: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if ua := strings.TrimSpace(c.userAgent); ua != "" {
		httpReq.Header.Set("User-Agent", ua)
	}

	resp, err := c.tokenExchangeHTTPClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("OIDC token exchange: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("OIDC token exchange: read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OIDC token exchange: %s (body: %s)", resp.Status, strings.TrimSpace(string(body)))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("OIDC token exchange: decode response: %w", err)
	}
	if tr.AccessToken == "" {
		return nil, fmt.Errorf("OIDC token exchange: empty access_token in response")
	}

	tok := &oauth2.Token{AccessToken: tr.AccessToken, TokenType: tr.TokenType}
	switch {
	case tr.ExpiresIn > 0:
		tok.Expiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	default:
		// Without expires_in the token is treated as valid until process restart; prefer the API to return expires_in.
		tok.Expiry = time.Now().Add(time.Hour)
	}
	return tok, nil
}
