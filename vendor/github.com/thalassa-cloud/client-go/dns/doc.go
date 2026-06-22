// Package dns provides a client for Thalassa Cloud DNS.
//
// Manage authoritative zones and records at /v1/dns. DNS is organisation- or
// project-scoped (no regional path segment in URLs). The dns feature gate must
// be enabled on your organisation.
//
// # Client setup
//
//	dnsClient := thalassa.NewClient(/* ... */).DNS()
//
// # Create a zone and record
//
//	zone, err := dnsClient.CreateZone(ctx, dns.CreateDnsZoneRequest{
//	    ZoneName: "example.com",
//	})
//	_, err = dnsClient.CreateRecord(ctx, zone.Identity, dns.CreateDnsRecordRequest{
//	    Name:   "www",
//	    Type:   dns.DnsRecordTypeA,
//	    TTL:    300,
//	    Values: []string{"192.0.2.1"},
//	})
//
// Record name and type are immutable after create; only TTL and values can be
// updated. Zone names cannot be renamed — create a new zone and import records.
//
// Changes publish to regional nameservers; expect ~25–40 seconds before
// authoritative answers update.
//
// See the package README and example_test.go for import/export, DNSSEC, and
// record value formats.
package dns
