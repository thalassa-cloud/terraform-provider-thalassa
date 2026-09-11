package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceFederatedIdentityProvider(t *testing.T) {
	resource := ResourceFederatedIdentityProvider()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		schema := resource.Schema
		assert.True(t, schema["name"].Required)
		assert.True(t, schema["provider_issuer"].Required)
		assert.True(t, schema["provider_issuer"].ForceNew)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["object_version"].Computed)
		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
		assert.True(t, schema["provider_jwks_uri"].Optional)
		assert.True(t, schema["provider_jwks_uri"].Computed)
		assert.True(t, schema["local_jwks"].Optional)
		assert.True(t, schema["local_jwks"].Sensitive)
		assert.True(t, schema["status"].Optional)
		assert.True(t, schema["status"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
		assert.True(t, schema["organisation_id"].ForceNew)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
		assert.NotNil(t, resource.Importer)
	})
}

func TestDataSourceFederatedIdentityProvider(t *testing.T) {
	dataSource := DataSourceFederatedIdentityProvider()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		schema := dataSource.Schema
		assert.True(t, schema["id"].Optional)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["name"].Optional)
		assert.True(t, schema["provider_issuer"].Computed)
		assert.True(t, schema["status"].Computed)
		assert.True(t, schema["created_at"].Computed)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}

func TestResourceFederatedIdentity(t *testing.T) {
	resource := ResourceFederatedIdentity()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		schema := resource.Schema
		assert.True(t, schema["name"].Required)
		assert.True(t, schema["service_account_id"].Required)
		assert.True(t, schema["service_account_id"].ForceNew)
		assert.True(t, schema["provider_id"].Required)
		assert.True(t, schema["provider_id"].ForceNew)
		assert.True(t, schema["provider_subject"].Required)
		assert.True(t, schema["provider_subject"].ForceNew)
		assert.True(t, schema["trusted_audiences"].Required)
		assert.True(t, schema["audience_match_mode"].Required)
		assert.True(t, schema["allowed_scopes"].Required)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["last_used_at"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["object_version"].Computed)
		assert.True(t, schema["status"].Optional)
		assert.True(t, schema["status"].Computed)
		assert.True(t, schema["expires_at"].Optional)
		assert.True(t, schema["conditions"].Optional)
		assert.True(t, schema["organisation_id"].Optional)
		assert.True(t, schema["organisation_id"].ForceNew)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
		assert.NotNil(t, resource.Importer)
	})
}

func TestDataSourceFederatedIdentity(t *testing.T) {
	dataSource := DataSourceFederatedIdentity()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		schema := dataSource.Schema
		assert.True(t, schema["id"].Optional)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["name"].Optional)
		assert.True(t, schema["service_account_id"].Computed)
		assert.True(t, schema["provider_id"].Computed)
		assert.True(t, schema["allowed_scopes"].Computed)
		assert.True(t, schema["created_at"].Computed)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}

func TestParseLocalJWKS(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		jwks, err := parseLocalJWKS("")
		assert.NoError(t, err)
		assert.Nil(t, jwks)
	})

	t.Run("valid", func(t *testing.T) {
		jwks, err := parseLocalJWKS(`{"keys":[{"kty":"RSA","kid":"1"}]}`)
		assert.NoError(t, err)
		assert.NotNil(t, jwks)
		assert.Len(t, jwks.Keys, 1)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseLocalJWKS(`not-json`)
		assert.Error(t, err)
	})
}
