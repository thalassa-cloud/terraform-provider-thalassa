package secrets

import (
	"fmt"
	"regexp"
	"strings"
)

const SecretsEndpoint = "/v1/secrets"

var secretPathPattern = regexp.MustCompile(`^/[a-zA-Z0-9/._+\-]*$`)

// NormalizePath ensures the path has a leading slash and contains only allowed characters.
func NormalizePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !secretPathPattern.MatchString(path) {
		return "", fmt.Errorf("invalid secret path %q", path)
	}
	return path, nil
}

// SecretResourceURL builds /v1/secrets/{region}/secret{path}[/suffix].
// There is no slash between "secret" and the path because the path already starts with "/".
func SecretResourceURL(region, path, suffix string) (string, error) {
	path, err := NormalizePath(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/secret%s%s", SecretsEndpoint, region, path, suffix), nil
}

func secretsCollectionURL(region string) string {
	return fmt.Sprintf("%s/%s/secrets", SecretsEndpoint, region)
}
