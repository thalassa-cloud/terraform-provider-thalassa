package dns

import (
	"context"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListRecords lists DNS records in a zone.
func (c *Client) ListRecords(ctx context.Context, zoneIdentity string, req *ListRecordsRequest) ([]DnsRecord, error) {
	records := []DnsRecord{}
	r := c.R().SetResult(&records)
	if req != nil {
		for _, filter := range req.Filters {
			for k, v := range filter.ToParams() {
				r = r.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, r, client.GET, zonePath(zoneIdentity, "records"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return records, err
	}
	return records, nil
}

// CreateRecord creates a DNS record in a zone.
func (c *Client) CreateRecord(ctx context.Context, zoneIdentity string, create CreateDnsRecordRequest) (*DnsRecord, error) {
	var record DnsRecord
	r := c.R().SetBody(create).SetResult(&record)
	resp, err := c.Do(ctx, r, client.POST, zonePath(zoneIdentity, "records"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &record, err
	}
	return &record, nil
}

// GetRecord retrieves a DNS record by identity.
func (c *Client) GetRecord(ctx context.Context, zoneIdentity, recordIdentity string) (*DnsRecord, error) {
	var record DnsRecord
	r := c.R().SetResult(&record)
	resp, err := c.Do(ctx, r, client.GET, zonePath(zoneIdentity, "records", recordIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &record, err
	}
	return &record, nil
}

// UpdateRecord updates a DNS record TTL and values.
func (c *Client) UpdateRecord(ctx context.Context, zoneIdentity, recordIdentity string, update UpdateDnsRecordRequest) (*DnsRecord, error) {
	var record DnsRecord
	r := c.R().SetBody(update).SetResult(&record)
	resp, err := c.Do(ctx, r, client.PUT, zonePath(zoneIdentity, "records", recordIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &record, err
	}
	return &record, nil
}

// DeleteRecord deletes a DNS record.
func (c *Client) DeleteRecord(ctx context.Context, zoneIdentity, recordIdentity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, zonePath(zoneIdentity, "records", recordIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

type listRecordsFilter struct {
	filter filters.Filter
}

func (f listRecordsFilter) ToParams() map[string]string {
	return f.filter.ToParams()
}

// ListRecordsFilterFromFilter wraps a shared filter for ListRecords.
func ListRecordsFilterFromFilter(filter filters.Filter) ListRecordsFilter {
	return listRecordsFilter{filter: filter}
}
