package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceIamPolicy(t *testing.T) {
	resource := ResourceIamPolicy()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Contains(t, resource.Description, "private beta")
		assert.Contains(t, resource.Description, "projects")

		schema := resource.Schema
		assert.True(t, schema["name"].Required)
		assert.True(t, schema["name"].ForceNew)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["slug"].Computed)
		assert.True(t, schema["system"].Computed)
		assert.True(t, schema["is_read_only"].Computed)
		assert.True(t, schema["object_version"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)
		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
		assert.True(t, schema["replicate_to_children"].Optional)
		assert.True(t, schema["replicate_to_children"].ForceNew)
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

func TestResourceIamPolicyRule(t *testing.T) {
	resource := ResourceIamPolicyRule()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Contains(t, resource.Description, "private beta")

		schema := resource.Schema
		assert.True(t, schema["policy_id"].Required)
		assert.True(t, schema["policy_id"].ForceNew)
		assert.True(t, schema["resources"].Required)
		assert.True(t, schema["resources"].ForceNew)
		assert.True(t, schema["permissions"].Required)
		assert.True(t, schema["permissions"].ForceNew)
		assert.True(t, schema["resource_identities"].Optional)
		assert.True(t, schema["resource_identities"].ForceNew)
		assert.True(t, schema["note"].Optional)
		assert.True(t, schema["note"].ForceNew)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
		assert.True(t, schema["organisation_id"].ForceNew)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.Nil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
	})
}

func TestResourceIamPolicyBinding(t *testing.T) {
	resource := ResourceIamPolicyBinding()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Contains(t, resource.Description, "private beta")

		schema := resource.Schema
		assert.True(t, schema["policy_id"].Required)
		assert.True(t, schema["policy_id"].ForceNew)
		assert.True(t, schema["name"].Required)
		assert.True(t, schema["name"].ForceNew)
		assert.True(t, schema["user_id"].Optional)
		assert.True(t, schema["user_id"].ForceNew)
		assert.True(t, schema["service_account_id"].Optional)
		assert.True(t, schema["service_account_id"].ForceNew)
		assert.Equal(t, []string{"user_id", "service_account_id"}, schema["user_id"].ExactlyOneOf)
		assert.Equal(t, []string{"user_id", "service_account_id"}, schema["service_account_id"].ExactlyOneOf)
		assert.True(t, schema["expires_at"].Optional)
		assert.True(t, schema["expires_at"].ForceNew)
		assert.True(t, schema["duration"].Optional)
		assert.True(t, schema["duration"].ForceNew)
		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["slug"].Computed)
		assert.True(t, schema["object_version"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
		assert.True(t, schema["organisation_id"].ForceNew)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
	})
}

func TestDataSourceIamResourceAccessBindings(t *testing.T) {
	dataSource := DataSourceIamResourceAccessBindings()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		assert.Contains(t, dataSource.Description, "private beta")

		schema := dataSource.Schema
		assert.True(t, schema["resource_identity"].Required)
		assert.True(t, schema["resource_types"].Required)
		assert.True(t, schema["permissions"].Optional)
		assert.True(t, schema["binding_types"].Optional)
		assert.True(t, schema["bindings"].Computed)
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["organisation_id"].Optional)
	})

	t.Run("data source read operation", func(t *testing.T) {
		assert.NotNil(t, dataSource.ReadContext)
	})
}

func TestIamPolicyResourcesRegistered(t *testing.T) {
	tests := []struct {
		name     string
		resource string
	}{
		{name: "policy", resource: "thalassa_iam_policy"},
		{name: "policy rule", resource: "thalassa_iam_policy_rule"},
		{name: "policy binding", resource: "thalassa_iam_policy_binding"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, ok := ResourcesMap[tt.resource]
			assert.True(t, ok)
			assert.NotNil(t, resource)
		})
	}

	t.Run("resource access bindings data source", func(t *testing.T) {
		ds, ok := DataSourcesMap["thalassa_iam_resource_access_bindings"]
		assert.True(t, ok)
		assert.NotNil(t, ds)
	})
}
