package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSnapshot_basic(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	snapshotName := acctest.RandomWithPrefix("tf-acc-snap")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfig(volumeName, snapshotName, region, volumeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "name", snapshotName),
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "description", "acceptance test snapshot"),
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "delete_protection", "false"),
					resource.TestCheckResourceAttrPair("thalassa_snapshot.test", "volume_identity", "thalassa_block_volume.test", "id"),
					resource.TestCheckResourceAttrPair("thalassa_snapshot.test", "source_volume_id", "thalassa_block_volume.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot.test", "status"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot.test", "region"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot.test", "size_gb"),
				),
			},
		},
	})
}

func TestAccSnapshot_update(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	snapshotName := acctest.RandomWithPrefix("tf-acc-snap")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfigWithDescription(volumeName, snapshotName, region, volumeType, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "description", "initial description"),
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "labels.environment", "test"),
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "annotations.managed_by", "terraform"),
				),
			},
			{
				Config: testAccSnapshotConfigWithDescription(volumeName, snapshotName, region, volumeType, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "description", "updated description"),
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "labels.environment", "test"),
					resource.TestCheckResourceAttr("thalassa_snapshot.test", "annotations.managed_by", "terraform"),
				),
			},
		},
	})
}

func TestAccSnapshot_import(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	snapshotName := acctest.RandomWithPrefix("tf-acc-snap")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfig(volumeName, snapshotName, region, volumeType),
			},
			{
				ResourceName:      "thalassa_snapshot.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"wait_until_available",
				},
			},
		},
	})
}

func TestAccSnapshotDataSource_byName(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	snapshotName := acctest.RandomWithPrefix("tf-acc-snap")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotDataSourceConfigByName(volumeName, snapshotName, region, volumeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_snapshot.test", "id", "thalassa_snapshot.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_snapshot.test", "slug", "thalassa_snapshot.test", "slug"),
					resource.TestCheckResourceAttr("data.thalassa_snapshot.test", "name", snapshotName),
					resource.TestCheckResourceAttrPair("data.thalassa_snapshot.test", "source_volume_id", "thalassa_block_volume.test", "id"),
					resource.TestCheckResourceAttrSet("data.thalassa_snapshot.test", "status"),
				),
			},
		},
	})
}

func testAccSnapshotConfig(volumeName, snapshotName, region, volumeType string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_snapshot" "test" {
  name                 = %q
  description          = "acceptance test snapshot"
  volume_identity      = thalassa_block_volume.test.id
  wait_until_available = true
}
`, testAccProviderBlock(), testAccBlockVolumeConfigBlock(volumeName, region, volumeType, 10), snapshotName)
}

func testAccSnapshotConfigWithDescription(volumeName, snapshotName, region, volumeType, description string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_snapshot" "test" {
  name                 = %q
  description          = %q
  volume_identity      = thalassa_block_volume.test.id
  wait_until_available = true

  labels = {
    environment = "test"
  }

  annotations = {
    managed_by = "terraform"
  }
}
`, testAccProviderBlock(), testAccBlockVolumeConfigBlock(volumeName, region, volumeType, 10), snapshotName, description)
}

func testAccSnapshotDataSourceConfigByName(volumeName, snapshotName, region, volumeType string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_snapshot" "test" {
  name                 = %q
  description          = "acceptance test snapshot"
  volume_identity      = thalassa_block_volume.test.id
  wait_until_available = true
}

data "thalassa_snapshot" "test" {
  name = thalassa_snapshot.test.name
}
`, testAccProviderBlock(), testAccBlockVolumeConfigBlock(volumeName, region, volumeType, 10), snapshotName)
}
