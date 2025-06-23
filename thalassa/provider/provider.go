package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func GetProvider(m interface{}) ConfiguredProvider {
	p, ok := m.(ConfiguredProvider)
	if !ok {
		panic("invalid configured provider")
	}
	return p
}

type ConfiguredProvider struct {
	Client       thalassa.Client
	Organisation string
	token        string
	apiEndpoint  string
	clientID     string
	clientSecret string
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	apiEndpoint := d.Get("api").(string)
	organisation := d.Get("organisation_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	opts := []client.Option{
		client.WithBaseURL(apiEndpoint),
		client.WithOrganisation(organisation),
		client.WithUserAgent("thalassa-cloud/terraform-provider-thalassa"),
	}

	hasAuth := false
	if token != "" {
		opts = append(opts, client.WithAuthPersonalToken(token))
		hasAuth = true
	}

	if clientID != "" && clientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(clientID, clientSecret, fmt.Sprintf("%s/oidc/token", apiEndpoint)))
		hasAuth = true
	}

	if !hasAuth {
		return nil, diag.FromErr(errors.New("no authentication method provided"))
	}

	internalClient, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return ConfiguredProvider{
		Client:       internalClient,
		Organisation: organisation,
		token:        token,
		apiEndpoint:  apiEndpoint,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

func GetClient(provider ConfiguredProvider, d *schema.ResourceData) (thalassa.Client, error) {
	organisation, err := getOrganisation(provider, d)
	if err != nil {
		return nil, err
	}

	opts := []client.Option{
		client.WithBaseURL(provider.apiEndpoint),
		client.WithOrganisation(organisation),
		client.WithUserAgent("thalassa-cloud/terraform-provider-thalassa"),
	}

	hasAuth := false
	if provider.token != "" {
		opts = append(opts, client.WithAuthPersonalToken(provider.token))
		hasAuth = true
	}

	if provider.clientID != "" && provider.clientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(provider.clientID, provider.clientSecret, fmt.Sprintf("%s/oidc/token", provider.apiEndpoint)))
		hasAuth = true
	}

	if !hasAuth {
		return nil, errors.New("no authentication method provided")
	}

	client, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getOrganisation(provider ConfiguredProvider, d *schema.ResourceData) (string, error) {
	organisation := provider.Organisation
	orgFromState := d.Get("organisation_id")
	if orgFromState != nil {
		if o, ok := orgFromState.(string); ok {
			organisation = o
		}
	}
	if organisation != "" {
		return organisation, nil
	}
	if provider.Organisation != "" {
		return provider.Organisation, nil
	}
	return "", errors.New("organisation is not set")
}
