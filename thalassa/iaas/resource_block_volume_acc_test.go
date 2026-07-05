package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockVolume_basic(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockVolumeConfig(volumeName, region, volumeType, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_block_volume.test", "name", volumeName),
					resource.TestCheckResourceAttr("thalassa_block_volume.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_block_volume.test", "volume_type", volumeType),
					resource.TestCheckResourceAttr("thalassa_block_volume.test", "size_gb", "10"),
					resource.TestCheckResourceAttrSet("thalassa_block_volume.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_block_volume.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_block_volume.test", "status"),
				),
			},
		},
	})
}

func TestAccBlockVolume_update(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockVolumeConfigWithDescription(volumeName, region, volumeType, 10, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_block_volume.test", "description", "initial description"),
				),
			},
			{
				Config: testAccBlockVolumeConfigWithDescription(volumeName, region, volumeType, 10, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_block_volume.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccBlockVolume_import(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockVolumeConfig(volumeName, region, volumeType, 10),
			},
			{
				ResourceName:      "thalassa_block_volume.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBlockVolumeConfig(name, region, volumeType string, sizeGB int) string {
	return fmt.Sprintf(`
%s

resource "thalassa_block_volume" "test" {
  name             = %q
  region           = %q
  volume_type      = %q
  size_gb          = %d
  wait_until_ready = true
}
`, testAccProviderBlock(), name, region, volumeType, sizeGB)
}

func testAccBlockVolumeConfigWithDescription(name, region, volumeType string, sizeGB int, description string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_block_volume" "test" {
  name             = %q
  region           = %q
  volume_type      = %q
  size_gb          = %d
  description      = %q
  wait_until_ready = true
}
`, testAccProviderBlock(), name, region, volumeType, sizeGB, description)
}
