package dns

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ExportZoneFile exports a zone as BIND-format text.
func (c *Client) ExportZoneFile(ctx context.Context, zoneIdentity string) (*ExportDnsZoneFileResponse, error) {
	var exported ExportDnsZoneFileResponse
	r := c.R().SetResult(&exported)
	resp, err := c.Do(ctx, r, client.GET, zonePath(zoneIdentity, "export"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &exported, err
	}
	return &exported, nil
}

// ImportZoneFile imports records from BIND-format zone file text.
func (c *Client) ImportZoneFile(ctx context.Context, zoneIdentity string, importReq ImportDnsZoneFileRequest) (*ImportDnsZoneFileResponse, error) {
	var result ImportDnsZoneFileResponse
	r := c.R().SetBody(importReq).SetResult(&result)
	resp, err := c.Do(ctx, r, client.POST, zonePath(zoneIdentity, "import"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}
