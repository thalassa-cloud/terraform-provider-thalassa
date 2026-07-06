package kms

import (
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type KmsKeyStatus string

const (
	KmsKeyStatusActive          KmsKeyStatus = "active"
	KmsKeyStatusDisabled        KmsKeyStatus = "disabled"
	KmsKeyStatusPendingDeletion KmsKeyStatus = "pending_deletion"
)

type KmsKeyType string

const (
	KmsKeyTypeAES128GCM96      KmsKeyType = "aes128-gcm96"
	KmsKeyTypeAES256GCM96      KmsKeyType = "aes256-gcm96"
	KmsKeyTypeChaCha20Poly1305 KmsKeyType = "chacha20-poly1305"
	KmsKeyTypeEd25519          KmsKeyType = "ed25519"
	KmsKeyTypeECDSAP256        KmsKeyType = "ecdsa-p256"
	KmsKeyTypeECDSAP384        KmsKeyType = "ecdsa-p384"
	KmsKeyTypeECDSAP521        KmsKeyType = "ecdsa-p521"
	KmsKeyTypeRSA2048          KmsKeyType = "rsa-2048"
	KmsKeyTypeRSA3072          KmsKeyType = "rsa-3072"
	KmsKeyTypeRSA4096          KmsKeyType = "rsa-4096"
	KmsKeyTypeHMAC             KmsKeyType = "hmac"
	KmsKeyTypeHMACSHA256       KmsKeyType = "hmac-sha256"
	KmsKeyTypeHMACSHA512       KmsKeyType = "hmac-sha512"
)

// IsValid reports whether the key type is supported.
func (t KmsKeyType) IsValid() bool {
	switch t {
	case KmsKeyTypeAES128GCM96,
		KmsKeyTypeAES256GCM96,
		KmsKeyTypeChaCha20Poly1305,
		KmsKeyTypeEd25519,
		KmsKeyTypeECDSAP256,
		KmsKeyTypeECDSAP384,
		KmsKeyTypeECDSAP521,
		KmsKeyTypeRSA2048,
		KmsKeyTypeRSA3072,
		KmsKeyTypeRSA4096,
		KmsKeyTypeHMAC,
		KmsKeyTypeHMACSHA256,
		KmsKeyTypeHMACSHA512:
		return true
	default:
		return false
	}
}

type KmsKeyVersion struct {
	Version   int       `json:"version"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type KmsKey struct {
	Identity             string             `json:"identity"`
	Name                 string             `json:"name"`
	Slug                 string             `json:"slug"`
	Description          string             `json:"description,omitempty"`
	Labels               map[string]string  `json:"labels,omitempty"`
	Annotations          map[string]string  `json:"annotations,omitempty"`
	ProjectID            string             `json:"projectId,omitempty"`
	KeyType              KmsKeyType         `json:"keyType"`
	Status               KmsKeyStatus       `json:"status"`
	ExportAllowed        bool               `json:"exportAllowed"`
	Imported             bool               `json:"imported,omitempty"`
	KeyRotationEnabled   bool               `json:"keyRotationEnabled"`
	RotationPeriodInDays *int               `json:"rotationPeriodInDays,omitempty"`
	LatestVersion        int                `json:"latestVersion,omitempty"`
	Versions             []KmsKeyVersion    `json:"versions,omitempty"`
	CreatedAt            time.Time          `json:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt"`
	ObjectVersion        int                `json:"objectVersion,omitempty"`
	Organisation         *base.Organisation `json:"organisation,omitempty"`
}

type KmsSummaryRegion struct {
	Identity            string `json:"identity"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	KmsAvailable        bool   `json:"kmsAvailable"`
	TotalKeys           int64  `json:"totalKeys"`
	ActiveKeys          int64  `json:"activeKeys"`
	DisabledKeys        int64  `json:"disabledKeys"`
	PendingDeletionKeys int64  `json:"pendingDeletionKeys"`
}

type KmsSummary struct {
	FeatureEnabled bool               `json:"featureEnabled"`
	Regions        []KmsSummaryRegion `json:"regions,omitempty"`
}

type ListKeysRequest struct {
	Filters []filters.Filter
}

type CreateKmsKeyRequest struct {
	Name                 string            `json:"name"`
	Description          string            `json:"description,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	Annotations          map[string]string `json:"annotations,omitempty"`
	KeyType              KmsKeyType        `json:"keyType,omitempty"`
	ExportAllowed        bool              `json:"exportAllowed,omitempty"`
	KeyRotationEnabled   bool              `json:"keyRotationEnabled,omitempty"`
	RotationPeriodInDays *int              `json:"rotationPeriodInDays,omitempty"`
	ImportKeyMaterial    string            `json:"importKeyMaterial,omitempty"`
	HashFunction         string            `json:"hashFunction,omitempty"`
	AllowRotation        bool              `json:"allowRotation,omitempty"`
}

type UpdateRotationRequest struct {
	KeyRotationEnabled   *bool `json:"keyRotationEnabled,omitempty"`
	RotationPeriodInDays *int  `json:"rotationPeriodInDays,omitempty"`
}

type EncryptRequest struct {
	Plaintext  string `json:"plaintext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type EncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type DecryptRequest struct {
	Ciphertext string `json:"ciphertext"`
}

type DecryptResponse struct {
	Plaintext  string `json:"plaintext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type SignRequest struct {
	Input         string `json:"input" validate:"required"`
	KeyVersion    string `json:"keyVersion,omitempty"`
	HashAlgorithm string `json:"hashAlgorithm,omitempty"`
	Prehashed     bool   `json:"prehashed,omitempty"`
	Context       string `json:"context,omitempty"`
}

type SignResponse struct {
	Signature  string `json:"signature"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type VerifySignatureRequest struct {
	Input         string `json:"input" validate:"required"`
	Signature     string `json:"signature" validate:"required"`
	HashAlgorithm string `json:"hashAlgorithm,omitempty"`
}

type VerifySignatureResponse struct {
	Valid bool `json:"valid"`
}

type HMACRequest struct {
	Input      string `json:"input" validate:"required"`
	KeyVersion string `json:"keyVersion,omitempty"`
	Algorithm  string `json:"algorithm,omitempty"`
}

type HMACResponse struct {
	HMAC       string `json:"hmac"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type VerifyHMACRequest struct {
	Input         string `json:"input" validate:"required"`
	HMAC          string `json:"hmac" validate:"required"`
	HashAlgorithm string `json:"hashAlgorithm,omitempty"`
}

type VerifyHMACResponse struct {
	Valid      bool   `json:"valid"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type GetPublicKeyResponse struct {
	Keys map[string]string `json:"keys"`
}

type ExportKeyRequest struct {
	KeyVersion string `json:"keyVersion,omitempty"`
}

type ExportKeyResponse struct {
	KeyMaterial string `json:"keyMaterial"`
	KeyVersion  string `json:"keyVersion,omitempty"`
}

type WrappingKeyResponse struct {
	PublicKey string `json:"publicKey"`
	Algorithm string `json:"algorithm,omitempty"`
}
