package convert_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
)

func TestSetReferenceField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		current  string
		identity string
		slug     string
		refName  string
		expected string
	}{
		{
			name:     "empty state uses slug when available",
			current:  "",
			identity: "id-1",
			slug:     "slug-1",
			refName:  "Name One",
			expected: "slug-1",
		},
		{
			name:     "empty state uses name when slug is unavailable",
			current:  "",
			identity: "id-1",
			slug:     "",
			refName:  "Name One",
			expected: "Name One",
		},
		{
			name:     "empty state uses identity when no slug or name",
			current:  "",
			identity: "id-1",
			slug:     "",
			refName:  "",
			expected: "id-1",
		},
		{
			name:     "matching identity preserved",
			current:  "id-1",
			identity: "id-1",
			slug:     "slug-1",
			refName:  "Name One",
			expected: "id-1",
		},
		{
			name:     "matching slug preserved",
			current:  "slug-1",
			identity: "id-1",
			slug:     "slug-1",
			refName:  "Name One",
			expected: "slug-1",
		},
		{
			name:     "matching name preserved case-insensitively",
			current:  "name one",
			identity: "id-1",
			slug:     "slug-1",
			refName:  "Name One",
			expected: "name one",
		},
		{
			name:     "drift detected from out-of-band change",
			current:  "old-reference",
			identity: "id-1",
			slug:     "slug-1",
			refName:  "Name One",
			expected: "id-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"reference": {Type: schema.TypeString, Optional: true},
			}, map[string]interface{}{
				"reference": tt.current,
			})

			convert.SetReferenceField(d, "reference", tt.identity, tt.slug, tt.refName)
			assert.Equal(t, tt.expected, d.Get("reference"))
		})
	}
}
