package secrets

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tcsecrets "github.com/thalassa-cloud/client-go/secrets"
)

const timeFormatRFC3339 = time.RFC3339

func secretID(region, path string) string {
	return region + path
}

func parseSecretID(id string) (region, path string, err error) {
	slashIdx := strings.Index(id, "/")
	if slashIdx <= 0 {
		return "", "", fmt.Errorf("invalid secret import ID %q, expected {region}{path} e.g. nl-01/app/prod/db/password", id)
	}
	region = id[:slashIdx]
	path = "/" + id[slashIdx+1:]
	return region, path, nil
}

func secretVersionID(region, path string, version int) string {
	return fmt.Sprintf("%s%s/%d", region, path, version)
}

func parseSecretVersionID(id string) (region, path string, version int, err error) {
	lastSlash := strings.LastIndex(id, "/")
	if lastSlash <= 0 {
		return "", "", 0, fmt.Errorf("invalid secret version import ID %q, expected {region}{path}/{version}", id)
	}
	versionStr := id[lastSlash+1:]
	version, err = strconv.Atoi(versionStr)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid secret version in import ID %q: %w", id, err)
	}
	region, path, err = parseSecretID(id[:lastSlash])
	return region, path, version, err
}

func validateSecretPath(path interface{}, _ string) (warns []string, errs []error) {
	p, ok := path.(string)
	if !ok {
		return nil, []error{fmt.Errorf("path must be a string")}
	}
	if !strings.HasPrefix(p, "/") {
		return nil, []error{fmt.Errorf("path must start with /, got %q", p)}
	}
	return nil, nil
}

func setSecretState(d interface {
	Set(string, interface{}) error
}, secret *tcsecrets.Secret, region string) error {
	if err := d.Set("region", region); err != nil {
		return err
	}
	if err := d.Set("path", secret.Path); err != nil {
		return err
	}
	if err := d.Set("description", secret.Description); err != nil {
		return err
	}
	if err := d.Set("labels", secret.Labels); err != nil {
		return err
	}
	if err := d.Set("annotations", secret.Annotations); err != nil {
		return err
	}
	if err := d.Set("current_version", secret.CurrentVersion); err != nil {
		return err
	}
	if !secret.CreatedAt.IsZero() {
		if err := d.Set("created_at", secret.CreatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	if !secret.UpdatedAt.IsZero() {
		if err := d.Set("updated_at", secret.UpdatedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	if secret.LastAccessedAt != nil {
		if err := d.Set("last_accessed_at", secret.LastAccessedAt.Format(timeFormatRFC3339)); err != nil {
			return err
		}
	}
	return nil
}

func expandGenerateSecret(raw []interface{}) *tcsecrets.GenerateSecret {
	if len(raw) == 0 {
		return nil
	}
	block := raw[0].(map[string]interface{})
	gen := &tcsecrets.GenerateSecret{}
	if v, ok := block["byte_length"].(int); ok && v > 0 {
		gen.ByteLength = v
	}
	return gen
}

func expandAccessPolicyStatements(raw []interface{}) []tcsecrets.SecretPolicyStatement {
	if len(raw) == 0 {
		return nil
	}
	statements := make([]tcsecrets.SecretPolicyStatement, 0, len(raw))
	for _, item := range raw {
		block := item.(map[string]interface{})
		stmt := tcsecrets.SecretPolicyStatement{
			Effect: block["effect"].(string),
		}
		if v, ok := block["actions"].([]interface{}); ok {
			stmt.Actions = make([]string, len(v))
			for i, a := range v {
				stmt.Actions[i] = a.(string)
			}
		}
		if v, ok := block["principals"].([]interface{}); ok {
			stmt.Principals = make([]string, len(v))
			for i, p := range v {
				stmt.Principals[i] = p.(string)
			}
		}
		statements = append(statements, stmt)
	}
	return statements
}

func flattenAccessPolicyStatements(statements []tcsecrets.SecretPolicyStatement) []map[string]interface{} {
	if len(statements) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(statements))
	for _, stmt := range statements {
		actions := make([]interface{}, len(stmt.Actions))
		for i, a := range stmt.Actions {
			actions[i] = a
		}
		principals := make([]interface{}, len(stmt.Principals))
		for i, p := range stmt.Principals {
			principals[i] = p
		}
		result = append(result, map[string]interface{}{
			"effect":     stmt.Effect,
			"actions":    actions,
			"principals": principals,
		})
	}
	return result
}
