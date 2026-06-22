// Package kms provides a client for the Thalassa Cloud Key Management Service (KMS).
//
// Use it to create and manage cryptographic keys, encrypt and decrypt data, sign and
// verify messages, and control key lifecycle. All regional operations take an explicit
// region argument (identity or slug, for example "nl-01").
//
// # Client setup
//
// Use the Thalassa facade client and call KMS():
//
//	tc, err := thalassa.NewClient(
//	    client.WithBaseURL("https://api.thalassa.cloud"),
//	    client.WithOrganisation("my-org"),
//	    client.WithAuthPersonalToken(os.Getenv("THALASSA_TOKEN")),
//	)
//	kmsClient := tc.KMS()
//
// Or construct a KMS client directly from a base client:
//
//	base, err := client.NewClient(/* options */)
//	kmsClient, err := kms.New(base)
//
// # Availability
//
// Call GetSummary before using regional endpoints to confirm KMS is enabled and
// which regions expose a KMS endpoint:
//
//	summary, err := kmsClient.GetSummary(ctx)
//	if !summary.FeatureEnabled {
//	    return fmt.Errorf("KMS not enabled")
//	}
//
// # Encrypt and decrypt
//
// Crypto payloads on the wire are base64-encoded. Prefer the byte helpers:
//
//	enc, err := kmsClient.EncryptBytes(ctx, "nl-01", key.Identity, []byte("hello"))
//	plain, err := kmsClient.DecryptBytes(ctx, "nl-01", key.Identity, enc.Ciphertext)
//
// See the package README and runnable examples in example_test.go for full usage.
package kms
