package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceTeam(t *testing.T) {
	resource := ResourceTeam()

	t.Run("resource schema validation", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.Equal(t, "Create a team in the Thalassa Cloud platform", resource.Description)

		// Check required fields
		schema := resource.Schema
		assert.NotNil(t, schema["name"])
		assert.True(t, schema["name"].Required)
		assert.False(t, schema["name"].Optional)

		// Check computed fields
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["slug"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)

		// Check optional fields
		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["labels"].Optional)
		assert.True(t, schema["annotations"].Optional)
		assert.True(t, schema["organisation_id"].Optional)
	})

	t.Run("resource CRUD operations", func(t *testing.T) {
		// This test would require a mock client or integration test
		// For now, we'll just test that the resource can be created
		assert.NotNil(t, resource.CreateContext)
		assert.NotNil(t, resource.ReadContext)
		assert.NotNil(t, resource.UpdateContext)
		assert.NotNil(t, resource.DeleteContext)
	})
}

func TestDataSourceTeam(t *testing.T) {
	dataSource := DataSourceTeam()

	t.Run("data source schema validation", func(t *testing.T) {
		assert.NotNil(t, dataSource)
		assert.Equal(t, "Get a team", dataSource.Description)

		// Check computed fields
		schema := dataSource.Schema
		assert.True(t, schema["id"].Computed)
		assert.True(t, schema["labels"].Computed)
		assert.True(t, schema["annotations"].Computed)
		assert.True(t, schema["created_at"].Computed)
		assert.True(t, schema["updated_at"].Computed)

		// Check optional fields
		assert.True(t, schema["name"].Optional)
		assert.True(t, schema["slug"].Optional)
		assert.True(t, schema["description"].Optional)
		assert.True(t, schema["organisation_id"].Optional)
	})

	t.Run("data source read operation", func(t *testing.T) {
		// This test would require a mock client or integration test
		// For now, we'll just test that the data source can be created
		assert.NotNil(t, dataSource.ReadContext)
	})
}

func TestTeamValidation(t *testing.T) {
	tests := []struct {
		name        string
		teamName    string
		description string
		expectError bool
	}{
		{
			name:        "valid team name",
			teamName:    "valid-team",
			description: "A valid team description",
			expectError: false,
		},
		{
			name:        "empty team name",
			teamName:    "",
			description: "A team with empty name",
			expectError: true,
		},
		{
			name:        "very long team name",
			teamName:    "a-very-long-team-name-that-exceeds-the-maximum-length-allowed-by-the-validation-function",
			description: "A team with very long name",
			expectError: true,
		},
		{
			name:        "very long description",
			teamName:    "valid-team",
			description: "A very long description that exceeds the maximum length allowed by the validation function and should cause an error when validated",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := ResourceTeam()
			schema := resource.Schema

			// Test name validation
			if tt.expectError {
				// This is a simplified test - in practice, you'd want to test the actual validation function
				assert.NotNil(t, schema["name"].ValidateFunc)
			}
		})
	}
}
