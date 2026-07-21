package objectstorage

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/objectstorage"
)

var thalassaPrincipalARNPattern = regexp.MustCompile(`^arn:thalassa:iam:::(?:serviceaccount|user)/[^:]+:(?:[^:]+|\*)$`)

func parseBucketPolicyJSON(raw string) (*objectstorage.PolicyDocument, error) {
	if raw == "" {
		return nil, nil
	}

	var doc objectstorage.PolicyDocument
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, fmt.Errorf("invalid policy JSON: %w", err)
	}

	if err := validateBucketPolicyDocument(doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func validateBucketPolicyDocument(doc objectstorage.PolicyDocument) error {
	for i, statement := range doc.Statement {
		for _, principalARN := range thalassaPrincipalARNs(statement.Principal.Thalassa) {
			if err := validateThalassaPrincipalARN(principalARN); err != nil {
				return fmt.Errorf("policy statement %d: %w", i, err)
			}
		}
	}

	return nil
}

func thalassaPrincipalARNs(raw any) []string {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil
		}
		return []string{v}
	case []string:
		return v
	case []any:
		arns := make([]string, 0, len(v))
		for _, item := range v {
			if arn, ok := item.(string); ok && arn != "" {
				arns = append(arns, arn)
			}
		}
		return arns
	default:
		return nil
	}
}

var (
	expectedPrincipalARNs = []string{
		"*",
		"arn:thalassa:iam:::serviceaccount/<organisation-id>:<service-account-id>",
		"arn:thalassa:iam:::serviceaccount/<organisation-id>/*",
		"arn:thalassa:iam:::user/<organisation-id>:<user-id>",
		"arn:thalassa:iam:::user/<organisation-id>/*",
	}
)

func validateThalassaPrincipalARN(arn string) error {
	arn = strings.TrimSpace(arn)
	if arn == "" {
		return fmt.Errorf("Principal.Thalassa ARN cannot be empty")
	}

	if arn == "*" || thalassaPrincipalARNPattern.MatchString(arn) {
		return nil
	}

	return fmt.Errorf(
		"invalid Principal.Thalassa ARN %q: expected one of %s",
		arn,
		strings.Join(expectedPrincipalARNs, ", "),
	)
}

func enrichBucketError(err error, action string) error {
	if err == nil {
		return nil
	}

	if strings.Contains(strings.ToLower(err.Error()), "invalid principal arn") {
		return fmt.Errorf(
			"failed to %s bucket: the policy contains an invalid Principal.Thalassa ARN. "+
				"Use \"*\", a service account principal such as "+
				"arn:thalassa:iam:::serviceaccount/<organisation-id>:<service-account-id>, "+
				"arn:thalassa:iam:::serviceaccount/<organisation-id>/*, "+
				"or a user principal such as arn:thalassa:iam:::user/<organisation-id>:<user-id> or "+
				"arn:thalassa:iam:::user/<organisation-id>/*: %w",
			action,
			err,
		)
	}

	return fmt.Errorf("failed to %s bucket: %w", action, err)
}

func suppressEquivalentPolicy(_, old, new string, _ *schema.ResourceData) bool {
	return equivalentPolicyJSON(old, new)
}

func equivalentPolicyJSON(a, b string) bool {
	if a == b {
		return true
	}

	aCanon, errA := canonicalPolicyJSONString(a)
	bCanon, errB := canonicalPolicyJSONString(b)
	if errA != nil || errB != nil {
		return false
	}

	return aCanon == bCanon
}

func canonicalPolicyJSONString(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}

	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return "", err
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func bucketPolicyStateValue(configured, apiPolicy string) string {
	if configured != "" && equivalentPolicyJSON(configured, apiPolicy) {
		return configured
	}

	return apiPolicy
}
