package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpc_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-vpc")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig(name, region, "10.0.0.0/16"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_vpc.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_vpc.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("thalassa_vpc.test", "cidrs.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttrSet("thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_vpc.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_vpc.test", "status"),
				),
			},
		},
	})
}

func TestAccVpc_withOptionalAttributes(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-vpc")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfigWithOptionalAttributes(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_vpc.test", "description", "acceptance test vpc"),
					resource.TestCheckResourceAttr("thalassa_vpc.test", "labels.environment", "test"),
					resource.TestCheckResourceAttr("thalassa_vpc.test", "annotations.managed_by", "terraform"),
				),
			},
		},
	})
}

func TestAccVpc_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-vpc")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfigWithDescription(name, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc.test", "description", "initial description"),
				),
			},
			{
				Config: testAccVpcConfigWithDescription(name, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccVpcDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-vpc")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDataSourceConfig(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_vpc.test", "id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_vpc.test", "name", name),
					resource.TestCheckResourceAttr("data.thalassa_vpc.test", "region", region),
					resource.TestCheckResourceAttr("data.thalassa_vpc.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("data.thalassa_vpc.test", "cidrs.0", "10.0.0.0/16"),
				),
			},
		},
	})
}

func testAccVpcConfig(name, region, cidr string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

resource "thalassa_vpc" "test" {
  name   = %q
  region = %q
  cidrs  = [%q]
}
`, name, region, cidr)
}

func testAccVpcConfigWithDescription(name, region, description string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

resource "thalassa_vpc" "test" {
  name        = %q
  region      = %q
  cidrs       = ["10.0.0.0/16"]
  description = %q
}
`, name, region, description)
}

func testAccVpcConfigWithOptionalAttributes(name, region string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

resource "thalassa_vpc" "test" {
  name        = %q
  region      = %q
  cidrs       = ["10.0.0.0/16"]
  description = "acceptance test vpc"

  labels = {
    environment = "test"
  }

  annotations = {
    managed_by = "terraform"
  }
}
`, name, region)
}

func testAccVpcDataSourceConfig(name, region string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

resource "thalassa_vpc" "test" {
  name   = %q
  region = %q
  cidrs  = ["10.0.0.0/16"]
}

data "thalassa_vpc" "test" {
  name   = thalassa_vpc.test.name
  region = thalassa_vpc.test.region
}
`, name, region)
}
