package objectstorage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/thalassa-cloud/client-go/objectstorage"
)

func lifecycleRuleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Set:      lifecycleRuleHash,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"prefix": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"status": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      string(objectstorage.BucketLifecycleRuleStatusEnabled),
					ValidateFunc: validate.StringInSlice([]string{string(objectstorage.BucketLifecycleRuleStatusEnabled), string(objectstorage.BucketLifecycleRuleStatusDisabled)}, false),
				},
				"filter": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"prefix": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"object_size_greater_than": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"object_size_less_than": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"tag": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"key":   {Type: schema.TypeString, Required: true},
										"value": {Type: schema.TypeString, Required: true},
									},
								},
							},
							"and": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"prefix": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"object_size_greater_than": {
											Type:     schema.TypeInt,
											Optional: true,
										},
										"object_size_less_than": {
											Type:     schema.TypeInt,
											Optional: true,
										},
										"tags": {
											Type:     schema.TypeSet,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"key":   {Type: schema.TypeString, Required: true},
													"value": {Type: schema.TypeString, Required: true},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"expiration": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"days": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"date": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validateRFC3339TimeString,
							},
							"expired_object_delete_marker": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
				"transition": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"days": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"date": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validateRFC3339TimeString,
							},
							"storage_class": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"noncurrent_version_expiration": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"noncurrent_days": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
				"noncurrent_version_transition": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"noncurrent_days": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"storage_class": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"abort_incomplete_multipart_upload": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"days_after_initiation": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func validateRFC3339TimeString(v any, _ string) (warns []string, errs []error) {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil, nil
	}
	if _, err := time.Parse(time.RFC3339, s); err != nil {
		return nil, []error{fmt.Errorf("expected RFC3339 timestamp, got %q: %w", s, err)}
	}
	return nil, nil
}

func lifecycleRuleHash(v any) int {
	m := v.(map[string]any)
	return schema.HashString(m["id"].(string))
}

func expandLifecycleRules(raw *schema.Set) []objectstorage.BucketLifecycleRule {
	if raw == nil || raw.Len() == 0 {
		return []objectstorage.BucketLifecycleRule{}
	}
	rules := make([]objectstorage.BucketLifecycleRule, 0, raw.Len())
	for _, item := range raw.List() {
		block := item.(map[string]any)
		rule := objectstorage.BucketLifecycleRule{
			ID:     block["id"].(string),
			Prefix: block["prefix"].(string),
			Status: objectstorage.BucketLifecycleRuleStatus(block["status"].(string)),
		}
		if v, ok := block["filter"].([]any); ok && len(v) > 0 {
			rule.Filter = expandLifecycleFilter(v[0].(map[string]any))
		}
		if v, ok := block["expiration"].([]any); ok && len(v) > 0 {
			rule.Expiration = expandLifecycleExpiration(v[0].(map[string]any))
		}
		if v, ok := block["transition"].(*schema.Set); ok {
			rule.Transitions = expandLifecycleTransitions(v.List())
		}
		if v, ok := block["noncurrent_version_expiration"].([]any); ok && len(v) > 0 {
			rule.NoncurrentVersionExpiration = expandNoncurrentVersionExpiration(v[0].(map[string]any))
		}
		if v, ok := block["noncurrent_version_transition"].(*schema.Set); ok {
			rule.NoncurrentVersionTransitions = expandNoncurrentVersionTransitions(v.List())
		}
		if v, ok := block["abort_incomplete_multipart_upload"].([]any); ok && len(v) > 0 {
			rule.AbortIncompleteMultipartUpload = expandAbortIncompleteMultipartUpload(v[0].(map[string]any))
		}
		rules = append(rules, rule)
	}
	return rules
}

func expandLifecycleFilter(block map[string]any) *objectstorage.BucketLifecycleRuleFilter {
	filter := &objectstorage.BucketLifecycleRuleFilter{}
	if v, ok := block["prefix"].(string); ok {
		filter.Prefix = v
	}
	if v, ok := block["object_size_greater_than"].(int); ok && v > 0 {
		val := int64(v)
		filter.ObjectSizeGreaterThan = &val
	}
	if v, ok := block["object_size_less_than"].(int); ok && v > 0 {
		val := int64(v)
		filter.ObjectSizeLessThan = &val
	}
	if v, ok := block["tag"].([]any); ok && len(v) > 0 {
		tagBlock := v[0].(map[string]any)
		filter.Tag = &objectstorage.BucketLifecycleRuleTag{
			Key:   tagBlock["key"].(string),
			Value: tagBlock["value"].(string),
		}
	}
	if v, ok := block["and"].([]any); ok && len(v) > 0 {
		andBlock := v[0].(map[string]any)
		and := &objectstorage.BucketLifecycleRuleAndOperator{}
		if p, ok := andBlock["prefix"].(string); ok {
			and.Prefix = p
		}
		if v, ok := andBlock["object_size_greater_than"].(int); ok && v > 0 {
			val := int64(v)
			and.ObjectSizeGreaterThan = &val
		}
		if v, ok := andBlock["object_size_less_than"].(int); ok && v > 0 {
			val := int64(v)
			and.ObjectSizeLessThan = &val
		}
		if tags, ok := andBlock["tags"].(*schema.Set); ok {
			and.Tags = expandLifecycleTags(tags.List())
		}
		filter.And = and
	}
	return filter
}

func expandLifecycleTags(raw []any) []objectstorage.BucketLifecycleRuleTag {
	tags := make([]objectstorage.BucketLifecycleRuleTag, 0, len(raw))
	for _, item := range raw {
		block := item.(map[string]any)
		tags = append(tags, objectstorage.BucketLifecycleRuleTag{
			Key:   block["key"].(string),
			Value: block["value"].(string),
		})
	}
	return tags
}

func expandLifecycleExpiration(block map[string]any) *objectstorage.BucketLifecycleRuleExpiration {
	exp := &objectstorage.BucketLifecycleRuleExpiration{}
	if v, ok := block["days"].(int); ok && v > 0 {
		days := int64(v)
		exp.Days = &days
	}
	if v, ok := block["date"].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			exp.Date = &t
		}
	}
	if v, ok := block["expired_object_delete_marker"].(bool); ok {
		exp.ExpiredObjectDeleteMarker = &v
	}
	return exp
}

func expandLifecycleTransitions(raw []any) []objectstorage.BucketLifecycleRuleTransition {
	transitions := make([]objectstorage.BucketLifecycleRuleTransition, 0, len(raw))
	for _, item := range raw {
		block := item.(map[string]any)
		t := objectstorage.BucketLifecycleRuleTransition{
			StorageClass: block["storage_class"].(string),
		}
		if v, ok := block["days"].(int); ok && v > 0 {
			days := int64(v)
			t.Days = &days
		}
		if v, ok := block["date"].(string); ok && v != "" {
			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
				t.Date = &parsed
			}
		}
		transitions = append(transitions, t)
	}
	return transitions
}

func expandNoncurrentVersionExpiration(block map[string]any) *objectstorage.BucketLifecycleRuleNoncurrentVersionExpiration {
	exp := &objectstorage.BucketLifecycleRuleNoncurrentVersionExpiration{}
	if v, ok := block["noncurrent_days"].(int); ok && v > 0 {
		days := int64(v)
		exp.NoncurrentDays = &days
	}
	return exp
}

func expandNoncurrentVersionTransitions(raw []any) []objectstorage.BucketLifecycleRuleNoncurrentVersionTransition {
	transitions := make([]objectstorage.BucketLifecycleRuleNoncurrentVersionTransition, 0, len(raw))
	for _, item := range raw {
		block := item.(map[string]any)
		t := objectstorage.BucketLifecycleRuleNoncurrentVersionTransition{
			StorageClass: block["storage_class"].(string),
		}
		if v, ok := block["noncurrent_days"].(int); ok && v > 0 {
			days := int64(v)
			t.NoncurrentDays = &days
		}
		transitions = append(transitions, t)
	}
	return transitions
}

func expandAbortIncompleteMultipartUpload(block map[string]any) *objectstorage.BucketLifecycleRuleAbortIncompleteMultipartUpload {
	abort := &objectstorage.BucketLifecycleRuleAbortIncompleteMultipartUpload{}
	if v, ok := block["days_after_initiation"].(int); ok && v > 0 {
		days := int64(v)
		abort.DaysAfterInitiation = &days
	}
	return abort
}

func flattenLifecycleRules(rules []objectstorage.BucketLifecycleRule) []any {
	result := make([]any, 0, len(rules))
	for _, rule := range rules {
		block := map[string]any{
			"id":     rule.ID,
			"prefix": rule.Prefix,
			"status": string(rule.Status),
		}
		if rule.Filter != nil {
			block["filter"] = []any{flattenLifecycleFilter(rule.Filter)}
		}
		if rule.Expiration != nil {
			block["expiration"] = []any{flattenLifecycleExpiration(rule.Expiration)}
		}
		if len(rule.Transitions) > 0 {
			block["transition"] = flattenLifecycleTransitions(rule.Transitions)
		}
		if rule.NoncurrentVersionExpiration != nil {
			block["noncurrent_version_expiration"] = []any{flattenNoncurrentVersionExpiration(rule.NoncurrentVersionExpiration)}
		}
		if len(rule.NoncurrentVersionTransitions) > 0 {
			block["noncurrent_version_transition"] = flattenNoncurrentVersionTransitions(rule.NoncurrentVersionTransitions)
		}
		if rule.AbortIncompleteMultipartUpload != nil {
			block["abort_incomplete_multipart_upload"] = []any{flattenAbortIncompleteMultipartUpload(rule.AbortIncompleteMultipartUpload)}
		}
		result = append(result, block)
	}
	return result
}

func flattenLifecycleFilter(filter *objectstorage.BucketLifecycleRuleFilter) map[string]any {
	block := map[string]any{}
	if filter.Prefix != "" {
		block["prefix"] = filter.Prefix
	}
	if filter.ObjectSizeGreaterThan != nil {
		block["object_size_greater_than"] = int(*filter.ObjectSizeGreaterThan)
	}
	if filter.ObjectSizeLessThan != nil {
		block["object_size_less_than"] = int(*filter.ObjectSizeLessThan)
	}
	if filter.Tag != nil {
		block["tag"] = []any{map[string]any{
			"key":   filter.Tag.Key,
			"value": filter.Tag.Value,
		}}
	}
	if filter.And != nil {
		andBlock := map[string]any{}
		if filter.And.Prefix != "" {
			andBlock["prefix"] = filter.And.Prefix
		}
		if filter.And.ObjectSizeGreaterThan != nil {
			andBlock["object_size_greater_than"] = int(*filter.And.ObjectSizeGreaterThan)
		}
		if filter.And.ObjectSizeLessThan != nil {
			andBlock["object_size_less_than"] = int(*filter.And.ObjectSizeLessThan)
		}
		if len(filter.And.Tags) > 0 {
			tags := make([]any, 0, len(filter.And.Tags))
			for _, tag := range filter.And.Tags {
				tags = append(tags, map[string]any{
					"key":   tag.Key,
					"value": tag.Value,
				})
			}
			andBlock["tags"] = tags
		}
		block["and"] = []any{andBlock}
	}
	return block
}

func flattenLifecycleExpiration(exp *objectstorage.BucketLifecycleRuleExpiration) map[string]any {
	block := map[string]any{}
	if exp.Days != nil {
		block["days"] = int(*exp.Days)
	}
	if exp.Date != nil {
		block["date"] = exp.Date.Format(time.RFC3339)
	}
	if exp.ExpiredObjectDeleteMarker != nil {
		block["expired_object_delete_marker"] = *exp.ExpiredObjectDeleteMarker
	}
	return block
}

func flattenLifecycleTransitions(transitions []objectstorage.BucketLifecycleRuleTransition) []any {
	result := make([]any, 0, len(transitions))
	for _, t := range transitions {
		block := map[string]any{
			"storage_class": t.StorageClass,
		}
		if t.Days != nil {
			block["days"] = int(*t.Days)
		}
		if t.Date != nil {
			block["date"] = t.Date.Format(time.RFC3339)
		}
		result = append(result, block)
	}
	return result
}

func flattenNoncurrentVersionExpiration(exp *objectstorage.BucketLifecycleRuleNoncurrentVersionExpiration) map[string]any {
	block := map[string]any{}
	if exp.NoncurrentDays != nil {
		block["noncurrent_days"] = int(*exp.NoncurrentDays)
	}
	return block
}

func flattenNoncurrentVersionTransitions(transitions []objectstorage.BucketLifecycleRuleNoncurrentVersionTransition) []any {
	result := make([]any, 0, len(transitions))
	for _, t := range transitions {
		block := map[string]any{
			"storage_class": t.StorageClass,
		}
		if t.NoncurrentDays != nil {
			block["noncurrent_days"] = int(*t.NoncurrentDays)
		}
		result = append(result, block)
	}
	return result
}

func flattenAbortIncompleteMultipartUpload(abort *objectstorage.BucketLifecycleRuleAbortIncompleteMultipartUpload) map[string]any {
	block := map[string]any{}
	if abort.DaysAfterInitiation != nil {
		block["days_after_initiation"] = int(*abort.DaysAfterInitiation)
	}
	return block
}

func lifecycleHasNoncurrentRules(rules []objectstorage.BucketLifecycleRule) bool {
	for _, rule := range rules {
		if rule.NoncurrentVersionExpiration != nil || len(rule.NoncurrentVersionTransitions) > 0 {
			return true
		}
	}
	return false
}
