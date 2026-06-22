package kms

import (
	"context"
	"strconv"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// Encrypt encrypts plaintext using a KMS key. Plaintext must be base64-encoded.
func (c *Client) Encrypt(ctx context.Context, region, identity string, encrypt EncryptRequest) (*EncryptResponse, error) {
	var result EncryptResponse
	req := c.R().SetBody(encrypt).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "encrypt"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// EncryptBytes encrypts raw bytes using a KMS key.
func (c *Client) EncryptBytes(ctx context.Context, region, identity string, plaintext []byte) (*EncryptResponse, error) {
	return c.Encrypt(ctx, region, identity, EncryptRequest{
		Plaintext: EncodeBytes(plaintext),
	})
}

// Decrypt decrypts ciphertext using a KMS key.
// Callers must not log or persist decrypted material from the response.
func (c *Client) Decrypt(ctx context.Context, region, identity string, decrypt DecryptRequest) (*DecryptResponse, error) {
	var result DecryptResponse
	req := c.R().SetBody(decrypt).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "decrypt"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// DecryptBytes decrypts ciphertext and returns the decoded plaintext bytes.
func (c *Client) DecryptBytes(ctx context.Context, region, identity string, ciphertext string) ([]byte, error) {
	result, err := c.Decrypt(ctx, region, identity, DecryptRequest{Ciphertext: ciphertext})
	if err != nil {
		return nil, err
	}
	return DecodeBytes("plaintext", result.Plaintext)
}

// Sign signs a message using an asymmetric KMS key.
func (c *Client) Sign(ctx context.Context, region, identity string, sign SignRequest) (*SignResponse, error) {
	var result SignResponse
	req := c.R().SetBody(sign).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "sign"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// VerifySignature verifies a signature using an asymmetric KMS key.
func (c *Client) VerifySignature(ctx context.Context, region, identity string, verify VerifySignatureRequest) (*VerifySignatureResponse, error) {
	var result VerifySignatureResponse
	req := c.R().SetBody(verify).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "verify"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// HMAC computes an HMAC using an HMAC KMS key.
func (c *Client) HMAC(ctx context.Context, region, identity string, hmacReq HMACRequest) (*HMACResponse, error) {
	var result HMACResponse
	req := c.R().SetBody(hmacReq).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "hmac"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// VerifyHMAC verifies an HMAC using an HMAC KMS key.
func (c *Client) VerifyHMAC(ctx context.Context, region, identity string, verify VerifyHMACRequest) (*VerifyHMACResponse, error) {
	var result VerifyHMACResponse
	req := c.R().SetBody(verify).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "verify-hmac"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// GetPublicKey returns the public key material for an asymmetric KMS key.
func (c *Client) GetPublicKey(ctx context.Context, region, identity string, version *int) (*GetPublicKeyResponse, error) {
	var result GetPublicKeyResponse
	req := c.R().SetResult(&result)
	if version != nil {
		req = req.SetQueryParam("version", strconv.Itoa(*version))
	}
	resp, err := c.Do(ctx, req, client.GET, regionPath(region, "keys", identity, "public-key"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// ExportKey exports key material when export is allowed on the key.
// Callers must not log or persist exported key material from the response.
func (c *Client) ExportKey(ctx context.Context, region, identity string, export ExportKeyRequest) (*ExportKeyResponse, error) {
	var result ExportKeyResponse
	req := c.R().SetBody(export).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "export"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}
