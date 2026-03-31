package iaas

import (
	"context"
	"fmt"
	"strings"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

// lookupMachineImageIdentity resolves a machine image reference to the API identity.
// The reference may be the image identity, slug, or name (name match is case-insensitive).
func lookupMachineImageIdentity(ctx context.Context, client *iaas.Client, ref string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("machine image is required")
	}
	images, err := client.ListMachineImages(ctx, &iaas.ListMachineImagesRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to list machine images: %w", err)
	}
	for _, img := range images {
		if img.Identity == ref || img.Slug == ref || strings.EqualFold(img.Name, ref) {
			return img.Identity, nil
		}
	}
	names := make([]string, len(images))
	for i, img := range images {
		names[i] = img.Name
	}
	return "", fmt.Errorf("machine image not found: %s. Available image names: %v", ref, strings.Join(names, ", "))
}
