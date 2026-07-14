package dbaas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRecoveryTargetTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:  "canonical barman format",
			input: "2023-08-11 11:14:21.00000+02",
		},
		{
			name:  "barman format with utc offset",
			input: "2026-07-06 10:00:00.00000+00",
		},
		{
			name:  "barman format without fractional seconds",
			input: "2026-07-06 10:00:00+00",
		},
		{
			name:  "barman format without timezone",
			input: "2026-07-06 10:00:00",
		},
		{
			name:  "rfc3339 with zulu suffix",
			input: "2026-07-06T10:00:00Z",
		},
		{
			name:  "rfc3339 with offset",
			input: "2026-07-06T10:00:00+00:00",
		},
		{
			name:    "invalid format",
			input:   "not-a-timestamp",
			wantErr: true,
		},
		{
			name:    "empty value",
			input:   "   ",
			wantErr: true,
			errMsg:  "target_time cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := parseRecoveryTargetTime(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				expectedMsg := tt.errMsg
				if expectedMsg == "" {
					expectedMsg = "target_time must be a valid barman format datetime"
				}
				assert.Contains(t, err.Error(), expectedMsg)
				return
			}

			assert.NoError(t, err)
			assert.False(t, parsed.IsZero())
		})
	}
}

func TestFormatBarmanTargetTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "canonical barman format is preserved",
			input: "2023-08-11 11:14:21.00000+02",
			want:  "2023-08-11 11:14:21.00000+02",
		},
		{
			name:  "rfc3339 is normalized to barman format",
			input: "2026-07-06T10:00:00Z",
			want:  "2026-07-06 10:00:00.00000+00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := parseRecoveryTargetTime(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, formatBarmanTargetTime(parsed))
		})
	}
}

func TestValidateBarmanTargetTimeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:  "valid canonical barman format",
			input: "2023-08-11 11:14:21.00000+02",
		},
		{
			name:  "valid rfc3339 format",
			input: "2026-07-06T10:00:00Z",
		},
		{
			name:  "empty string is ignored",
			input: "",
		},
		{
			name:    "invalid timestamp",
			input:   "yesterday",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, errs := validateBarmanTargetTimeString(tt.input, "target_time")
			if tt.wantErr {
				assert.Len(t, errs, 1)
				assert.Contains(t, errs[0].Error(), "target_time must be a valid barman format datetime")
				return
			}

			assert.Empty(t, errs)
		})
	}
}

func TestValidateTargetLSNString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:  "valid lsn",
			input: "0/1234567",
		},
		{
			name:  "valid uppercase hex lsn",
			input: "1A/ABCDEF",
		},
		{
			name:  "empty string is ignored",
			input: "",
		},
		{
			name:    "invalid lsn",
			input:   "1234567",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, errs := validateTargetLSNString(tt.input, "target_lsn")
			if tt.wantErr {
				assert.Len(t, errs, 1)
				assert.Contains(t, errs[0].Error(), "target_lsn must be a valid PostgreSQL LSN")
				return
			}

			assert.Empty(t, errs)
		})
	}
}

func TestValidateRestoreRecoveryTargetBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		block   map[string]any
		wantErr string
	}{
		{
			name: "target time only",
			block: map[string]any{
				"target_time": "2023-08-11 11:14:21.00000+02",
			},
		},
		{
			name: "target lsn only",
			block: map[string]any{
				"target_lsn": "0/1234567",
			},
		},
		{
			name:    "missing both targets",
			block:   map[string]any{},
			wantErr: "requires either target_time or target_lsn",
		},
		{
			name: "both targets set",
			block: map[string]any{
				"target_time": "2023-08-11 11:14:21.00000+02",
				"target_lsn":  "0/1234567",
			},
			wantErr: "cannot specify both target_time and target_lsn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRestoreRecoveryTargetBlock(tt.block)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
