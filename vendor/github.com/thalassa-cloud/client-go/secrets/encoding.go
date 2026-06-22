package secrets

import (
	"encoding/base64"
	"fmt"
)

// EncodeBytes encodes bytes as standard base64 for secret payloads.
func EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// DecodeBytes decodes a standard base64 secret payload field.
func DecodeBytes(field, encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("%s must be valid base64: %w", field, err)
	}
	return data, nil
}
