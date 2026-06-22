# DNS client

Go client for [Thalassa Cloud DNS](https://docs.thalassa.cloud/docs/dns/).

Zones and records are org/project-scoped at `/v1/dns`. See [package docs](https://pkg.go.dev/github.com/thalassa-cloud/client-go/dns) and `example_test.go` for the full API.

```go
tc, err := thalassa.NewClient(
    client.WithBaseURL("https://api.thalassa.cloud"),
    client.WithOrganisation("my-org"),
    client.WithAuthPersonalToken(os.Getenv("THALASSA_TOKEN")),
)
dnsClient := tc.DNS()
```

```go
zone, err := dnsClient.CreateZone(ctx, dns.CreateDnsZoneRequest{
    ZoneName: "example.com",
})
_, err = dnsClient.CreateRecord(ctx, zone.Identity, dns.CreateDnsRecordRequest{
    Name:   "www",
    Type:   dns.DnsRecordTypeA,
    TTL:    300,
    Values: []string{"192.0.2.1"},
})
```

Record name and type cannot change after create. Expect ~25–40s for DNS propagation after changes.
