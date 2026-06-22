package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	TargetGroupEndpoint = "/v1/loadbalancer-target-groups"
)

// ListTargetGroups lists all loadbalancer target groups for a given organisation.
func (c *Client) ListTargetGroups(ctx context.Context, listRequest *ListTargetGroupsRequest) ([]VpcLoadbalancerTargetGroup, error) {
	targetGroups := []VpcLoadbalancerTargetGroup{}
	req := c.R().SetResult(&targetGroups)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, TargetGroupEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return targetGroups, err
	}
	return targetGroups, nil
}

// GetTargetGroup retrieves a specific loadbalancer target group by its identity.
func (c *Client) GetTargetGroup(ctx context.Context, getRequest GetTargetGroupRequest) (*VpcLoadbalancerTargetGroup, error) {
	if getRequest.Identity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var targetGroup *VpcLoadbalancerTargetGroup
	req := c.R().SetResult(&targetGroup)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", TargetGroupEndpoint, getRequest.Identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return targetGroup, err
	}
	return targetGroup, nil
}

// CreateTargetGroup creates a new loadbalancer target group.
func (c *Client) CreateTargetGroup(ctx context.Context, create CreateTargetGroup) (*VpcLoadbalancerTargetGroup, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if create.Vpc == "" {
		return nil, fmt.Errorf("vpc is required")
	}
	if create.TargetPort == 0 {
		return nil, fmt.Errorf("targetPort is required")
	}
	if create.Protocol == "" {
		return nil, fmt.Errorf("protocol is required")
	}

	var targetGroup *VpcLoadbalancerTargetGroup
	req := c.R().
		SetBody(create).SetResult(&targetGroup)

	resp, err := c.Do(ctx, req, client.POST, TargetGroupEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return targetGroup, err
	}
	return targetGroup, nil
}

// UpdateTargetGroup updates an existing loadbalancer target group.
func (c *Client) UpdateTargetGroup(ctx context.Context, update UpdateTargetGroupRequest) (*VpcLoadbalancerTargetGroup, error) {
	if update.Identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	var targetGroup *VpcLoadbalancerTargetGroup
	req := c.R().
		SetBody(update.UpdateTargetGroup).SetResult(&targetGroup)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", TargetGroupEndpoint, update.Identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return targetGroup, err
	}
	return targetGroup, nil
}

// DeleteTargetGroup deletes a specific loadbalancer target group by its identity.
func (c *Client) DeleteTargetGroup(ctx context.Context, deleteRequest DeleteTargetGroupRequest) error {
	if deleteRequest.Identity == "" {
		return fmt.Errorf("identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", TargetGroupEndpoint, deleteRequest.Identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// SetTargetGroupServerAttachments sets the server attachments for a target group.
// This will replace the existing attachments with the ones provided in the request.
// Note: Any existing attachments not present in the request will be detached.
func (c *Client) SetTargetGroupServerAttachments(ctx context.Context, setRequest TargetGroupAttachmentsBatch) error {
	if setRequest.TargetGroupID == "" {
		return fmt.Errorf("targetGroupID is required")
	}
	req := c.R().SetBody(setRequest)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/attachments", TargetGroupEndpoint, setRequest.TargetGroupID))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// AttachServerToTargetGroup attaches a server to a target group.
func (c *Client) AttachServerToTargetGroup(ctx context.Context, attachRequest AttachTargetGroupRequest) (*LoadbalancerTargetGroupAttachment, error) {
	if attachRequest.ServerIdentity == "" && attachRequest.EndpointIdentity == "" {
		return nil, fmt.Errorf("serverIdentity or endpointIdentity is required")
	}
	if attachRequest.TargetGroupID == "" {
		return nil, fmt.Errorf("targetGroupID is required")
	}

	var result *LoadbalancerTargetGroupAttachment
	req := c.R().
		SetBody(attachRequest.AttachTarget).SetResult(&result)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/attach", TargetGroupEndpoint, attachRequest.TargetGroupID))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// DetachServerFromTargetGroup detaches a server from a target group.
func (c *Client) DetachServerFromTargetGroup(ctx context.Context, detachRequest DetachTargetRequest) error {
	if detachRequest.TargetGroupID == "" {
		return fmt.Errorf("targetGroupID is required")
	}
	if detachRequest.AttachmentID == "" {
		return fmt.Errorf("attachmentID is required")
	}

	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/detach/%s", TargetGroupEndpoint, detachRequest.TargetGroupID, detachRequest.AttachmentID))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

type VpcLoadbalancerTargetGroup struct {
	// Identity is the identity of the target group.
	Identity string `json:"identity"`
	// Name is the name of the target group.
	Name string `json:"name"`
	// Slug is the slug of the target group.
	Slug string `json:"slug"`
	// Description is a human-readable description of the target group.
	Description string `json:"description"`
	// CreatedAt is the time the target group was created.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the time the target group was updated.
	UpdatedAt time.Time `json:"updatedAt"`
	// ObjectVersion is the version of the target group.
	ObjectVersion int `json:"objectVersion"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the target group, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the target group.
	Annotations Annotations `json:"annotations,omitempty"`

	// Organisation is the organisation the target group belongs to.
	Organisation *base.Organisation `json:"organisation"`
	// Vpc is the VPC the target group belongs to.
	Vpc *Vpc `json:"vpc"`
	// TargetPort is the port to use for the target group.
	TargetPort int `json:"targetPort"`
	// Protocol is the protocol to use for the target group.
	Protocol LoadbalancerProtocol `json:"protocol"`
	// TargetSelector is a map of labels to match against the server labels.
	// If a server matches the labels, it will be added to the target group.
	// If no target selector is provided, target must be assigned manually
	TargetSelector map[string]string `json:"targetSelector"`

	// EnableProxyProtocol enables proxy protocol on the target group. When enabled, the load balancer will use the proxy protocol to communicate with the target group.
	// Enabling proxy protocl means all targets within the target group must support proxy protocol, otherwise connections may fail.
	EnableProxyProtocol *bool `json:"enableProxyProtocol,omitempty"`

	// LoadbalancingPolicy is the load balancing policy for the target group. Must be one of ROUND_ROBIN, RANDOM, or MAGLEV.
	// The default policy is ROUND_ROBIN.
	// ROUND_ROBIN: Connections from a listener to the target group are distributed across all target group attachments.
	// RANDOM: Connections from a listener to the target group are distributed across all target group attachments in a random manner.
	// MAGLEV: Connections from a listener to the target group are distributed across all target group attachments based on the MAGLEV algorithm.
	// +optional
	LoadbalancingPolicy *LoadbalancingPolicy `json:"loadbalancingPolicy,omitempty"`

	// HealthCheck is the health check settings for the target group
	// +optional
	HealthCheck *BackendHealthCheck `json:"healthCheck,omitempty"`

	// LoadbalancerListeners are the listeners the target group is attached to.
	LoadbalancerListeners []VpcLoadbalancerListener `json:"loadbalancerListeners"`
	// LoadbalancerTargetGroupAttachments are the attachments to the target group.
	LoadbalancerTargetGroupAttachments []LoadbalancerTargetGroupAttachment `json:"loadbalancerTargetGroupAttachments"`
}

type LoadbalancerTargetGroupAttachment struct {
	// Identity is the identity of the attachment.
	Identity string `json:"identity"`
	// CreatedAt is the time the attachment was created.
	CreatedAt time.Time `json:"createdAt"`
	// LoadbalancerTargetGroup is the target group the attachment belongs to.
	LoadbalancerTargetGroup *VpcLoadbalancerTargetGroup `json:"loadbalancerTargetGroup"`
	// VirtualMachineInstance is the server the attachment belongs to. Either VirtualMachineInstance or Endpoint is set.
	VirtualMachineInstance *Machine `json:"virtualMachineInstance,omitempty"`
	// Endpoint is the endpoint the attachment belongs to. Either VirtualMachineInstance or Endpoint is set.
	Endpoint *Endpoint `json:"endpoint,omitempty"`
}

type DetachTargetRequest struct {
	// TargetGroupID is the identity of the target group to detach the server from.
	TargetGroupID string `json:"targetGroupID"`
	// AttachmentID is the identity of the LoadbalancerTargetGroupAttachment to detach.
	AttachmentID string `json:"attachmentID"`
}

type GetTargetGroupRequest struct {
	Identity string
}

type ListTargetGroupsRequest struct {
	Filters []filters.Filter
}

// CreateTargetGroupRequest is the request body for creating a new loadbalancer target group.
type CreateTargetGroup struct {
	// Name is the name of the target group.
	Name string `json:"name"`
	// Description is a human-readable description of the target group.
	Description string `json:"description"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the target group, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the target group.
	Annotations Annotations `json:"annotations,omitempty"`
	// Vpc is the identity of the VPC to create the target group in.
	Vpc string `json:"vpc"`
	// TargetPort is the port to use for the target group. Must be between 1 and 65535.
	TargetPort int `json:"targetPort"`
	// Protocol is the protocol to use for the target group. Must be one of TCP, UDP, HTTP, or HTTPS.
	Protocol LoadbalancerProtocol `json:"protocol"`
	// TargetSelector is a map of labels to match against the server labels.
	// If a server matches the labels, it will be added to the target group.
	// If no target selector is provided, target must be assigned manually
	// +optional
	TargetSelector map[string]string `json:"targetSelector,omitempty"`

	// EnableProxyProtocol enables proxy protocol on the target group. When enabled, the load balancer will use the proxy protocol to communicate with the target group.
	// Enabling proxy protocl means all targets within the target group must support proxy protocol, otherwise connections may fail.
	EnableProxyProtocol *bool `json:"enableProxyProtocol,omitempty"`

	// LoadbalancingPolicy is the load balancing policy for the target group. Must be one of ROUND_ROBIN, RANDOM, or MAGLEV.
	// The default policy is ROUND_ROBIN.
	// ROUND_ROBIN: Connections from a listener to the target group are distributed across all target group attachments.
	// RANDOM: Connections from a listener to the target group are distributed across all target group attachments in a random manner.
	// MAGLEV: Connections from a listener to the target group are distributed across all target group attachments based on the MAGLEV algorithm.
	// +optional
	LoadbalancingPolicy *LoadbalancingPolicy `json:"loadbalancingPolicy,omitempty"`

	// HealthCheck is the health check settings for the target group
	// +optional
	HealthCheck *BackendHealthCheck `json:"healthCheck,omitempty"`
}

type LoadbalancingPolicy string

const (
	LoadbalancingPolicyRoundRobin LoadbalancingPolicy = "ROUND_ROBIN"
	LoadbalancingPolicyRandom     LoadbalancingPolicy = "RANDOM"
	LoadbalancingPolicyMagLev     LoadbalancingPolicy = "MAGLEV"
)

type BackendHealthCheck struct {
	// Protocol is the protocol to use for the health check
	Protocol LoadbalancerProtocol `json:"protocol"`
	// Port is the port to use for the health check. Must be between 1 and 65535.
	Port int32 `json:"port"`
	// Path is the path to use for the health check
	// If provided, must be a valid URL path.
	Path string `json:"path"`
	// Interval is the interval for the health check. Time is in seconds.
	// Minimum value is 5, maximum value is 300.
	PeriodSeconds int `json:"periodSeconds"`
	// Timeout is the timeout for the health check. Time is in seconds
	// Minimum value is 1, maximum value is 300.
	TimeoutSeconds int `json:"timeoutSeconds"`
	// UnhealthyThreshold is the number of failures before marking the server as unhealthy
	// Minimum value is 1, maximum value is 10.
	UnhealthyThreshold int32 `json:"unhealthyThreshold"`
	// HealthyThreshold is the number of successes before marking the server as healthy
	// Minimum value is 1, maximum value is 10.
	HealthyThreshold int32 `json:"healthyThreshold"`
}

type UpdateTargetGroupRequest struct {
	Identity string
	UpdateTargetGroup
}

type UpdateTargetGroup struct {
	// Name is the name of the target group.
	Name string `json:"name"`
	// Description is a human-readable description of the target group.
	Description string `json:"description"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the target group, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the target group.
	Annotations Annotations `json:"annotations,omitempty"`
	// TargetPort is the port to use for the target group. Must be between 1 and 65535.
	TargetPort int `json:"targetPort"`
	// Protocol is the protocol to use for the target group. Must be one of TCP, UDP, HTTP, or HTTPS.
	Protocol LoadbalancerProtocol `json:"protocol"`
	// TargetSelector is a map of labels to match against the server labels.
	// If a server matches the labels, it will be added to the target group.
	// If no target selector is provided, target must be assigned manually
	// +optional
	TargetSelector map[string]string `json:"targetSelector,omitempty"`

	// EnableProxyProtocol enables proxy protocol on the target group. When enabled, the load balancer will use the proxy protocol to communicate with the target group.
	// Enabling proxy protocl means all targets within the target group must support proxy protocol, otherwise connections may fail.
	EnableProxyProtocol *bool `json:"enableProxyProtocol,omitempty"`

	// LoadbalancingPolicy is the load balancing policy for the target group. Must be one of ROUND_ROBIN, RANDOM, or MAGLEV.
	// The default policy is ROUND_ROBIN.
	// ROUND_ROBIN: Connections from a listener to the target group are distributed across all target group attachments.
	// RANDOM: Connections from a listener to the target group are distributed across all target group attachments in a random manner.
	// MAGLEV: Connections from a listener to the target group are distributed across all target group attachments based on the MAGLEV algorithm.
	// +optional
	LoadbalancingPolicy *LoadbalancingPolicy `json:"loadbalancingPolicy,omitempty"`

	// HealthCheck is the health check settings for the target group
	// +optional
	HealthCheck *BackendHealthCheck `json:"healthCheck,omitempty"`
}

type AttachTarget struct {
	// ServerIdentity is the identity of the server to attach.
	ServerIdentity string `json:"serverIdentity"`
	// EndpointIdentity is the identity of the endpoint to attach.
	EndpointIdentity string `json:"endpointIdentity,omitempty"`
}

type AttachTargetGroupRequest struct {
	TargetGroupID string
	AttachTarget
}

type DeleteTargetGroupRequest struct {
	Identity string
}

type TargetGroupAttachmentsBatch struct {
	TargetGroupID string
	Attachments   []AttachTarget `json:"attachments"`
}
