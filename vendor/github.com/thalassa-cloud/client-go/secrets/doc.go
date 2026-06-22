// Package secrets provides a client for the Thalassa Cloud Secrets Manager.
//
// Secrets are path-based resources (for example "/app/prod/db/password") encrypted with
// KMS. The SDK embeds the path directly in request URLs:
//
//	GET  /v1/secrets/nl-01/secret/app/prod/db/password
//	POST /v1/secrets/nl-01/secret/app/prod/db/password/value
//
// There is no slash between "secret" and the path because the path already starts with "/".
//
// # Client setup
//
// Use the Thalassa facade client and call Secrets():
//
//	tc, err := thalassa.NewClient(
//	    client.WithBaseURL("https://api.thalassa.cloud"),
//	    client.WithOrganisation("my-org"),
//	    client.WithAuthPersonalToken(os.Getenv("THALASSA_TOKEN")),
//	)
//	secretsClient := tc.Secrets()
//
// Or construct a Secrets client directly from a base client:
//
//	base, err := client.NewClient(/* options */)
//	secretsClient, err := secrets.New(base)
//
// # Create and read secrets
//
// Exactly one of SecretString, SecretKeyValues, or GenerateSecret must be supplied on
// create. String values are base64-encoded on the wire:
//
//	_, err := secretsClient.CreateSecret(ctx, "nl-01", secrets.CreateSecretRequest{
//	    Path:           "/app/prod/db/password",
//	    KmsKeyIdentity: "kms-abc123",
//	    SecretString:   secrets.EncodeBytes([]byte("super-secret")),
//	})
//	val, version, err := secretsClient.GetSecretString(ctx, "nl-01", "/app/prod/db/password", nil)
//
// See the package README and runnable examples in example_test.go for full usage.
package secrets
