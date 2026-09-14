// Package observability provides a client for Thalassa Cloud Observability workspaces.
//
// Workspaces unify Prometheus (metrics) and Loki (logs) under one identity at
// /v1/observability/workspaces. Create a workspace, then use the endpoint URLs
// on the resource (remoteWriteUrl, pushUrl, prometheusQueryUrl, lokiQueryUrl, …)
// with your preferred Prometheus/Loki clients.
//
// Query/proxy routes under /{identity}/proxy/… are not wrapped by this package.
//
//	obs := thalassa.NewClient(/* ... */).Observability()
//	ws, err := obs.CreateObservabilityWorkspace(ctx, observability.CreateObservabilityWorkspaceRequest{
//	    Name:           "prod-metrics",
//	    RegionIdentity: "nl-01",
//	})
package observability
