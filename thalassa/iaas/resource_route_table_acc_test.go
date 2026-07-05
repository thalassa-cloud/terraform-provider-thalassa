package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRouteTable_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	rtName := acctest.RandomWithPrefix("tf-acc-rt")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig(vpcName, rtName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_route_table.test", "name", rtName),
					resource.TestCheckResourceAttrPair("thalassa_route_table.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_route_table.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_route_table.test", "slug"),
				),
			},
		},
	})
}

func TestAccRouteTable_update(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	rtName := acctest.RandomWithPrefix("tf-acc-rt")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfigWithDescription(vpcName, rtName, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_route_table.test", "description", "initial description"),
				),
			},
			{
				Config: testAccRouteTableConfigWithDescription(vpcName, rtName, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_route_table.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccRouteTable_import(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	rtName := acctest.RandomWithPrefix("tf-acc-rt")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig(vpcName, rtName, region),
			},
			{
				ResourceName:      "thalassa_route_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRouteTableConfig(vpcName, rtName, region string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_route_table" "test" {
  name   = %q
  vpc_id = thalassa_vpc.test.id
}
`, testAccProviderBlock(), testAccVpcConfigBlock(vpcName, region, "10.0.0.0/16"), rtName)
}

func testAccRouteTableConfigWithDescription(vpcName, rtName, region, description string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_route_table" "test" {
  name        = %q
  vpc_id      = thalassa_vpc.test.id
  description = %q
}
`, testAccProviderBlock(), testAccVpcConfigBlock(vpcName, region, "10.0.0.0/16"), rtName, description)
}
