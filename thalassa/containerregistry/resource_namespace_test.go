package containerregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceNamespace(t *testing.T) {
	resource := ResourceNamespace()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Equal(t, "Create and manage a container registry namespace in Thalassa Cloud", resource.Description)

		schema := resource.Schema
		assert.NotNil(t, schema["region"])
		assert.True(t, schema["region"].Required)
		assert.True(t, schema["region"].ForceNew)
		assert.NotNil(t, schema["namespace"])
		assert.True(t, schema["namespace"].Required)
		assert.True(t, schema["namespace"].ForceNew)

		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["object_version"].Computed)
		assert.True(t, schema["total_size_bytes"].Computed)

		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
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

func TestDataSourceNamespace(t *testing.T) {
	dataSource := DataSourceNamespace()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		schema := dataSource.Schema
		assert.True(t, schema["id"].Optional)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["region"].Optional)
		assert.True(t, schema["namespace"].Optional)
		assert.True(t, schema["description"].Computed)
		assert.True(t, schema["total_size_bytes"].Computed)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}

func TestResourceNamespaceConfiguration(t *testing.T) {
	resource := ResourceNamespaceConfiguration()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		schema := resource.Schema
		assert.NotNil(t, schema["namespace_id"])
		assert.True(t, schema["namespace_id"].Required)
		assert.True(t, schema["namespace_id"].ForceNew)
		assert.NotNil(t, schema["visibility"])
		assert.True(t, schema["visibility"].Required)

		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["object_version"].Computed)

		assert.True(t, schema["retention_policy"].Optional)
		assert.Equal(t, 1, schema["retention_policy"].MaxItems)
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

func TestDataSourceNamespaceConfiguration(t *testing.T) {
	dataSource := DataSourceNamespaceConfiguration()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		schema := dataSource.Schema
		assert.True(t, schema["namespace_id"].Required)
		assert.True(t, schema["visibility"].Computed)
		assert.True(t, schema["retention_policy"].Computed)
		assert.True(t, schema["created_at"].Computed)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}
