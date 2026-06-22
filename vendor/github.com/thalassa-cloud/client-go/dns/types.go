package dns

import "time"

type DnsRecordType string

const (
	DnsRecordTypeTXT   DnsRecordType = "TXT"
	DnsRecordTypeA     DnsRecordType = "A"
	DnsRecordTypeCNAME DnsRecordType = "CNAME"
	DnsRecordTypeCAA   DnsRecordType = "CAA"
	DnsRecordTypeAAAA  DnsRecordType = "AAAA"
	DnsRecordTypeMX    DnsRecordType = "MX"
	DnsRecordTypeNS    DnsRecordType = "NS"
	DnsRecordTypeSRV   DnsRecordType = "SRV"
)

type DnsZone struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	ProjectID     string            `json:"projectId,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"`
	ObjectVersion int64             `json:"objectVersion"`
}

type DnsRecord struct {
	Identity  string        `json:"identity"`
	Name      string        `json:"name"`
	Type      DnsRecordType `json:"type"`
	TTL       int           `json:"ttl"`
	Values    []string      `json:"values"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt *time.Time    `json:"updatedAt,omitempty"`
}

type CreateDnsZoneRequest struct {
	ZoneName    string            `json:"zoneName"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type UpdateDnsZoneRequest struct {
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type CreateDnsRecordRequest struct {
	Name   string        `json:"name"`
	Type   DnsRecordType `json:"type"`
	TTL    int           `json:"ttl,omitempty"`
	Values []string      `json:"values"`
}

type UpdateDnsRecordRequest struct {
	TTL    int      `json:"ttl,omitempty"`
	Values []string `json:"values"`
}

type ExportDnsZoneFileResponse struct {
	ZoneName string `json:"zoneName"`
	ZoneFile string `json:"zoneFile"`
}

type ImportDnsZoneFileRequest struct {
	ZoneFile        string `json:"zoneFile"`
	ReplaceExisting bool   `json:"replaceExisting,omitempty"`
}

type ImportDnsZoneFileResponse struct {
	Created int         `json:"created"`
	Updated int         `json:"updated"`
	Deleted int         `json:"deleted"`
	Skipped int         `json:"skipped"`
	Records []DnsRecord `json:"records"`
}

type ListZonesRequest struct {
	Filters []ListZonesFilter
}

type ListRecordsRequest struct {
	Filters []ListRecordsFilter
}

type ListZonesFilter interface {
	ToParams() map[string]string
}

type ListRecordsFilter interface {
	ToParams() map[string]string
}

type DnsZoneDsRecord struct {
	Record         string `json:"record,omitempty"`
	DigestTypeName string `json:"digestTypeName,omitempty"`
	KeyTag         int    `json:"keyTag,omitempty"`
	Algorithm      int    `json:"algorithm,omitempty"`
	KeyRole        string `json:"keyRole,omitempty"`
	PublicKey      string `json:"publicKey,omitempty"`
}

type DnsZoneDnssecStatus struct {
	Enabled        bool              `json:"enabled"`
	DsDelegated    bool              `json:"dsDelegated"`
	DsRecords      []DnsZoneDsRecord `json:"dsRecords,omitempty"`
	LastSignedAt   *time.Time        `json:"lastSignedAt,omitempty"`
	LastSignError  *string           `json:"lastSignError,omitempty"`
	NextDsProbeAt  *time.Time        `json:"nextDsProbeAt,omitempty"`
	KmsKeyIdentity string            `json:"kmsKeyIdentity,omitempty"`
	Region         string            `json:"region,omitempty"`
}

type SetDnssecRequest struct {
	Region         string `json:"region"`
	KmsKeyIdentity string `json:"kmsKeyIdentity,omitempty"`
}
