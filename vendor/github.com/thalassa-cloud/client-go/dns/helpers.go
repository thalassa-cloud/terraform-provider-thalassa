package dns

import "fmt"

// FormatMX formats an MX record value.
func FormatMX(priority int, host string) string {
	return fmt.Sprintf("%d %s", priority, host)
}

// FormatCAA formats a CAA record value.
func FormatCAA(flags int, tag, value string) string {
	return fmt.Sprintf("%d %s %s", flags, tag, value)
}

// FormatSRV formats an SRV record value.
func FormatSRV(priority, weight, port int, target string) string {
	return fmt.Sprintf("%d %d %d %s", priority, weight, port, target)
}

// RecordFQDN returns the fully qualified record name for display.
func RecordFQDN(zoneName, recordName string) string {
	switch recordName {
	case "@":
		return zoneName
	case "*":
		return "*." + zoneName
	default:
		return recordName + "." + zoneName
	}
}
