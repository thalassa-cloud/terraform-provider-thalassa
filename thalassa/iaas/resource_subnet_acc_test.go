package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSubnet_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()
	subnetCIDR := "10.0.1.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig(vpcName, subnetName, region, subnetCIDR),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("thalassa_subnet.test", "cidr", subnetCIDR),
					resource.TestCheckResourceAttrPair("thalassa_subnet.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_subnet.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_subnet.test", "status"),
					resource.TestCheckResourceAttrSet("thalassa_subnet.test", "type"),
					resource.TestCheckResourceAttrSet("thalassa_subnet.test", "ipv4_addresses_available"),
				),
			},
		},
	})
}

func TestAccSubnet_withOptionalAttributes(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigWithOptionalAttributes(vpcName, subnetName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("thalassa_subnet.test", "description", "acceptance test subnet"),
					resource.TestCheckResourceAttr("thalassa_subnet.test", "labels.environment", "test"),
					resource.TestCheckResourceAttr("thalassa_subnet.test", "annotations.managed_by", "terraform"),
				),
			},
		},
	})
}

func TestAccSubnet_update(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigWithDescription(vpcName, subnetName, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_subnet.test", "description", "initial description"),
				),
			},
			{
				Config: testAccSubnetConfigWithDescription(vpcName, subnetName, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_subnet.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccSubnet_updateLabelsAndAnnotations(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigWithLabels(vpcName, subnetName, region, "initial", "terraform"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_subnet.test", "labels.environment", "initial"),
					resource.TestCheckResourceAttr("thalassa_subnet.test", "annotations.managed_by", "terraform"),
				),
			},
			{
				Config: testAccSubnetConfigWithLabels(vpcName, subnetName, region, "updated", "terraform-acc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_subnet.test", "labels.environment", "updated"),
					resource.TestCheckResourceAttr("thalassa_subnet.test", "annotations.managed_by", "terraform-acc"),
				),
			},
		},
	})
}

func TestAccSubnet_import(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig(vpcName, subnetName, region, "10.0.1.0/24"),
			},
			{
				ResourceName:      "thalassa_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSubnetDataSource_byName(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfigByName(vpcName, subnetName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_subnet.test", "id", "thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_subnet.test", "slug", "thalassa_subnet.test", "slug"),
					resource.TestCheckResourceAttr("data.thalassa_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttrPair("data.thalassa_subnet.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_subnet.test", "cidr", "10.0.1.0/24"),
				),
			},
		},
	})
}

func TestAccSubnetDataSource_bySlug(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfigBySlug(vpcName, subnetName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_subnet.test", "id", "thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_subnet.test", "slug", "thalassa_subnet.test", "slug"),
					resource.TestCheckResourceAttr("data.thalassa_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("data.thalassa_subnet.test", "cidr", "10.0.1.0/24"),
				),
			},
		},
	})
}

func testAccSubnetVpcConfig(vpcName, region string) string {
	return fmt.Sprintf(`
resource "thalassa_vpc" "test" {
  name   = %q
  region = %q
  cidrs  = ["10.0.0.0/16"]
}
`, vpcName, region)
}

func testAccSubnetConfig(vpcName, subnetName, region, subnetCIDR string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

%s

resource "thalassa_subnet" "test" {
  vpc_id = thalassa_vpc.test.id
  name   = %q
  cidr   = %q
}
`, testAccSubnetVpcConfig(vpcName, region), subnetName, subnetCIDR)
}

func testAccSubnetConfigWithDescription(vpcName, subnetName, region, description string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

%s

resource "thalassa_subnet" "test" {
  vpc_id      = thalassa_vpc.test.id
  name        = %q
  cidr        = "10.0.1.0/24"
  description = %q
}
`, testAccSubnetVpcConfig(vpcName, region), subnetName, description)
}

func testAccSubnetConfigWithOptionalAttributes(vpcName, subnetName, region string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

%s

resource "thalassa_subnet" "test" {
  vpc_id      = thalassa_vpc.test.id
  name        = %q
  cidr        = "10.0.1.0/24"
  description = "acceptance test subnet"

  labels = {
    environment = "test"
  }

  annotations = {
    managed_by = "terraform"
  }
}
`, testAccSubnetVpcConfig(vpcName, region), subnetName)
}

func testAccSubnetConfigWithLabels(vpcName, subnetName, region, labelValue, annotationValue string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

%s

resource "thalassa_subnet" "test" {
  vpc_id = thalassa_vpc.test.id
  name   = %q
  cidr   = "10.0.1.0/24"

  labels = {
    environment = %q
  }

  annotations = {
    managed_by = %q
  }
}
`, testAccSubnetVpcConfig(vpcName, region), subnetName, labelValue, annotationValue)
}

func testAccSubnetDataSourceConfigByName(vpcName, subnetName, region string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

%s

resource "thalassa_subnet" "test" {
  vpc_id = thalassa_vpc.test.id
  name   = %q
  cidr   = "10.0.1.0/24"
}

data "thalassa_subnet" "test" {
  name   = thalassa_subnet.test.name
  vpc_id = thalassa_vpc.test.id
}
`, testAccSubnetVpcConfig(vpcName, region), subnetName)
}

func testAccSubnetDataSourceConfigBySlug(vpcName, subnetName, region string) string {
	return fmt.Sprintf(`
provider "thalassa" {}

%s

resource "thalassa_subnet" "test" {
  vpc_id = thalassa_vpc.test.id
  name   = %q
  cidr   = "10.0.1.0/24"
}

data "thalassa_subnet" "test" {
  name   = thalassa_subnet.test.name
  vpc_id = thalassa_vpc.test.id
  slug   = thalassa_subnet.test.slug
}
`, testAccSubnetVpcConfig(vpcName, region), subnetName)
}
