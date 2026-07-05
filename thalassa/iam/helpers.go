package iam

import (
	"time"

	iam "github.com/thalassa-cloud/client-go/iam"
)

const (
	TimeFormatRFC3339 = time.RFC3339
)

func convertPermissionsToStrings(permissions []iam.PermissionType) []string {
	result := make([]string, len(permissions))
	for i, perm := range permissions {
		result[i] = string(perm)
	}
	return result
}

func convertStringsToPermissions(permissions []string) []iam.PermissionType {
	result := make([]iam.PermissionType, len(permissions))
	for i, perm := range permissions {
		result[i] = iam.PermissionType(perm)
	}
	return result
}

func valueOrEmptySlice(value []string) []string {
	if value == nil {
		return []string{}
	}
	return value
}

func toListOfInterfaces(value []string) []any {
	if value == nil {
		return []any{}
	}
	result := make([]any, len(value))
	for i, v := range value {
		result[i] = v
	}
	return result
}
