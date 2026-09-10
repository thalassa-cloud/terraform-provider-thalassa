package projects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceProject(t *testing.T) {
	resource := ResourceProject()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Equal(t, "Create and manage a project in the Thalassa Cloud platform", resource.Description)

		schema := resource.Schema
		assert.NotNil(t, schema["name"])
		assert.True(t, schema["name"].Required)
		assert.False(t, schema["name"].Optional)

		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["slug"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["object_version"].Computed)

		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
		assert.True(t, schema["organisation_id"].Optional)
		assert.True(t, schema["organisation_id"].ForceNew)
		assert.True(t, schema["parent_project_identity"].Optional)
		assert.True(t, schema["parent_project_identity"].ForceNew)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
		assert.NotNil(t, resource.Importer)
	})
}

func TestDataSourceProject(t *testing.T) {
	dataSource := DataSourceProject()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		assert.Equal(t, "Get a project by identity, name, or slug", dataSource.Description)

		schema := dataSource.Schema
		assert.True(t, schema["id"].Optional)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["name"].Optional)
		assert.True(t, schema["slug"].Optional)
		assert.True(t, schema["description"].Computed)
		assert.True(t, schema["labels"].Computed)
		assert.True(t, schema["annotations"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["object_version"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}
