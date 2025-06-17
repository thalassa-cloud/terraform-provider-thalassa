package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	VolumeEndpoint = "/v1/volumes"
)

type ListVolumesRequest struct {
	Filters []filters.Filter
}

// ListVolumes lists all volumes for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListVolumes(ctx context.Context, listRequest *ListVolumesRequest) ([]Volume, error) {
	volumes := []Volume{}
	req := c.R().SetResult(&volumes)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, VolumeEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return volumes, err
	}

	return volumes, nil
}

// GetVolume retrieves a specific volume by its identity.
// The identity is the unique identifier for the volume.
func (c *Client) GetVolume(ctx context.Context, identity string) (*Volume, error) {
	var volume *Volume
	req := c.R().SetResult(&volume)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VolumeEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return volume, err
	}

	return volume, nil
}

// CreateVolume creates a new volume.
func (c *Client) CreateVolume(ctx context.Context, create CreateVolume) (*Volume, error) {
	var volume *Volume
	req := c.R().SetResult(&volume).SetBody(create)

	resp, err := c.Do(ctx, req, client.POST, VolumeEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return volume, err
	}

	return volume, nil
}

// UpdateVolume updates a volume.
func (c *Client) UpdateVolume(ctx context.Context, identity string, update UpdateVolume) (*Volume, error) {
	var volume *Volume
	req := c.R().SetResult(&volume).SetBody(update)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", VolumeEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return volume, err
	}

	return volume, nil
}

// DeleteVolume deletes a volume.
func (c *Client) DeleteVolume(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", VolumeEndpoint, identity))
	if err != nil {
		return err
	}

	if err := c.Check(resp); err != nil {
		return err
	}

	return nil
}

// AttachVolume attaches a volume to a machine.
func (c *Client) AttachVolume(ctx context.Context, volumeIdentity string, attach AttachVolumeRequest) (*VolumeAttachment, error) {
	var attachment *VolumeAttachment
	req := c.R().SetResult(&attachment).SetBody(attach)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/attach", VolumeEndpoint, volumeIdentity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return attachment, err
	}

	return attachment, nil
}

// DetachVolume detaches a volume from a machine.
func (c *Client) DetachVolume(ctx context.Context, volumeIdentity string, detach DetachVolumeRequest) error {
	req := c.R().SetBody(detach)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/detach", VolumeEndpoint, volumeIdentity))
	if err != nil {
		return err
	}

	if err := c.Check(resp); err != nil {
		return err
	}

	return nil
}
