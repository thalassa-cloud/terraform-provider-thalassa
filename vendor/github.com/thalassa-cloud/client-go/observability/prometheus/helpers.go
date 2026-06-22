package prometheus

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// validatePrometheusDuration validates that the string is a valid Prometheus duration format
// and that it is at least 1 day and not longer than 3 years
// Prometheus duration format: [0-9]+(ms|[smhdwy])
// Examples: "1h", "24h", "7d", "30d", "90d", "1y"
func validatePrometheusDuration(duration string) error {
	if strings.TrimSpace(duration) == "" {
		return fmt.Errorf("duration cannot be empty")
	}

	// Prometheus duration regex: number followed by unit (h, d, w, y)
	// The number must be positive
	pattern := `^([1-9]\d*)([hdwy])$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(duration)
	if matches == nil {
		return fmt.Errorf("duration must be in Prometheus format: [number][unit] where unit is one of: h, d, w, y (e.g., 7d, 30d, 1y)")
	}

	// Extract number and unit
	numberStr := matches[1]
	unit := matches[2]

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid duration number: %w", err)
	}

	// Convert to hours for comparison
	// Minimum: 1 day = 24 hours
	// Maximum: 3 years = 3 * 365 * 24 = 26280 hours (assuming 365 days per year)
	const minHours = 24    // 1 day
	const maxHours = 26280 // 3 years (3 * 365 * 24)

	var totalHours int64
	switch unit {
	case "h":
		totalHours = number
	case "d":
		totalHours = number * 24
	case "w":
		totalHours = number * 24 * 7
	case "y":
		// Assuming 365 days per year
		totalHours = number * 24 * 365
	default:
		return fmt.Errorf("unsupported duration unit: %s (supported: h, d, w, y)", unit)
	}

	if totalHours < minHours {
		return fmt.Errorf("retention must be at least 1 day (24h or 1d), got %s which is less than 1 day", duration)
	}

	if totalHours > maxHours {
		return fmt.Errorf("retention cannot be longer than 3 years (3y), got %s which exceeds 3 years", duration)
	}

	return nil
}
