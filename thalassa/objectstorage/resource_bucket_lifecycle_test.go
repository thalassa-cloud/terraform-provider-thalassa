package objectstorage

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/thalassa-cloud/client-go/objectstorage"
)

func TestResourceBucketLifecycle(t *testing.T) {
	resource := resourceBucketLifecycle()
	assert.True(t, resource.Schema["bucket_name"].ForceNew)
	assert.Equal(t, schema.TypeList, resource.Schema["rule"].Type)
}

func TestExpandFlattenLifecycleRules(t *testing.T) {
	raw := map[string]any{
		"id":     "expire-logs",
		"prefix": "logs/",
		"status": string(objectstorage.BucketLifecycleRuleStatusEnabled),
		"expiration": []any{
			map[string]any{"days": 30},
		},
	}
	rules := expandLifecycleRules([]any{raw})
	assert.Len(t, rules, 1)
	assert.Equal(t, "expire-logs", rules[0].ID)
	assert.NotNil(t, rules[0].Expiration.Days)
	assert.Equal(t, int64(30), *rules[0].Expiration.Days)

	flat := flattenLifecycleRules(rules)
	assert.Len(t, flat, 1)

	filterOnly := flattenLifecycleRules([]objectstorage.BucketLifecycleRule{
		{
			ID:     "expire-logs",
			Status: objectstorage.BucketLifecycleRuleStatusEnabled,
			Filter: &objectstorage.BucketLifecycleRuleFilter{Prefix: "logs/"},
			Expiration: &objectstorage.BucketLifecycleRuleExpiration{
				Days: ptrInt64(30),
			},
		},
	})
	assert.Len(t, filterOnly, 1)
	block := filterOnly[0].(map[string]any)
	assert.Equal(t, "logs/", block["prefix"])
	assert.Nil(t, block["filter"])

	zeroSizeFilter := flattenLifecycleRules([]objectstorage.BucketLifecycleRule{
		{
			ID:     "expire-noncurrent",
			Prefix: "archive/",
			Status: objectstorage.BucketLifecycleRuleStatusEnabled,
			Filter: &objectstorage.BucketLifecycleRuleFilter{
				Prefix:                "archive/",
				ObjectSizeGreaterThan: ptrInt64(0),
				ObjectSizeLessThan:    ptrInt64(0),
			},
			NoncurrentVersionExpiration: &objectstorage.BucketLifecycleRuleNoncurrentVersionExpiration{
				NoncurrentDays: ptrInt64(7),
			},
		},
	})
	assert.Len(t, zeroSizeFilter, 1)
	noncurrentBlock := zeroSizeFilter[0].(map[string]any)
	assert.Equal(t, "archive/", noncurrentBlock["prefix"])
	assert.Nil(t, noncurrentBlock["filter"])
}

func ptrInt64(v int64) *int64 {
	return &v
}

func TestLifecycleHasNoncurrentRules(t *testing.T) {
	assert.False(t, lifecycleHasNoncurrentRules([]objectstorage.BucketLifecycleRule{
		{ID: "a"},
	}))
	assert.True(t, lifecycleHasNoncurrentRules([]objectstorage.BucketLifecycleRule{
		{ID: "a", NoncurrentVersionExpiration: &objectstorage.BucketLifecycleRuleNoncurrentVersionExpiration{}},
	}))
}
