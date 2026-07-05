package dns

import (
	"time"

	tcdns "github.com/thalassa-cloud/client-go/dns"
)

const timeFormatRFC3339 = time.RFC3339

var dnsRecordTypes = []string{
	string(tcdns.DnsRecordTypeA),
	string(tcdns.DnsRecordTypeAAAA),
	string(tcdns.DnsRecordTypeCNAME),
	string(tcdns.DnsRecordTypeMX),
	string(tcdns.DnsRecordTypeTXT),
	string(tcdns.DnsRecordTypeCAA),
	string(tcdns.DnsRecordTypeSRV),
	string(tcdns.DnsRecordTypeNS),
}

func setDnsZoneState(d interface {
	Set(string, any) error
}, zone *tcdns.DnsZone) error {
	_ = d.Set("zone_name", zone.Name)
	_ = d.Set("slug", zone.Slug)
	_ = d.Set("description", zone.Description)
	_ = d.Set("labels", zone.Labels)
	_ = d.Set("annotations", zone.Annotations)
	_ = d.Set("object_version", zone.ObjectVersion)
	if !zone.CreatedAt.IsZero() {
		_ = d.Set("created_at", zone.CreatedAt.Format(timeFormatRFC3339))
	}
	if zone.UpdatedAt != nil {
		_ = d.Set("updated_at", zone.UpdatedAt.Format(timeFormatRFC3339))
	}
	return nil
}

func setDnsRecordState(d interface {
	Set(string, any) error
}, record *tcdns.DnsRecord, zoneID string) error {
	_ = d.Set("zone_id", zoneID)
	_ = d.Set("name", record.Name)
	_ = d.Set("type", string(record.Type))
	_ = d.Set("ttl", record.TTL)
	_ = d.Set("values", record.Values)
	if !record.CreatedAt.IsZero() {
		_ = d.Set("created_at", record.CreatedAt.Format(timeFormatRFC3339))
	}
	if record.UpdatedAt != nil {
		_ = d.Set("updated_at", record.UpdatedAt.Format(timeFormatRFC3339))
	}
	return nil
}

func flattenDsRecords(records []tcdns.DnsZoneDsRecord) []map[string]any {
	result := make([]map[string]any, 0, len(records))
	for _, r := range records {
		result = append(result, map[string]any{
			"record":           r.Record,
			"digest_type_name": r.DigestTypeName,
			"key_tag":          r.KeyTag,
			"algorithm":        r.Algorithm,
			"key_role":         r.KeyRole,
			"public_key":       r.PublicKey,
		})
	}
	return result
}
