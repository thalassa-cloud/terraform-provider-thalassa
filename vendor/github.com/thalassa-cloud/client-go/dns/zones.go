package dns

import (
	"context"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListZones lists DNS zones for the organisation or project scope.
func (c *Client) ListZones(ctx context.Context, req *ListZonesRequest) ([]DnsZone, error) {
	zones := []DnsZone{}
	r := c.R().SetResult(&zones)
	if req != nil {
		for _, filter := range req.Filters {
			for k, v := range filter.ToParams() {
				r = r.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, r, client.GET, DnsEndpoint+"/zones")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return zones, err
	}
	return zones, nil
}

// CreateZone creates a DNS zone.
func (c *Client) CreateZone(ctx context.Context, create CreateDnsZoneRequest) (*DnsZone, error) {
	var zone DnsZone
	r := c.R().SetBody(create).SetResult(&zone)
	resp, err := c.Do(ctx, r, client.POST, DnsEndpoint+"/zones")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &zone, err
	}
	return &zone, nil
}

// GetZone retrieves a DNS zone by identity.
func (c *Client) GetZone(ctx context.Context, zoneIdentity string) (*DnsZone, error) {
	var zone DnsZone
	r := c.R().SetResult(&zone)
	resp, err := c.Do(ctx, r, client.GET, zonePath(zoneIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &zone, err
	}
	return &zone, nil
}

// UpdateZone updates zone metadata.
func (c *Client) UpdateZone(ctx context.Context, zoneIdentity string, update UpdateDnsZoneRequest) (*DnsZone, error) {
	var zone DnsZone
	r := c.R().SetBody(update).SetResult(&zone)
	resp, err := c.Do(ctx, r, client.PUT, zonePath(zoneIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &zone, err
	}
	return &zone, nil
}

// DeleteZone deletes a DNS zone and all of its records.
func (c *Client) DeleteZone(ctx context.Context, zoneIdentity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, zonePath(zoneIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

type listZonesFilter struct {
	filter filters.Filter
}

func (f listZonesFilter) ToParams() map[string]string {
	return f.filter.ToParams()
}

// ListZonesFilterFromFilter wraps a shared filter for ListZones.
func ListZonesFilterFromFilter(filter filters.Filter) ListZonesFilter {
	return listZonesFilter{filter: filter}
}
