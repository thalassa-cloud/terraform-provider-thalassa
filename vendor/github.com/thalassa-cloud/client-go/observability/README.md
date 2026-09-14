# Observability client

Go client for [Thalassa Cloud Observability workspaces](https://docs.thalassa.cloud).

Workspaces are org/project-scoped at `/v1/observability/workspaces` and cover
both Prometheus (metrics) and Loki (logs). See [package docs](https://pkg.go.dev/github.com/thalassa-cloud/client-go/observability)
and `example_test.go` for the full API.

This package replaces the former `/v1/observability/prometheus/tenants` client.
