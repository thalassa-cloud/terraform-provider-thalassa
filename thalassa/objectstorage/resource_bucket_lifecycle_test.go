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
	assert.Equal(t, schema.TypeSet, resource.Schema["rule"].Type)
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
	set := schema.NewSet(lifecycleRuleHash, []any{raw})
	rules := expandLifecycleRules(set)
	assert.Len(t, rules, 1)
	assert.Equal(t, "expire-logs", rules[0].ID)
	assert.NotNil(t, rules[0].Expiration.Days)
	assert.Equal(t, int64(30), *rules[0].Expiration.Days)

	flat := flattenLifecycleRules(rules)
	assert.Len(t, flat, 1)
}

func TestLifecycleHasNoncurrentRules(t *testing.T) {
	assert.False(t, lifecycleHasNoncurrentRules([]objectstorage.BucketLifecycleRule{
		{ID: "a"},
	}))
	assert.True(t, lifecycleHasNoncurrentRules([]objectstorage.BucketLifecycleRule{
		{ID: "a", NoncurrentVersionExpiration: &objectstorage.BucketLifecycleRuleNoncurrentVersionExpiration{}},
	}))
}
