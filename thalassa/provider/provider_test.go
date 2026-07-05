package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func TestProviderConfigurePreservesAccessTokenForGetClient(t *testing.T) {
	t.Parallel()

	p := thalassa.Provider()
	rd := schema.TestResourceDataRaw(t, p.Schema, map[string]any{
		"access_token":    "test-access-token",
		"api":             "https://api.thalassa.cloud",
		"organisation_id": "org-test",
	})

	configured, diags := provider.ProviderConfigure(context.Background(), rd)
	assert.Empty(t, diags)

	resourceData := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"organisation_id": {Type: schema.TypeString, Optional: true},
	}, map[string]any{})

	client, err := provider.GetClient(provider.GetProvider(configured), resourceData)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestProviderConfigurePreservesAllowInsecureOIDC(t *testing.T) {
	t.Parallel()

	p := thalassa.Provider()
	rd := schema.TestResourceDataRaw(t, p.Schema, map[string]any{
		"client_id":           "client-id",
		"client_secret":       "client-secret",
		"allow_insecure_oidc": true,
		"api":                 "https://api.thalassa.cloud",
		"organisation_id":     "org-test",
	})

	configured, diags := provider.ProviderConfigure(context.Background(), rd)
	assert.Empty(t, diags)

	resourceData := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"organisation_id": {Type: schema.TypeString, Optional: true},
	}, map[string]any{})

	client, err := provider.GetClient(provider.GetProvider(configured), resourceData)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
