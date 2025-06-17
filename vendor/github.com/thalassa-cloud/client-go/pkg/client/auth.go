package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

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
