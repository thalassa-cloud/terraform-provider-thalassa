package kms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceKmsKey(t *testing.T) {
	resource := ResourceKmsKey()

	t.Run("schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		schema := resource.Schema
		assert.True(t, schema["region"].Required)
		assert.True(t, schema["name"].Required)
		assert.True(t, schema["key_type"].Required)
		assert.Nil(t, schema["identity"])
		assert.NotNil(t, resource.Importer)
	})

	t.Run("CRUD handlers", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
	})
}

func TestDataSourceKmsKey(t *testing.T) {
	dataSource := DataSourceKmsKey()
	assert.Nil(t, dataSource.Schema["identity"])
}

func TestDataSourceKmsSummary(t *testing.T) {
	dataSource := DataSourceKmsSummary()
	assert.NotNil(t, dataSource.ReadContext)
	assert.True(t, dataSource.Schema["feature_enabled"].Computed)
}

func TestParseKmsKeyImportID(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		wantRegion   string
		wantIdentity string
	}{
		{
			name:         "composite id",
			id:           "nl-01/kms-abc123",
			wantRegion:   "nl-01",
			wantIdentity: "kms-abc123",
		},
		{
			name:         "identity only",
			id:           "kms-abc123",
			wantRegion:   "",
			wantIdentity: "kms-abc123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, identity := parseKmsKeyImportID(tt.id)
			assert.Equal(t, tt.wantRegion, region)
			assert.Equal(t, tt.wantIdentity, identity)
		})
	}
}
