package tfs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTfsInstance_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	tfsName := acctest.RandomWithPrefix("tf-acc-tfs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTfsInstanceConfig(vpcName, subnetName, tfsName, region, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_tfs_instance.test", "name", tfsName),
					resource.TestCheckResourceAttr("thalassa_tfs_instance.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_tfs_instance.test", "size_gb", "10"),
					resource.TestCheckResourceAttrPair("thalassa_tfs_instance.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttrPair("thalassa_tfs_instance.test", "subnet_id", "thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_tfs_instance.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_tfs_instance.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_tfs_instance.test", "status"),
					resource.TestCheckResourceAttr("thalassa_tfs_instance.test", "endpoints.#", "1"),
					resource.TestCheckResourceAttrSet("thalassa_tfs_instance.test", "endpoints.0.address"),
				),
			},
		},
	})
}

func TestAccTfsInstance_update(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	tfsName := acctest.RandomWithPrefix("tf-acc-tfs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTfsInstanceConfigWithDescription(vpcName, subnetName, tfsName, region, 10, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_tfs_instance.test", "description", "initial description"),
				),
			},
			{
				Config: testAccTfsInstanceConfigWithDescription(vpcName, subnetName, tfsName, region, 10, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_tfs_instance.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccTfsInstance_import(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	tfsName := acctest.RandomWithPrefix("tf-acc-tfs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTfsInstanceConfig(vpcName, subnetName, tfsName, region, 10),
			},
			{
				ResourceName:            "thalassa_tfs_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_until_available"},
			},
		},
	})
}

func TestAccTfsInstanceDataSource_byName(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	tfsName := acctest.RandomWithPrefix("tf-acc-tfs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTfsInstanceDataSourceConfig(vpcName, subnetName, tfsName, region, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_tfs_instance.test", "id", "thalassa_tfs_instance.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_tfs_instance.test", "name", tfsName),
					resource.TestCheckResourceAttrPair("data.thalassa_tfs_instance.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_tfs_instance.test", "subnet_id", "thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_tfs_instance.test", "size_gb", "10"),
					resource.TestCheckResourceAttrSet("data.thalassa_tfs_instance.test", "status"),
				),
			},
		},
	})
}

func testAccProviderBlock() string {
	return `provider "thalassa" {}`
}

func testAccVpcSubnetConfigBlock(vpcName, subnetName, region string) string {
	return fmt.Sprintf(`
resource "thalassa_vpc" "test" {
  name   = %q
  region = %q
  cidrs  = ["10.0.0.0/16"]
}

resource "thalassa_subnet" "test" {
  vpc_id = thalassa_vpc.test.id
  name   = %q
  cidr   = "10.0.1.0/24"
}
`, vpcName, region, subnetName)
}

func testAccTfsInstanceConfig(vpcName, subnetName, tfsName, region string, sizeGB int) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_tfs_instance" "test" {
  name                 = %q
  region               = %q
  vpc_id               = thalassa_vpc.test.id
  subnet_id            = thalassa_subnet.test.id
  size_gb              = %d
  wait_until_available = true
}
`, testAccProviderBlock(), testAccVpcSubnetConfigBlock(vpcName, subnetName, region), tfsName, region, sizeGB)
}

func testAccTfsInstanceConfigWithDescription(vpcName, subnetName, tfsName, region string, sizeGB int, description string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_tfs_instance" "test" {
  name                 = %q
  region               = %q
  vpc_id               = thalassa_vpc.test.id
  subnet_id            = thalassa_subnet.test.id
  size_gb              = %d
  description          = %q
  wait_until_available = true
}
`, testAccProviderBlock(), testAccVpcSubnetConfigBlock(vpcName, subnetName, region), tfsName, region, sizeGB, description)
}

func testAccTfsInstanceDataSourceConfig(vpcName, subnetName, tfsName, region string, sizeGB int) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_tfs_instance" "test" {
  name                 = %q
  region               = %q
  vpc_id               = thalassa_vpc.test.id
  subnet_id            = thalassa_subnet.test.id
  size_gb              = %d
  wait_until_available = true
}

data "thalassa_tfs_instance" "test" {
  name = thalassa_tfs_instance.test.name
}
`, testAccProviderBlock(), testAccVpcSubnetConfigBlock(vpcName, subnetName, region), tfsName, region, sizeGB)
}
