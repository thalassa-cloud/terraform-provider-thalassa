package secrets

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	tckms "github.com/thalassa-cloud/client-go/kms"
	tcsecrets "github.com/thalassa-cloud/client-go/secrets"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
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

func TestSecretKeyValuesPlaintextIsNotAlwaysValidBase64(t *testing.T) {
	// Issue #91: the API base64-decodes each map value. Unencoded plaintext
	// either fails DecodeBytes or coincidentally decodes to the wrong bytes.
	invalid := []string{"192.0.2.10", "p#ss&w*rd"}
	for _, v := range invalid {
		_, err := tcsecrets.DecodeBytes("secretKeyValues", v)
		assert.Error(t, err, "plaintext %q must fail DecodeBytes; the API requires encoding", v)
	}

	coincidentallyValid := []string{"username", "5432", "databasename"}
	for _, v := range coincidentallyValid {
		decoded, err := tcsecrets.DecodeBytes("secretKeyValues", v)
		assert.NoError(t, err, "plaintext %q happens to be valid base64", v)
		assert.NotEqual(t, v, string(decoded), "coincidentally valid base64 must not round-trip to the original plaintext")
	}
}

func TestEncodeSecretKeyValues(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{name: "nil", input: nil},
		{name: "empty map", input: map[string]any{}},
		{
			name: "issue 91 reproduction values",
			input: map[string]any{
				"host":     "192.0.2.10",
				"port":     "5432",
				"dbname":   "databasename",
				"username": "username",
			},
		},
		{
			name: "special characters",
			input: map[string]any{
				"password": "p#ss&w*rd",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := encodeSecretKeyValues(tt.input)
			plain := convert.ConvertToMap(tt.input)
			assert.Len(t, encoded, len(plain))

			for k, v := range plain {
				decoded, err := tcsecrets.DecodeBytes("secretKeyValues["+k+"]", encoded[k])
				if assert.NoError(t, err, "value for %q must be valid base64", k) {
					assert.Equal(t, v, string(decoded))
				}
			}
		})
	}
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
