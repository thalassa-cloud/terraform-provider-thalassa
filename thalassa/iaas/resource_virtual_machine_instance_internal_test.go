package iaas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	iaasclient "github.com/thalassa-cloud/client-go/iaas"
)

func TestSetMachineTypeField(t *testing.T) {
	t.Parallel()

	machineType := &iaasclient.MachineType{
		Identity: "mt-identity",
		Slug:     "mt-slug",
		Name:     "MT Name",
	}

	tests := []struct {
		name        string
		current     string
		expected    string
		description string
	}{
		{
			name:        "empty state uses identity",
			current:     "",
			expected:    "mt-identity",
			description: "uses API identity when state is empty",
		},
		{
			name:        "matching identity preserved",
			current:     "mt-identity",
			expected:    "mt-identity",
			description: "keeps identity reference when unchanged",
		},
		{
			name:        "matching slug preserved",
			current:     "mt-slug",
			expected:    "mt-slug",
			description: "keeps slug reference when unchanged",
		},
		{
			name:        "matching name preserved case-insensitively",
			current:     "mt name",
			expected:    "mt name",
			description: "keeps name reference when unchanged",
		},
		{
			name:        "drift detected from out-of-band change",
			current:     "old-machine-type",
			expected:    "mt-identity",
			description: "uses API identity when state no longer matches",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"machine_type": {Type: schema.TypeString, Optional: true},
			}, map[string]any{
				"machine_type": tt.current,
			})

			setMachineTypeField(d, machineType)
			assert.Equal(t, tt.expected, d.Get("machine_type"), tt.description)
		})
	}
}
