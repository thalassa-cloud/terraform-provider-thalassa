package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListListeners lists all listeners for a specific loadbalancer.
func (c *Client) ListListeners(ctx context.Context, listRequest *ListLoadbalancerListenersRequest) ([]VpcLoadbalancerListener, error) {
	if listRequest == nil {
		return nil, fmt.Errorf("listRequest is required")
	}
	if listRequest.Loadbalancer == "" {
		return nil, fmt.Errorf("loadbalancer is required")
	}

	listeners := []VpcLoadbalancerListener{}
	req := c.R().SetResult(&listeners)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/listeners", LoadbalancerEndpoint, listRequest.Loadbalancer))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return listeners, err
	}
	return listeners, nil
}

// GetListener retrieves a specific loadbalancer listener by its identity.
func (c *Client) GetListener(ctx context.Context, getRequest GetLoadbalancerListenerRequest) (*VpcLoadbalancerListener, error) {
	if getRequest.Loadbalancer == "" {
		return nil, fmt.Errorf("loadbalancer is required")
	}
	if getRequest.Listener == "" {
		return nil, fmt.Errorf("listener is required")
	}

	var listener *VpcLoadbalancerListener
	req := c.R().SetResult(&listener)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/listeners/%s", LoadbalancerEndpoint, getRequest.Loadbalancer, getRequest.Listener))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return listener, err
	}
	return listener, nil
}

// CreateListener creates a new loadbalancer listener.
func (c *Client) CreateListener(ctx context.Context, loadbalancerID string, create CreateListener) (*VpcLoadbalancerListener, error) {
	var listener *VpcLoadbalancerListener
	req := c.R().
		SetBody(create).SetResult(&listener)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/listeners", LoadbalancerEndpoint, loadbalancerID))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return listener, err
	}
	return listener, nil
}

// UpdateListener updates an existing loadbalancer listener.
func (c *Client) UpdateListener(ctx context.Context, loadbalancerID string, listenerID string, update UpdateListener) (*VpcLoadbalancerListener, error) {
	var listener *VpcLoadbalancerListener
	req := c.R().
		SetBody(update).SetResult(&listener)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/listeners/%s", LoadbalancerEndpoint, loadbalancerID, listenerID))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return listener, err
	}
	return listener, nil
}

// DeleteListener deletes a specific loadbalancer listener by its identity.
func (c *Client) DeleteListener(ctx context.Context, loadbalancerID string, listenerID string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/listeners/%s", LoadbalancerEndpoint, loadbalancerID, listenerID))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

type VpcLoadbalancerListener struct {
	// Identity is the unique identifier for the listener.
	Identity string `json:"identity"`
	// Name is the name of the listener.
	Name string `json:"name"`
	// Slug is the slug of the listener.
	Slug string `json:"slug"`
	// Description is the description of the listener.
	Description string `json:"description"`
	// CreatedAt is the date and time the listener was created.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the date and time the listener was updated.
	UpdatedAt time.Time `json:"updatedAt"`
	// ObjectVersion is the version of the listener.
	ObjectVersion int `json:"objectVersion"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the listener, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the listener, and are used for matching resources.
	Annotations Annotations `json:"annotations,omitempty"`

	// Port is the port the listener is listening on.
	Port int `json:"port"`
	// Protocol is the protocol the listener is using.
	Protocol LoadbalancerProtocol `json:"protocol"`
	// TargetGroup is the target group attached to the listener.
	TargetGroup *VpcLoadbalancerTargetGroup `json:"targetGroup"`
	// MaxConnections is the maximum number of connections that the listener can handle
	MaxConnections *uint32 `json:"maxConnections,omitempty"`
	// ConnectionIdleTimeout is the amount of seconds used for configuring the idle connection timeout on a listener
	ConnectionIdleTimeout *uint32 `json:"connectionIdleTimeout,omitempty"`

	// AllowedSources is a list of CIDR blocks that are allowed to connect to the listener.
	AllowedSources []string `json:"allowedSources"`
}

type CreateListener struct {
	// Name is the name of the listener.
	Name string `json:"name"`
	// Description is the description of the listener.
	Description string `json:"description"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the listener, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the listener, and are used for matching resources.
	Annotations Annotations `json:"annotations,omitempty"`
	// Port is the port the listener is listening on.
	Port int `json:"port"`
	// Protocol is the protocol the listener is using.
	Protocol LoadbalancerProtocol `json:"protocol"`
	// MaxConnections is the maximum number of connections that the listener can handle
	MaxConnections *uint32 `json:"maxConnections,omitempty"`
	// ConnectionIdleTimeout is the amount of seconds used for configuring the idle connection timeout on a listener
	ConnectionIdleTimeout *uint32 `json:"connectionIdleTimeout,omitempty"`
	// TargetGroup is the target group attached to the listener.
	TargetGroup string `json:"targetGroup"`
	// AllowedSources is a list of CIDR blocks that are allowed to connect to the listener.
	AllowedSources []string `json:"allowedSources,omitempty"`
}

type UpdateListener struct {
	// Name is the name of the listener.
	Name string `json:"name"`
	// Description is the description of the listener.
	Description string `json:"description"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the listener, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the listener, and are used for matching resources.
	Annotations Annotations `json:"annotations,omitempty"`
	// Port is the port the listener is listening on.
	Port int `json:"port"`
	// Protocol is the protocol the listener is using.
	Protocol LoadbalancerProtocol `json:"protocol"`
	// TargetGroup is the target group attached to the listener.
	TargetGroup string `json:"targetGroup"`
	// MaxConnections is the maximum number of connections that the listener can handle
	MaxConnections *uint32 `json:"maxConnections,omitempty"`
	// ConnectionIdleTimeout is the amount of seconds used for configuring the idle connection timeout on a listener
	ConnectionIdleTimeout *uint32 `json:"connectionIdleTimeout,omitempty"`
	// AllowedSources is a list of CIDR blocks that are allowed to connect to the listener.
	AllowedSources []string `json:"allowedSources,omitempty"`
}

type ListLoadbalancerListenersRequest struct {
	// Loadbalancer is the identity of the loadbalancer to list listeners for.
	Loadbalancer string
	// Filters are the filters to apply to the list of listeners.
	Filters []filters.Filter
}

type GetLoadbalancerListenerRequest struct {
	// Loadbalancer is the identity of the loadbalancer to get a listener for.
	Loadbalancer string
	// Listener is the identity of the listener to get.
	Listener string
}
