package kms

import (
	"time"

	tckms "github.com/thalassa-cloud/client-go/kms"
)

const timeFormatRFC3339 = time.RFC3339

var kmsKeyTypes = []string{
	string(tckms.KmsKeyTypeAES128GCM96),
	string(tckms.KmsKeyTypeAES256GCM96),
	string(tckms.KmsKeyTypeChaCha20Poly1305),
	string(tckms.KmsKeyTypeEd25519),
	string(tckms.KmsKeyTypeECDSAP256),
	string(tckms.KmsKeyTypeECDSAP384),
	string(tckms.KmsKeyTypeECDSAP521),
	string(tckms.KmsKeyTypeRSA2048),
	string(tckms.KmsKeyTypeRSA3072),
	string(tckms.KmsKeyTypeRSA4096),
	string(tckms.KmsKeyTypeHMAC),
	string(tckms.KmsKeyTypeHMACSHA256),
	string(tckms.KmsKeyTypeHMACSHA512),
}

var kmsKeyStatuses = []string{
	string(tckms.KmsKeyStatusActive),
	string(tckms.KmsKeyStatusDisabled),
}

func setKmsKeyState(d interface {
	Set(string, any) error
}, key *tckms.KmsKey, region string) error {
	_ = d.Set("region", region)
	_ = d.Set("name", key.Name)
	_ = d.Set("slug", key.Slug)
	_ = d.Set("description", key.Description)
	_ = d.Set("labels", key.Labels)
	_ = d.Set("annotations", key.Annotations)
	_ = d.Set("key_type", string(key.KeyType))
	_ = d.Set("status", string(key.Status))
	_ = d.Set("export_allowed", key.ExportAllowed)
	_ = d.Set("imported", key.Imported)
	_ = d.Set("key_rotation_enabled", key.KeyRotationEnabled)
	_ = d.Set("rotation_period_in_days", key.RotationPeriodInDays)
	_ = d.Set("latest_version", key.LatestVersion)
	_ = d.Set("object_version", key.ObjectVersion)
	if !key.CreatedAt.IsZero() {
		_ = d.Set("created_at", key.CreatedAt.Format(timeFormatRFC3339))
	}
	if !key.UpdatedAt.IsZero() {
		_ = d.Set("updated_at", key.UpdatedAt.Format(timeFormatRFC3339))
	}
	return nil
}

func regionKmsAvailable(summary *tckms.KmsSummary, region string) bool {
	if summary == nil || !summary.FeatureEnabled {
		return false
	}
	for _, r := range summary.Regions {
		if r.Slug == region || r.Identity == region {
			return r.KmsAvailable
		}
	}
	return false
}
