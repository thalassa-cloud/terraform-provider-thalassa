package iaas_test

import "fmt"

func testAccProviderBlock() string {
	return `provider "thalassa" {}`
}

func testAccVpcConfigBlock(vpcName, region, vpcCIDR string) string {
	return fmt.Sprintf(`
resource "thalassa_vpc" "test" {
  name   = %q
  region = %q
  cidrs  = [%q]
}
`, vpcName, region, vpcCIDR)
}

func testAccVpcSubnetConfigBlock(vpcName, subnetName, region string) string {
	return testAccVpcConfigBlock(vpcName, region, "10.0.0.0/16") + fmt.Sprintf(`
resource "thalassa_subnet" "test" {
  vpc_id = thalassa_vpc.test.id
  name   = %q
  cidr   = "10.0.1.0/24"
}
`, subnetName)
}

func testAccBlockVolumeConfigBlock(name, region, volumeType string, sizeGB int) string {
	return fmt.Sprintf(`
resource "thalassa_block_volume" "test" {
  name             = %q
  region           = %q
  volume_type      = %q
  size_gb          = %d
  wait_until_ready = true
}
`, name, region, volumeType, sizeGB)
}

func testAccDualVpcConfigBlock(requesterName, accepterName, region string) string {
	return fmt.Sprintf(`
resource "thalassa_vpc" "requester" {
  name   = %q
  region = %q
  cidrs  = ["10.0.0.0/16"]
}

resource "thalassa_vpc" "accepter" {
  name   = %q
  region = %q
  cidrs  = ["10.1.0.0/16"]
}
`, requesterName, region, accepterName, region)
}
