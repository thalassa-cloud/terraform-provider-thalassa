package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
)

func TestStringValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "non-nil pointer",
			input:    convert.Ptr("route table description"),
			expected: "route table description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, convert.StringValue(tt.input))
		})
	}
}
