package objectstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBucketLifecycleResponse(t *testing.T) {
	t.Parallel()

	t.Run("flat rules", func(t *testing.T) {
		t.Parallel()

		lifecycle, err := decodeBucketLifecycleResponse([]byte(`{
			"rules": [{
				"id": "expire-logs",
				"prefix": "logs/",
				"status": "Enabled",
				"expiration": {"days": 30}
			}]
		}`))
		assert.NoError(t, err)
		assert.Len(t, lifecycle.Rules, 1)
		assert.Equal(t, "expire-logs", lifecycle.Rules[0].ID)
	})

	t.Run("wrapped lifecycle", func(t *testing.T) {
		t.Parallel()

		lifecycle, err := decodeBucketLifecycleResponse([]byte(`{
			"lifecycle": {
				"rules": [{
					"id": "expire-logs",
					"prefix": "logs/",
					"status": "Enabled",
					"expiration": {"days": 30}
				}]
			}
		}`))
		assert.NoError(t, err)
		assert.Len(t, lifecycle.Rules, 1)
		assert.Equal(t, "logs/", lifecycle.Rules[0].Prefix)
	})

	t.Run("empty body", func(t *testing.T) {
		t.Parallel()

		lifecycle, err := decodeBucketLifecycleResponse(nil)
		assert.NoError(t, err)
		assert.Empty(t, lifecycle.Rules)
	})
}

func TestDecodeBucketLifecycleResponseNoncurrent(t *testing.T) {
	t.Parallel()

	lifecycle, err := decodeBucketLifecycleResponse([]byte(`{
		"lifecycle": {
			"rules": [{
				"id": "expire-noncurrent",
				"prefix": "archive/",
				"status": "Enabled",
				"noncurrentVersionExpiration": {"noncurrentDays": 7}
			}]
		}
	}`))
	assert.NoError(t, err)
	assert.Len(t, lifecycle.Rules, 1)
	assert.NotNil(t, lifecycle.Rules[0].NoncurrentVersionExpiration)
	assert.Equal(t, int64(7), *lifecycle.Rules[0].NoncurrentVersionExpiration.NoncurrentDays)
}
