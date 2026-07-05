package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecurityGroup_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	sgName := acctest.RandomWithPrefix("sg")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig(vpcName, sgName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_security_group.test", "name", sgName),
					resource.TestCheckResourceAttrPair("thalassa_security_group.test", "vpc_id", "thalassa_vpc.test", "id"),
					resource.TestCheckResourceAttr("thalassa_security_group.test", "allow_same_group_traffic", "false"),
					resource.TestCheckResourceAttrSet("thalassa_security_group.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_security_group.test", "identity"),
					resource.TestCheckResourceAttrSet("thalassa_security_group.test", "status"),
				),
			},
		},
	})
}

func TestAccSecurityGroup_update(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	sgName := acctest.RandomWithPrefix("sg")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfigWithDescription(vpcName, sgName, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_security_group.test", "description", "initial description"),
				),
			},
			{
				Config: testAccSecurityGroupConfigWithDescription(vpcName, sgName, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_security_group.test", "description", "updated description"),
				),
			},
		},
	})
}

// func TestAccSecurityGroup_import(t *testing.T) {
// 	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
// 	sgName := acctest.RandomWithPrefix("sg")
// 	region := testAccRegion()

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSecurityGroupConfig(vpcName, sgName, region),
// 			},
// 			{
// 				ResourceName:            "thalassa_security_group.test",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"ingress_rule", "egress_rule"},
// 			},
// 		},
// 	})
// }

func testAccSecurityGroupConfig(vpcName, sgName, region string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_security_group" "test" {
  name   = %q
  vpc_id = thalassa_vpc.test.id
}
`, testAccProviderBlock(), testAccVpcConfigBlock(vpcName, region, "10.0.0.0/16"), sgName)
}

func testAccSecurityGroupConfigWithDescription(vpcName, sgName, region, description string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_security_group" "test" {
  name        = %q
  vpc_id      = thalassa_vpc.test.id
  description = %q
}
`, testAccProviderBlock(), testAccVpcConfigBlock(vpcName, region, "10.0.0.0/16"), sgName, description)
}
