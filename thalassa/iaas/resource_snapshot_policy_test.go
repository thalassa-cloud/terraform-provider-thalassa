package iaas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{name: "hours", input: "24h", want: 24 * time.Hour},
		{name: "go duration with minutes", input: "168h0m0s", want: 168 * time.Hour},
		{name: "days", input: "7d", want: 7 * 24 * time.Hour},
		{name: "weeks", input: "1w", want: 7 * 24 * time.Hour},
		{name: "invalid days", input: "xd", wantErr: true},
		{name: "invalid weeks", input: "xw", wantErr: true},
		{name: "invalid duration", input: "not-a-duration", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSuppressEquivalentDuration(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{name: "identical", old: "7d", new: "7d", want: true},
		{name: "equivalent day formats", old: "7d", new: "168h0m0s", want: true},
		{name: "equivalent week and days", old: "1w", new: "7d", want: true},
		{name: "different durations", old: "7d", new: "14d", want: false},
		{name: "invalid old", old: "nope", new: "7d", want: false},
		{name: "invalid new", old: "7d", new: "nope", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, suppressEquivalentDuration("ttl", tt.old, tt.new, nil))
		})
	}
}

func TestSnapshotPolicyTTLState(t *testing.T) {
	tests := []struct {
		name       string
		configured string
		ttl        time.Duration
		want       string
	}{
		{name: "preserves configured days", configured: "7d", ttl: 7 * 24 * time.Hour, want: "7d"},
		{name: "preserves configured weeks", configured: "1w", ttl: 7 * 24 * time.Hour, want: "1w"},
		{name: "falls back when empty", configured: "", ttl: 24 * time.Hour, want: "24h0m0s"},
		{name: "falls back when mismatched", configured: "1d", ttl: 7 * 24 * time.Hour, want: "168h0m0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, snapshotPolicyTTLState(tt.configured, tt.ttl))
		})
	}
}
