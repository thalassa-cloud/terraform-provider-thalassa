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
	assert.NotNil(t, resource.Importer)
}

func TestParseDnsRecordImportID(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		wantZoneID   string
		wantRecordID string
		wantErr      bool
	}{
		{
			name:         "valid composite id",
			id:           "dnsz-abc/dnsr-xyz",
			wantZoneID:   "dnsz-abc",
			wantRecordID: "dnsr-xyz",
		},
		{
			name:    "record id only",
			id:      "dnsr-xyz",
			wantErr: true,
		},
		{
			name:    "empty",
			id:      "",
			wantErr: true,
		},
		{
			name:    "missing record id",
			id:      "dnsz-abc/",
			wantErr: true,
		},
		{
			name:    "missing zone id",
			id:      "/dnsr-xyz",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zoneID, recordID, err := parseDnsRecordImportID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantZoneID, zoneID)
			assert.Equal(t, tt.wantRecordID, recordID)
		})
	}
}
