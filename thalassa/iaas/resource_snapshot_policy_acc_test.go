package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSnapshotPolicy_basic(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tf-acc-snap-pol")
	region := testAccRegion()
	selectorValue := acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotPolicySelectorConfig(policyName, region, selectorValue, "7d", 3, true, "0 2 * * *"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "name", policyName),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "description", "acceptance test snapshot policy"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "ttl", "7d"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "keep_count", "3"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "enabled", "true"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "schedule", "0 2 * * *"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "timezone", "UTC"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "target.0.type", "selector"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "target.0.selector.backup", selectorValue),
					resource.TestCheckResourceAttrSet("thalassa_snapshot_policy.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot_policy.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot_policy.test", "next_snapshot_at"),
				),
			},
		},
	})
}

func TestAccSnapshotPolicy_explicitTarget(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tf-acc-vol")
	policyName := acctest.RandomWithPrefix("tf-acc-snap-pol")
	region := testAccRegion()
	volumeType := testAccVolumeType()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotPolicyExplicitConfig(volumeName, policyName, region, volumeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "name", policyName),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "target.0.type", "explicit"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "target.0.volume_identities.#", "1"),
					resource.TestCheckResourceAttrPair("thalassa_snapshot_policy.test", "target.0.volume_identities.0", "thalassa_block_volume.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_snapshot_policy.test", "id"),
				),
			},
		},
	})
}

func TestAccSnapshotPolicy_update(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tf-acc-snap-pol")
	region := testAccRegion()
	selectorValue := acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotPolicySelectorConfig(policyName, region, selectorValue, "7d", 3, true, "0 2 * * *"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "ttl", "7d"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "keep_count", "3"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "enabled", "true"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "schedule", "0 2 * * *"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "description", "acceptance test snapshot policy"),
				),
			},
			{
				Config: testAccSnapshotPolicySelectorConfig(policyName, region, selectorValue, "14d", 5, false, "0 4 * * *"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "ttl", "14d"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "keep_count", "5"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "schedule", "0 4 * * *"),
					resource.TestCheckResourceAttr("thalassa_snapshot_policy.test", "description", "acceptance test snapshot policy"),
				),
			},
		},
	})
}

func TestAccSnapshotPolicy_import(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tf-acc-snap-pol")
	region := testAccRegion()
	selectorValue := acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotPolicySelectorConfig(policyName, region, selectorValue, "7d", 3, true, "0 2 * * *"),
			},
			{
				ResourceName:      "thalassa_snapshot_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ttl",
				},
			},
		},
	})
}

func TestAccSnapshotPolicyDataSource_byName(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tf-acc-snap-pol")
	region := testAccRegion()
	selectorValue := acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotPolicyDataSourceConfigByName(policyName, region, selectorValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_snapshot_policy.test", "id", "thalassa_snapshot_policy.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_snapshot_policy.test", "slug", "thalassa_snapshot_policy.test", "slug"),
					resource.TestCheckResourceAttr("data.thalassa_snapshot_policy.test", "name", policyName),
					resource.TestCheckResourceAttr("data.thalassa_snapshot_policy.test", "schedule", "0 2 * * *"),
					resource.TestCheckResourceAttr("data.thalassa_snapshot_policy.test", "timezone", "UTC"),
					resource.TestCheckResourceAttr("data.thalassa_snapshot_policy.test", "target.0.type", "selector"),
					resource.TestCheckResourceAttrSet("data.thalassa_snapshot_policy.test", "enabled"),
				),
			},
		},
	})
}

func testAccSnapshotPolicySelectorConfig(name, region, selectorValue, ttl string, keepCount int, enabled bool, schedule string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_snapshot_policy" "test" {
  name        = %q
  description = "acceptance test snapshot policy"
  region      = %q
  ttl         = %q
  keep_count  = %d
  enabled     = %t
  schedule    = %q
  timezone    = "UTC"

  target {
    type = "selector"
    selector = {
      backup = %q
    }
  }
}
`, testAccProviderBlock(), name, region, ttl, keepCount, enabled, schedule, selectorValue)
}

func testAccSnapshotPolicyExplicitConfig(volumeName, policyName, region, volumeType string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_snapshot_policy" "test" {
  name        = %q
  description = "acceptance test snapshot policy"
  region      = %q
  ttl         = "7d"
  keep_count  = 3
  enabled     = true
  schedule    = "0 2 * * *"
  timezone    = "UTC"

  target {
    type              = "explicit"
    volume_identities = [thalassa_block_volume.test.id]
  }
}
`, testAccProviderBlock(), testAccBlockVolumeConfigBlock(volumeName, region, volumeType, 10), policyName, region)
}

func testAccSnapshotPolicyDataSourceConfigByName(name, region, selectorValue string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_snapshot_policy" "test" {
  name        = %q
  description = "acceptance test snapshot policy"
  region      = %q
  ttl         = "7d"
  keep_count  = 3
  enabled     = true
  schedule    = "0 2 * * *"
  timezone    = "UTC"

  target {
    type = "selector"
    selector = {
      backup = %q
    }
  }
}

data "thalassa_snapshot_policy" "test" {
  name = thalassa_snapshot_policy.test.name
}
`, testAccProviderBlock(), name, region, selectorValue)
}
