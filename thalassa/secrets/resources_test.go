package secrets

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	tckms "github.com/thalassa-cloud/client-go/kms"
	tcsecrets "github.com/thalassa-cloud/client-go/secrets"
)

func TestResourceSecret(t *testing.T) {
	resource := ResourceSecret()
	schema := resource.Schema

	assert.True(t, schema["region"].Required)
	assert.True(t, schema["path"].Required)
	assert.True(t, schema["kms_key_id"].Required)
	assert.True(t, schema["secret_string"].Sensitive)
	assert.NotNil(t, resource.Importer)
}

func TestResourceSecretVersion(t *testing.T) {
	resource := ResourceSecretVersion()
	assert.True(t, resource.Schema["secret_string"].Sensitive)
	assert.NotNil(t, resource.Schema["generate_secret"].Elem.(*schema.Resource).Schema["byte_length"])
	assert.NotNil(t, resource.DeleteContext)
}

func TestParseSecretID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantRegion string
		wantPath   string
		wantErr    bool
	}{
		{
			name:       "valid",
			id:         "nl-01/app/prod/db/password",
			wantRegion: "nl-01",
			wantPath:   "/app/prod/db/password",
		},
		{name: "invalid", id: "nopath", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, path, err := parseSecretID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRegion, region)
			assert.Equal(t, tt.wantPath, path)
		})
	}
}

func TestParseSecretVersionID(t *testing.T) {
	region, path, version, err := parseSecretVersionID("nl-01/app/prod/db/password/3")
	assert.NoError(t, err)
	assert.Equal(t, "nl-01", region)
	assert.Equal(t, "/app/prod/db/password", path)
	assert.Equal(t, 3, version)
}

func TestValidateSecretPath(t *testing.T) {
	_, errs := validateSecretPath("/app/prod", "path")
	assert.Empty(t, errs)

	_, errs = validateSecretPath("app/prod", "path")
	assert.NotEmpty(t, errs)
}

func TestSetSecretStateKmsKey(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceSecret().Schema, map[string]any{})
	secret := &tcsecrets.Secret{
		Path:   "/app/prod/db/password",
		KmsKey: &tckms.KmsKey{Identity: "kms-abc123"},
	}

	err := setSecretState(d, secret, "nl-01")
	assert.NoError(t, err)
	assert.Equal(t, "kms-abc123", d.Get("kms_key_id"))
}
