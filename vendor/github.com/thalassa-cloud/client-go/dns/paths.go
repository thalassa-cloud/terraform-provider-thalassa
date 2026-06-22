package dns

import "strings"

const DnsEndpoint = "/v1/dns"

func zonePath(zoneIdentity string, segments ...string) string {
	parts := append([]string{DnsEndpoint, "zones", zoneIdentity}, segments...)
	return strings.Join(parts, "/")
}
