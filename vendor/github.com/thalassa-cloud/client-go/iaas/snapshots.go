package iaas

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	SnapshotEndpoint       = "/v1/snapshots"
	SnapshotPolicyEndpoint = "/v1/snapshot-policies"
)

type ListSnapshotsRequest struct {
	Filters []filters.Filter
}

// ListSnapshots lists all snapshots for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListSnapshots(ctx context.Context, listRequest *ListSnapshotsRequest) ([]Snapshot, error) {
	snapshots := []Snapshot{}
	req := c.R().SetResult(&snapshots)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, SnapshotEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return snapshots, err
	}

	return snapshots, nil
}

// GetSnapshot retrieves a specific snapshot by its identity.
// The identity is the unique identifier for the snapshot.
func (c *Client) GetSnapshot(ctx context.Context, identity string) (*Snapshot, error) {
	var snapshot *Snapshot
	req := c.R().SetResult(&snapshot)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", SnapshotEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return snapshot, err
	}

	return snapshot, nil
}

// CreateSnapshot creates a new snapshot.
func (c *Client) CreateSnapshot(ctx context.Context, create CreateSnapshotRequest) (*Snapshot, error) {
	var snapshot *Snapshot
	req := c.R().SetResult(&snapshot).SetBody(create)

	resp, err := c.Do(ctx, req, client.POST, SnapshotEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return snapshot, err
	}

	return snapshot, nil
}

// UpdateSnapshot updates a snapshot.
func (c *Client) UpdateSnapshot(ctx context.Context, identity string, update UpdateSnapshotRequest) (*Snapshot, error) {
	var snapshot *Snapshot
	req := c.R().SetResult(&snapshot).SetBody(update)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", SnapshotEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return snapshot, err
	}

	return snapshot, nil
}

// DeleteSnapshot deletes a snapshot.
func (c *Client) DeleteSnapshot(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", SnapshotEndpoint, identity))
	if err != nil {
		return err
	}

	if err := c.Check(resp); err != nil {
		return err
	}

	return nil
}

// WaitUntilSnapshotIsAvailable waits until a snapshot is available.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilSnapshotIsAvailable(ctx context.Context, snapshotIdentity string) error {
	return c.WaitUntilSnapshotIsStatus(ctx, snapshotIdentity, SnapshotStatusAvailable)
}

// WaitUntilSnapshotIsStatus waits until a snapshot is in a specific status.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilSnapshotIsStatus(ctx context.Context, snapshotIdentity string, status SnapshotStatus) error {
	snapshot, err := c.GetSnapshot(ctx, snapshotIdentity)
	if err != nil {
		return err
	}
	if snapshot.Status == status {
		return nil
	}
	// wait until the snapshot is in the desired status
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
		}

		snapshot, err = c.GetSnapshot(ctx, snapshotIdentity)
		if err != nil {
			return err
		}
		if snapshot.Status == status {
			return nil
		}
	}
}

// WaitUntilSnapshotIsDeleted waits until a snapshot is deleted.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilSnapshotIsDeleted(ctx context.Context, snapshotIdentity string) error {
	snapshot, err := c.GetSnapshot(ctx, snapshotIdentity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if snapshot.Status == SnapshotStatusDeleted {
		return nil
	}
	if snapshot.Status != SnapshotStatusDeleting {
		return fmt.Errorf("snapshot %s is not being deleted (status: %s)", snapshotIdentity, snapshot.Status)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			snapshot, err := c.GetSnapshot(ctx, snapshotIdentity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if snapshot.Status == SnapshotStatusDeleted {
				return nil
			}
		}
	}
}

type ListSnapshotPoliciesRequest struct {
	Filters []filters.Filter
}

// ListSnapshotPolicies lists all snapshot policies for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListSnapshotPolicies(ctx context.Context, listRequest *ListSnapshotPoliciesRequest) ([]SnapshotPolicy, error) {
	policies := []SnapshotPolicy{}
	req := c.R().SetResult(&policies)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, SnapshotPolicyEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return policies, err
	}

	return policies, nil
}

// GetSnapshotPolicy retrieves a specific snapshot policy by its identity.
// The identity is the unique identifier for the snapshot policy.
func (c *Client) GetSnapshotPolicy(ctx context.Context, identity string) (*SnapshotPolicy, error) {
	var policy *SnapshotPolicy
	req := c.R().SetResult(&policy)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", SnapshotPolicyEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return policy, err
	}

	return policy, nil
}

// CreateSnapshotPolicy creates a new snapshot policy.
func (c *Client) CreateSnapshotPolicy(ctx context.Context, create CreateSnapshotPolicyRequest) (*SnapshotPolicy, error) {
	var policy *SnapshotPolicy
	req := c.R().SetResult(&policy).SetBody(create)

	resp, err := c.Do(ctx, req, client.POST, SnapshotPolicyEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return policy, err
	}

	return policy, nil
}

// UpdateSnapshotPolicy updates a snapshot policy.
func (c *Client) UpdateSnapshotPolicy(ctx context.Context, identity string, update UpdateSnapshotPolicyRequest) (*SnapshotPolicy, error) {
	var policy *SnapshotPolicy
	req := c.R().SetResult(&policy).SetBody(update)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", SnapshotPolicyEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return policy, err
	}

	return policy, nil
}

// DeleteSnapshotPolicy deletes a snapshot policy.
func (c *Client) DeleteSnapshotPolicy(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", SnapshotPolicyEndpoint, identity))
	if err != nil {
		return err
	}

	if err := c.Check(resp); err != nil {
		return err
	}

	return nil
}

type CreateSnapshotRequest struct {
	// Name is the name of the snapshot. Must be unique within the organisation.
	// The name is used for identification and display purposes.
	Name string `json:"name"`
	// Description provides additional context about the snapshot.
	// This field is optional but recommended for better resource management.
	Description string `json:"description"`
	// Labels are key-value pairs that can be used for categorizing and filtering snapshots.
	// Common labels include environment (prod, staging, dev), project, or team.
	Labels Labels `json:"labels"`
	// Annotations are key-value pairs for storing additional metadata about the snapshot.
	// Unlike labels, annotations are not used for filtering or querying.
	Annotations Annotations `json:"annotations"`
	// VolumeIdentity is the unique identifier of the volume to create a snapshot from.
	// The volume must be in an available state and not currently attached to a machine.
	VolumeIdentity string `json:"volumeIdentity"`
	// DeleteProtection prevents the snapshot from being accidentally deleted.
	// When enabled, the snapshot cannot be deleted until this flag is set to false.
	DeleteProtection bool `json:"deleteProtection"`
}

type UpdateSnapshotRequest struct {
	// Name is the name of the snapshot. Must be unique within the organisation.
	// The name is used for identification and display purposes.
	Name string `json:"name"`
	// Description provides additional context about the snapshot.
	// This field is optional but recommended for better resource management.
	Description string `json:"description"`
	// Labels are key-value pairs that can be used for categorizing and filtering snapshots.
	// Common labels include environment (prod, staging, dev), project, or team.
	Labels Labels `json:"labels"`
	// Annotations are key-value pairs for storing additional metadata about the snapshot.
	// Unlike labels, annotations are not used for filtering or querying.
	Annotations Annotations `json:"annotations"`
	// DeleteProtection prevents the snapshot from being accidentally deleted.
	// When enabled, the snapshot cannot be deleted until this flag is set to false.
	DeleteProtection bool `json:"deleteProtection"`
}

type Snapshot struct {
	// Identity is the unique identifier for the snapshot across the entire system.
	// This field is automatically generated and cannot be modified.
	Identity string `json:"identity"`
	// Name is the human-readable name of the snapshot.
	// This field can be updated and is used for identification and display purposes.
	Name string `json:"name"`
	// Slug is a URL-friendly version of the name, automatically generated from the name.
	// This field is used for identification and display purposes.
	Slug string `json:"slug"`
	// Description provides additional context about the snapshot.
	// This field can be updated and is useful for documentation purposes.
	Description string `json:"description"`
	// CreatedAt is the timestamp when the snapshot was created.
	// This field is automatically set and cannot be modified.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the snapshot was last modified.
	// This field is automatically updated whenever the snapshot is modified.
	UpdatedAt time.Time `json:"updatedAt"`
	// ObjectVersion is used for optimistic locking to prevent concurrent modifications.
	// This field is automatically incremented on each update.
	ObjectVersion int `json:"objectVersion"`

	// Labels are key-value pairs for categorizing and filtering snapshots.
	// Common labels include environment (prod, staging, dev), project, or team.
	Labels Labels `json:"labels"`
	// Annotations are key-value pairs for storing additional metadata.
	// Unlike labels, annotations are not used for filtering or querying.
	Annotations Annotations `json:"annotations"`

	// Region contains information about the region where the snapshot is stored.
	// This field is populated when the snapshot is created.
	Region *Region `json:"region,omitempty"`
	// Organisation contains information about the organisation that owns the snapshot.
	// This field is populated when the snapshot is created.
	Organisation *base.Organisation `json:"organisation,omitempty"`

	// Status indicates the current state of the snapshot.
	Status SnapshotStatus `json:"status"`

	// SourceVolumeId is the unique identifier of the volume that this snapshot was created from.
	// This field is set when the snapshot is created and cannot be changed.
	SourceVolumeId *string `json:"sourceVolumeId"`
	// SourceVolume contains the full volume information if the volume still exists.
	// This field is populated when the volume is available and accessible.
	SourceVolume *Volume `json:"sourceVolume,omitempty"`

	// SizeGB is the size of the snapshot in GB.
	// This field is set once the snapshot creation is complete.
	SizeGB *int `json:"sizeGB,omitempty"`

	// DeleteProtection prevents the snapshot from being accidentally deleted.
	// When enabled, the snapshot cannot be deleted until this flag is set to false.
	DeleteProtection bool `json:"deleteProtection"`

	// SnapshotPolicyId is the unique identifier of the snapshot policy that created this snapshot.
	// This field is set when the snapshot is created by an automated policy.
	SnapshotPolicyId *string `json:"snapshotPolicyId,omitempty"`
	// SnapshotPolicy contains the full snapshot policy information if the policy still exists.
	// This field is populated when the policy is available and accessible.
	SnapshotPolicy *SnapshotPolicy `json:"snapshotPolicy,omitempty"`
}

// SnapshotStatus represents the current state of a snapshot.
type SnapshotStatus string

const (
	// SnapshotStatusCreating indicates that the snapshot is currently being created.
	// The snapshot is being copied from the source volume and is not yet available for use.
	SnapshotStatusCreating SnapshotStatus = "Creating"
	// SnapshotStatusAvailable indicates that the snapshot has been successfully created
	// and is ready for use. The snapshot can be used to create new volumes or restore data.
	SnapshotStatusAvailable SnapshotStatus = "Available"
	// SnapshotStatusDeleting indicates that the snapshot is currently being deleted.
	// The snapshot is being removed from storage and is no longer accessible.
	SnapshotStatusDeleting SnapshotStatus = "Deleting"
	// SnapshotStatusDeleted indicates that the snapshot has been successfully deleted.
	SnapshotStatusDeleted SnapshotStatus = "Deleted"
	// SnapshotStatusFailed indicates that the snapshot creation or deletion failed.
	// The snapshot is in an error state and may need manual intervention.
	SnapshotStatusFailed SnapshotStatus = "Failed"
)

type SnapshotPolicy struct {
	// Identity is the unique identifier for the snapshot policy across the entire system.
	// This field is automatically generated and cannot be modified.
	Identity string `json:"identity"`
	// Name is the human-readable name of the snapshot policy.
	// This field can be updated and is used for identification and display purposes.
	Name string `json:"name"`
	// Slug is a URL-friendly version of the name, automatically generated from the name.
	// This field is used in URLs and API endpoints.
	Slug string `json:"slug"`
	// Description provides additional context about the snapshot policy.
	// This field can be updated and is useful for documentation purposes.
	Description string `json:"description"`
	// CreatedAt is the timestamp when the snapshot policy was created.
	// This field is automatically set and cannot be modified.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the snapshot policy was last modified.
	// This field is automatically updated whenever the policy is modified.
	UpdatedAt time.Time `json:"updatedAt"`
	// ObjectVersion is used for optimistic locking to prevent concurrent modifications.
	// This field is automatically incremented on each update.
	ObjectVersion int `json:"objectVersion"`

	// Labels are key-value pairs for categorizing and filtering snapshot policies.
	// Common labels include environment (prod, staging, dev), project, or team.
	Labels Labels `json:"labels"`
	// Annotations are key-value pairs for storing additional metadata.
	// Unlike labels, annotations are not used for filtering or querying.
	Annotations Annotations `json:"annotations"`

	// Region contains information about the region where the snapshot policy operates.
	// This field is populated when the policy is created.
	Region *Region `json:"region,omitempty"`
	// Organisation contains information about the organisation that owns the policy.
	// This field is populated when the policy is created.
	Organisation *base.Organisation `json:"organisation,omitempty"`

	// Ttl (Time To Live) specifies how long snapshots created by this policy should be retained.
	// After this duration, snapshots will be automatically deleted unless protected.
	Ttl time.Duration `json:"ttl"`
	// KeepCount specifies the maximum number of snapshots to retain.
	// When this limit is reached, the oldest snapshots will be deleted.
	// If not set, the policy will keep all snapshots until the TTL expires.
	KeepCount *int `json:"keepCount,omitempty"`

	// Enabled indicates whether the snapshot policy is active and will create snapshots.
	// When disabled, the policy will not create new snapshots but existing ones remain.
	Enabled bool `json:"enabled"`

	// NextSnapshotAt indicates when the next snapshot will be created according to the schedule.
	// This field is calculated based on the schedule and timezone settings.
	NextSnapshotAt *time.Time `json:"nextSnapshotAt,omitempty"`

	// LastSnapshotAt indicates when the most recent snapshot was created by this policy.
	// This field is updated each time a snapshot is successfully created.
	LastSnapshotAt *time.Time `json:"lastSnapshotAt,omitempty"`

	// Schedule defines when snapshots should be created using a cron expression.
	// Examples: "0 2 * * *" (daily at 2 AM), "0 */6 * * *" (every 6 hours).
	Schedule string `json:"schedule"`

	// Timezone specifies the timezone for interpreting the schedule.
	// Must be a valid IANA timezone identifier (e.g., "UTC", "America/New_York").
	Timezone string `json:"timezone"`

	// Target is the target of the snapshot policy
	Target SnapshotPolicyTarget `json:"target"`

	// Snapshots is the list of snapshots created by the snapshot policy
	Snapshots []Snapshot `json:"snapshots,omitempty"`
}

// SnapshotPolicyTarget defines which volumes should be included in snapshots created by a policy.
type SnapshotPolicyTarget struct {
	// Type specifies how the target volumes are identified.
	// Use "selector" to target volumes based on labels, or "explicit" to target specific volumes.
	Type SnapshotPolicyTargetType `json:"type"`
	// Selector is a map of label key-value pairs used to identify volumes when Type is "selector".
	// Only volumes that match all the specified labels will be included in snapshots.
	// Example: {"environment": "production", "backup": "true"}
	Selector map[string]string `json:"selector"`
	// VolumeIdentities is a list of specific volume identifiers when Type is "explicit".
	// Only the volumes with these exact identities will be included in snapshots.
	VolumeIdentities []string `json:"volumeIdentities"`
}

// SnapshotPolicyTargetType defines the method used to identify target volumes for a snapshot policy.
type SnapshotPolicyTargetType string

const (
	// SnapshotPolicyTargetTypeSelector uses label selectors to dynamically identify volumes.
	// This approach is flexible and will automatically include new volumes that match the labels.
	SnapshotPolicyTargetTypeSelector SnapshotPolicyTargetType = "selector"
	// SnapshotPolicyTargetTypeExplicit uses a fixed list of volume identities.
	// This approach provides precise control over which volumes are included.
	SnapshotPolicyTargetTypeExplicit SnapshotPolicyTargetType = "explicit"
)

type CreateSnapshotPolicyRequest struct {
	// Name is the name of the snapshot policy. Must be unique within the organisation.
	// The name is used for identification and display purposes.
	Name string `json:"name"`
	// Description provides additional context about the snapshot policy.
	// This field is optional but recommended for better resource management.
	Description string `json:"description"`
	// Labels are key-value pairs that can be used for categorizing and filtering snapshot policies.
	// Common labels include environment (prod, staging, dev), project, or team.
	Labels Labels `json:"labels"`
	// Annotations are key-value pairs for storing additional metadata about the snapshot policy.
	// Unlike labels, annotations are not used for filtering or querying.
	Annotations Annotations `json:"annotations"`
	// Region is the identity or slug of the region where the snapshot policy will be created.
	// The policy will only operate on volumes in this region.
	Region string `json:"region"`
	// Ttl (Time To Live) specifies how long snapshots created by this policy should be retained.
	// After this duration, snapshots will be automatically deleted unless protected.
	Ttl time.Duration `json:"ttl"`
	// KeepCount specifies the maximum number of snapshots to retain.
	// When this limit is reached, the oldest snapshots will be deleted.
	// If not set, the policy will keep all snapshots until the TTL expires.
	KeepCount *int `json:"keepCount,omitempty"`
	// Enabled indicates whether the snapshot policy should be active immediately after creation.
	// When disabled, the policy will not create new snapshots until enabled.
	Enabled bool `json:"enabled"`
	// Schedule is the schedule of the snapshot policy. This is a cron expression.
	Schedule string `json:"schedule"`
	// Timezone is the timezone of the snapshot policy
	Timezone string `json:"timezone"`
	// Target is the target of the snapshot policy
	Target SnapshotPolicyTarget `json:"target"`
}

type UpdateSnapshotPolicyRequest struct {
	// Name is the name of the snapshot policy. Must be unique within the organisation.
	// The name is used for identification and display purposes.
	Name string `json:"name"`
	// Description provides additional context about the snapshot policy.
	// This field is optional but recommended for better resource management.
	Description string `json:"description"`
	// Labels are key-value pairs that can be used for categorizing and filtering snapshot policies.
	// Common labels include environment (prod, staging, dev), project, or team.
	Labels Labels `json:"labels"`
	// Annotations are key-value pairs for storing additional metadata about the snapshot policy.
	// Unlike labels, annotations are not used for filtering or querying.
	Annotations Annotations `json:"annotations"`
	// Ttl (Time To Live) specifies how long snapshots created by this policy should be retained.
	// After this duration, snapshots will be automatically deleted unless protected.
	Ttl time.Duration `json:"ttl"`
	// KeepCount specifies the maximum number of snapshots to retain.
	// When this limit is reached, the oldest snapshots will be deleted.
	// If not set, the policy will keep all snapshots until the TTL expires.
	KeepCount *int `json:"keepCount,omitempty"`
	// Enabled is a flag that indicates if the snapshot policy is enabled
	Enabled bool `json:"enabled"`
	// Schedule is the schedule of the snapshot policy. This is a cron expression.
	Schedule string `json:"schedule"`
	// Timezone is the timezone of the snapshot policy
	Timezone string `json:"timezone"`
	// Target is the target of the snapshot policy
	Target SnapshotPolicyTarget `json:"target"`
}
