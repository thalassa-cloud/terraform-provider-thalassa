package iaas

import (
	"context"
	"fmt"
	"strings"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func lookupVolumeType(ctx context.Context, client *iaas.Client, volumeType string) (*iaas.VolumeType, error) {
	// lets find the volume type
	volumeTypes, err := client.ListVolumeTypes(ctx, &iaas.ListVolumeTypesRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volume types: %w", err)
	}
	var foundVolumeType *iaas.VolumeType
	for _, v := range volumeTypes {
		if v.Identity == volumeType || strings.EqualFold(v.Name, volumeType) {
			foundVolumeType = &v
			break
		}
	}

	if foundVolumeType == nil {
		availableVolumeTypes := make([]string, len(volumeTypes))
		for i, v := range volumeTypes {
			availableVolumeTypes[i] = v.Name
		}
		return nil, fmt.Errorf("volume type not found: %s. Available volume types: %v", volumeType, strings.Join(availableVolumeTypes, ", "))
	}
	return foundVolumeType, nil
}
