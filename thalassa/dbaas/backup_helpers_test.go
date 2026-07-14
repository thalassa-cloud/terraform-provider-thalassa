package dbaas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalassa-cloud/client-go/dbaas"
)

func TestIsBackupComplete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		backup *dbaas.DbClusterBackup
		want   bool
	}{
		{
			name:   "nil backup",
			backup: nil,
			want:   false,
		},
		{
			name:   "completed status",
			backup: &dbaas.DbClusterBackup{Status: dbaas.ObjectStatus("completed")},
			want:   true,
		},
		{
			name:   "ready status",
			backup: &dbaas.DbClusterBackup{Status: dbaas.ObjectStatusReady},
			want:   true,
		},
		{
			name:   "failed status",
			backup: &dbaas.DbClusterBackup{Status: dbaas.ObjectStatusFailed},
			want:   false,
		},
		{
			name:   "still creating",
			backup: &dbaas.DbClusterBackup{Status: dbaas.ObjectStatusCreating},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isBackupComplete(tt.backup))
		})
	}
}
