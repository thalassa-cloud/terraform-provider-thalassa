package iaas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

// AttachVolumeAndWaitUntilAttached attaches a volume to a machine and waits until it is attached.
// The user is expected to provide a timeout context.
func (c *Client) AttachVolumeAndWaitUntilAttached(ctx context.Context, volumeIdentity string, attach AttachVolumeRequest) error {
	_, err := c.AttachVolume(ctx, volumeIdentity, attach)
	if err != nil {
		return err
	}
	return c.WaitUntilVolumeIsAttached(ctx, volumeIdentity)
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

// DetachVolumeAndWaitUntilAvailable detaches a volume from a machine and waits until it is available.
// The user is expected to provide a timeout context.
func (c *Client) DetachVolumeAndWaitUntilAvailable(ctx context.Context, volumeIdentity string, detach DetachVolumeRequest) error {
	err := c.DetachVolume(ctx, volumeIdentity, detach)
	if err != nil {
		return err
	}
	return c.WaitUntilVolumeIsAvailable(ctx, volumeIdentity)
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

func (c *Client) WaitUntilVolumeIsAttached(ctx context.Context, volumeIdentity string) error {
	return c.WaitUntilVolumeIsStatus(ctx, volumeIdentity, "attached")
}

func (c *Client) WaitUntilVolumeIsAvailable(ctx context.Context, volumeIdentity string) error {
	return c.WaitUntilVolumeIsStatus(ctx, volumeIdentity, "available")
}

// WaitUntilVolumeIsStatus waits until a volume is in a specific status.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilVolumeIsStatus(ctx context.Context, volumeIdentity string, status string) error {
	volume, err := c.GetVolume(ctx, volumeIdentity)
	if err != nil {
		return err
	}
	if strings.EqualFold(volume.Status, status) {
		return nil
	}
	// wait until the volume is unattached
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}

		volume, err = c.GetVolume(ctx, volumeIdentity)
		if err != nil {
			return err
		}
		if strings.EqualFold(volume.Status, status) {
			return nil
		}
	}
}

func (c *Client) WaitUntilVolumeIsDeleted(ctx context.Context, volumeIdentity string) error {
	volume, err := c.GetVolume(ctx, volumeIdentity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if strings.EqualFold(volume.Status, "deleted") {
		return nil
	}
	if !strings.EqualFold(volume.Status, "deleting") {
		return fmt.Errorf("volume %s is not being deleted (status: %s)", volumeIdentity, volume.Status)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			volume, err := c.GetVolume(ctx, volumeIdentity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if strings.EqualFold(volume.Status, "deleted") {
				return nil
			}
		}
	}
}
