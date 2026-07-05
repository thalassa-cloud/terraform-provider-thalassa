package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNatGateway_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	natName := acctest.RandomWithPrefix("tf-acc-nat")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig(vpcName, subnetName, natName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_natgateway.test", "name", natName),
					resource.TestCheckResourceAttrPair("thalassa_natgateway.test", "subnet_id", "thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttrPair("thalassa_natgateway.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_natgateway.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_natgateway.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_natgateway.test", "status"),
					resource.TestCheckResourceAttrSet("thalassa_natgateway.test", "endpoint_ip"),
				),
			},
		},
	})
}

func TestAccNatGateway_update(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	natName := acctest.RandomWithPrefix("tf-acc-nat")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfigWithDescription(vpcName, subnetName, natName, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_natgateway.test", "description", "initial description"),
				),
			},
			{
				Config: testAccNatGatewayConfigWithDescription(vpcName, subnetName, natName, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_natgateway.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccNatGateway_import(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	natName := acctest.RandomWithPrefix("tf-acc-nat")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig(vpcName, subnetName, natName, region),
			},
			{
				ResourceName:      "thalassa_natgateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNatGatewayConfig(vpcName, subnetName, natName, region string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_natgateway" "test" {
  name      = %q
  subnet_id = thalassa_subnet.test.id
}
`, testAccProviderBlock(), testAccVpcSubnetConfigBlock(vpcName, subnetName, region), natName)
}

func testAccNatGatewayConfigWithDescription(vpcName, subnetName, natName, region, description string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_natgateway" "test" {
  name        = %q
  subnet_id   = thalassa_subnet.test.id
  description = %q
}
`, testAccProviderBlock(), testAccVpcSubnetConfigBlock(vpcName, subnetName, region), natName, description)
}
