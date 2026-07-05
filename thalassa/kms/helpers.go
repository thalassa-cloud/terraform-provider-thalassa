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
	if err := d.Set("region", region); err != nil {
		return err
	}
	if err := d.Set("name", key.Name); err != nil {
		return err
	}
	if err := d.Set("slug", key.Slug); err != nil {
		return err
	}
	if err := d.Set("description", key.Description); err != nil {
		return err
	}
	if err := d.Set("labels", key.Labels); err != nil {
		return err
	}
	if err := d.Set("annotations", key.Annotations); err != nil {
		return err
	}
	if err := d.Set("key_type", string(key.KeyType)); err != nil {
		return err
	}
	if err := d.Set("status", string(key.Status)); err != nil {
		return err
	}
	if err := d.Set("export_allowed", key.ExportAllowed); err != nil {
		return err
	}
	if err := d.Set("imported", key.Imported); err != nil {
		return err
	}
	if err := d.Set("key_rotation_enabled", key.KeyRotationEnabled); err != nil {
		return err
	}
	if err := d.Set("rotation_period_in_days", key.RotationPeriodInDays); err != nil {
		return err
	}
	if err := d.Set("latest_version", key.LatestVersion); err != nil {
		return err
	}
	if err := d.Set("object_version", key.ObjectVersion); err != nil {
		return err
	}
	if !key.CreatedAt.IsZero() {
		if err := d.Set("created_at", key.CreatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	if !key.UpdatedAt.IsZero() {
		if err := d.Set("updated_at", key.UpdatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
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
