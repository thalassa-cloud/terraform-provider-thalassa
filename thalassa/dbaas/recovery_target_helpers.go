package dbaas

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	barmanTargetTimeExample     = "2023-08-11 11:14:21.00000+02"
	barmanTargetTimeDescription = "Timestamp to restore to in barman format (YYYY-MM-DD HH:MM:SS.00000±TZ). Example: '" + barmanTargetTimeExample + "'"
	barmanTargetTimeLayout      = "2006-01-02 15:04:05.00000-07"
)

var (
	barmanTargetTimeLayouts = []string{
		barmanTargetTimeLayout,
		"2006-01-02 15:04:05.999999-07:00",
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05-07",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
	}
	rfc3339TargetTimeLayouts = []string{
		"2006-01-02T15:04:05.000000Z07:00",
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
	}
	targetLSNPattern = regexp.MustCompile(`^[0-9A-Fa-f]+/[0-9A-Fa-f]+$`)
)

func parseRecoveryTargetTime(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, fmt.Errorf("target_time cannot be empty")
	}

	for _, layout := range barmanTargetTimeLayouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			return parsed, nil
		}
	}

	for _, layout := range rfc3339TargetTimeLayouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"target_time must be a valid barman format datetime (YYYY-MM-DD HH:MM:SS.00000±TZ); got %q; example: %q",
		value,
		barmanTargetTimeExample,
	)
}

func formatBarmanTargetTime(value time.Time) string {
	return value.Format(barmanTargetTimeLayout)
}

func validateBarmanTargetTimeString(value any, _ string) (warns []string, errs []error) {
	raw, ok := value.(string)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	if _, err := parseRecoveryTargetTime(raw); err != nil {
		errs = append(errs, err)
	}

	return nil, errs
}

func validateTargetLSNString(value any, _ string) (warns []string, errs []error) {
	raw, ok := value.(string)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	if !targetLSNPattern.MatchString(strings.TrimSpace(raw)) {
		errs = append(errs, fmt.Errorf(
			"target_lsn must be a valid PostgreSQL LSN in the form <segment>/<offset>; got %q; example: %q",
			raw,
			"0/1234567",
		))
	}

	return nil, errs
}

func validateRestoreRecoveryTargetBlock(block map[string]any) error {
	targetTime, hasTargetTime := block["target_time"].(string)
	targetTime = strings.TrimSpace(targetTime)
	hasTargetTime = hasTargetTime && targetTime != ""

	targetLSN, hasTargetLSN := block["target_lsn"].(string)
	targetLSN = strings.TrimSpace(targetLSN)
	hasTargetLSN = hasTargetLSN && targetLSN != ""

	switch {
	case !hasTargetTime && !hasTargetLSN:
		return fmt.Errorf(
			"restore_recovery_target requires either target_time or target_lsn to be set",
		)
	case hasTargetTime && hasTargetLSN:
		return fmt.Errorf(
			"restore_recovery_target cannot specify both target_time and target_lsn; choose one recovery target",
		)
	}

	return nil
}
