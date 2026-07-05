package iaas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestResourceVpc(t *testing.T) {
	resource := resourceVpc()

	t.Run("schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Equal(t, "Create an vpc", resource.Description)

		s := resource.Schema
		assert.True(t, s["name"].Required)
		assert.True(t, s["name"].ForceNew)
		assert.NotNil(t, s["name"].ValidateFunc)

		assert.True(t, s["region"].Required)
		assert.True(t, s["region"].ForceNew)

		assert.True(t, s["cidrs"].Required)
		assert.Equal(t, schema.TypeList, s["cidrs"].Type)

		assert.True(t, s["description"].Optional)
		assert.NotNil(t, s["description"].ValidateFunc)

		assert.True(t, s["organisation_id"].Optional)
		assert.True(t, s["organisation_id"].ForceNew)

		assert.True(t, s["labels"].Optional)
		assert.True(t, s["annotations"].Optional)

		assert.True(t, s["slug"].Computed)
		assert.True(t, s["status"].Computed)
		assert.True(t, s["id"].Computed)

		assert.NotNil(t, resource.Importer)
	})

	t.Run("CRUD handlers", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
	})
}

func TestDataSourceVpc(t *testing.T) {
	dataSource := DataSourceVpc()

	t.Run("schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		assert.Equal(t, "Get an vpc", dataSource.Description)

		s := dataSource.Schema
		assert.True(t, s["id"].Computed)
		assert.True(t, s["slug"].Optional)
		assert.True(t, s["name"].Optional)
		assert.True(t, s["region"].Optional)
		assert.True(t, s["cidrs"].Computed)
		assert.True(t, s["status"].Computed)
		assert.True(t, s["labels"].Computed)
		assert.True(t, s["annotations"].Computed)
	})

	t.Run("read handler", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}
