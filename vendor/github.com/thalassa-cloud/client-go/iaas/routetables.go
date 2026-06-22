package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	RouteTableEndpoint = "/v1/route-tables"
)

type ListRouteTablesRequest struct {
	Filters []filters.Filter
}

// ListRouteTables lists all RouteTables for a given organisation.
func (c *Client) ListRouteTables(ctx context.Context, listRequest *ListRouteTablesRequest) ([]RouteTable, error) {
	routeTables := []RouteTable{}
	req := c.R().SetResult(&routeTables)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, RouteTableEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return routeTables, err
	}
	return routeTables, nil
}

// GetRouteTable retrieves a specific RouteTable by its identity.
func (c *Client) GetRouteTable(ctx context.Context, identity string) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().SetResult(&routeTable)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
}

// CreateRouteTable creates a new RouteTable.
func (c *Client) CreateRouteTable(ctx context.Context, create CreateRouteTable) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().
		SetBody(create).SetResult(&routeTable)

	resp, err := c.Do(ctx, req, client.POST, RouteTableEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
}

// UpdateRouteTable updates an existing RouteTable.
func (c *Client) UpdateRouteTable(ctx context.Context, identity string, update UpdateRouteTable) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().
		SetBody(update).SetResult(&routeTable)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
}

// DeleteRouteTable deletes a specific RouteTable by its identity.
func (c *Client) DeleteRouteTable(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// CreateRouteTableRoute creates a new route for a specific RouteTable.
func (c *Client) CreateRouteTableRoute(ctx context.Context, identity string, create CreateRouteTableRoute) (*RouteEntry, error) {
	var routeEntry *RouteEntry
	req := c.R().
		SetBody(create).SetResult(&routeEntry)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/routes", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeEntry, err
	}
	return routeEntry, nil
}

// GetRouteTableRoute retrieves a specific route for a specific RouteTable.
func (c *Client) GetRouteTableRoute(ctx context.Context, identity string, routeIdentity string) (*RouteEntry, error) {
	var routeEntry *RouteEntry
	req := c.R().
		SetResult(&routeEntry)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/routes/%s", RouteTableEndpoint, identity, routeIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeEntry, err
	}
	return routeEntry, nil
}

// DeleteRouteTableRoute deletes a specific route for a specific RouteTable.
func (c *Client) DeleteRouteTableRoute(ctx context.Context, identity string, routeIdentity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/routes/%s", RouteTableEndpoint, identity, routeIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// UpdateRouteTableRoute updates a specific route for a specific RouteTable.
func (c *Client) UpdateRouteTableRoute(ctx context.Context, identity string, routeIdentity string, update UpdateRouteTableRoute) (*RouteEntry, error) {
	var routeEntry *RouteEntry
	req := c.R().
		SetBody(update).SetResult(&routeEntry)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/routes/%s", RouteTableEndpoint, identity, routeIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeEntry, err
	}
	return routeEntry, nil
}

// UpdateRouteTableRoutes updates the routes for a specific RouteTable.
func (c *Client) UpdateRouteTableRoutes(ctx context.Context, identity string, update UpdateRouteTableRoutes) ([]RouteEntry, error) {
	var routeEntries []RouteEntry
	req := c.R().
		SetBody(update).SetResult(&routeEntries)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/routes", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeEntries, err
	}
	return routeEntries, nil
}
