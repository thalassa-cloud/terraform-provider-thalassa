package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceWorkspace(t *testing.T) {
	resource := ResourceWorkspace()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Equal(t, "Create and manage an observability workspace (metrics and logs) in Thalassa Cloud", resource.Description)

		schema := resource.Schema
		assert.True(t, schema["name"].Required)
		assert.True(t, schema["region"].Required)
		assert.True(t, schema["region"].ForceNew)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["status"].Computed)
		assert.True(t, schema["remote_write_url"].Computed)
		assert.True(t, schema["loki_query_url"].Computed)
		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
		assert.True(t, schema["retention_days"].Optional)
		assert.True(t, schema["retention_days"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
		assert.True(t, schema["organisation_id"].ForceNew)
		assert.True(t, schema["wait_for_deleted_timeout"].Optional)
		assert.Equal(t, 0, schema["wait_for_deleted_timeout"].Default)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
		assert.NotNil(t, resource.Importer)
	})
}

func TestDataSourceWorkspace(t *testing.T) {
	dataSource := DataSourceWorkspace()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		assert.Equal(t, "Get an observability workspace by identity or by name (optionally scoped by region)", dataSource.Description)

		schema := dataSource.Schema
		assert.True(t, schema["id"].Optional)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["name"].Optional)
		assert.True(t, schema["region"].Optional)
		assert.True(t, schema["region"].Computed)
		assert.True(t, schema["description"].Computed)
		assert.True(t, schema["status"].Computed)
		assert.True(t, schema["remote_write_url"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}

func TestParseWorkspaceImportID(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		expectedRegion   string
		expectedIdentity string
	}{
		{name: "identity only", id: "obsw-abc", expectedIdentity: "obsw-abc"},
		{name: "region and identity", id: "nl-01/obsw-abc", expectedRegion: "nl-01", expectedIdentity: "obsw-abc"},
		{name: "empty", id: "", expectedIdentity: ""},
		{name: "whitespace", id: "  ", expectedIdentity: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, identity := parseWorkspaceImportID(tt.id)
			assert.Equal(t, tt.expectedRegion, region)
			assert.Equal(t, tt.expectedIdentity, identity)
		})
	}
}
