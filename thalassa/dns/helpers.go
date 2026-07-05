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
	if err := d.Set("zone_name", zone.Name); err != nil {
		return err
	}
	if err := d.Set("slug", zone.Slug); err != nil {
		return err
	}
	if err := d.Set("description", zone.Description); err != nil {
		return err
	}
	if err := d.Set("labels", zone.Labels); err != nil {
		return err
	}
	if err := d.Set("annotations", zone.Annotations); err != nil {
		return err
	}
	if err := d.Set("object_version", zone.ObjectVersion); err != nil {
		return err
	}
	if !zone.CreatedAt.IsZero() {
		if err := d.Set("created_at", zone.CreatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	if zone.UpdatedAt != nil {
		if err := d.Set("updated_at", zone.UpdatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	return nil
}

func setDnsRecordState(d interface {
	Set(string, any) error
}, record *tcdns.DnsRecord, zoneID string) error {
	if err := d.Set("zone_id", zoneID); err != nil {
		return err
	}
	if err := d.Set("name", record.Name); err != nil {
		return err
	}
	if err := d.Set("type", string(record.Type)); err != nil {
		return err
	}
	if err := d.Set("ttl", record.TTL); err != nil {
		return err
	}
	if err := d.Set("values", record.Values); err != nil {
		return err
	}
	if !record.CreatedAt.IsZero() {
		if err := d.Set("created_at", record.CreatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	if record.UpdatedAt != nil {
		if err := d.Set("updated_at", record.UpdatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
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
