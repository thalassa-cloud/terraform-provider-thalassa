# KMS client

Go client for the [Thalassa Cloud Key Management Service](https://docs.thalassa.cloud/docs/kms/).

Regional calls take an explicit region (e.g. `nl-01`). See [package docs](https://pkg.go.dev/github.com/thalassa-cloud/client-go/kms) and `example_test.go` for the full API.

```go
tc, err := thalassa.NewClient(
    client.WithBaseURL("https://api.thalassa.cloud"),
    client.WithOrganisation("my-org"),
    client.WithAuthPersonalToken(os.Getenv("THALASSA_TOKEN")),
)
kmsClient := tc.KMS()
```

```go
ctx := context.Background()
region := "nl-01"

summary, err := kmsClient.GetSummary(ctx)
if err != nil || !summary.FeatureEnabled {
    log.Fatal("KMS not available")
}

key, err := kmsClient.CreateKey(ctx, region, kms.CreateKmsKeyRequest{
    Name:    "app-secrets",
    KeyType: kms.KmsKeyTypeAES256GCM96,
})

enc, err := kmsClient.EncryptBytes(ctx, region, key.Identity, []byte("hello"))
plain, err := kmsClient.DecryptBytes(ctx, region, key.Identity, enc.Ciphertext)
```
