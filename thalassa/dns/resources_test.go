package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceDnsZone(t *testing.T) {
	resource := ResourceDnsZone()
	assert.True(t, resource.Schema["zone_name"].ForceNew)
	assert.NotNil(t, resource.Importer)
}

func TestResourceDnsRecord(t *testing.T) {
	resource := ResourceDnsRecord()
	schema := resource.Schema
	assert.True(t, schema["zone_id"].ForceNew)
	assert.True(t, schema["name"].ForceNew)
	assert.True(t, schema["type"].ForceNew)
	assert.False(t, schema["values"].ForceNew)
}
