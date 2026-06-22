package kms

import "strings"

const KmsEndpoint = "/v1/kms"

func regionPath(region string, segments ...string) string {
	parts := append([]string{KmsEndpoint, region}, segments...)
	return strings.Join(parts, "/")
}
