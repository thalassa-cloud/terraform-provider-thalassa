# Secrets Manager client

Go client for the [Thalassa Cloud Secrets Manager](https://docs.thalassa.cloud/docs/secrets-manager/).

Paths are embedded in URLs as `/v1/secrets/{region}/secret{path}` (no extra slash before the path). See [package docs](https://pkg.go.dev/github.com/thalassa-cloud/client-go/secrets) and `example_test.go` for the full API.

```go
tc, err := thalassa.NewClient(
    client.WithBaseURL("https://api.thalassa.cloud"),
    client.WithOrganisation("my-org"),
    client.WithAuthPersonalToken(os.Getenv("THALASSA_TOKEN")),
)
secretsClient := tc.Secrets()
```

```go
ctx := context.Background()
region := "nl-01"

_, err := secretsClient.CreateSecret(ctx, region, secrets.CreateSecretRequest{
    Path:           "/app/prod/db/password",
    KmsKeyIdentity: "kms-abc123",
    SecretString:   secrets.EncodeBytes([]byte("super-secret")),
})

val, version, err := secretsClient.GetSecretString(ctx, region, "/app/prod/db/password", nil)
browse, err := secretsClient.BrowseSecrets(ctx, region, "/app/prod/")
```

Exactly one of `SecretString`, `SecretKeyValues`, or `GenerateSecret` is required on create. 
