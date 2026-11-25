package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLoadBalancer_securityGroupAttachmentUpdate(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tf-acc-vpc")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	lbName := acctest.RandomWithPrefix("tf-acc-lb")
	// Security group names are limited to 16 characters.
	sgName := fmt.Sprintf("sg%d", acctest.RandIntRange(100000, 999999))
	region := testAccRegion()

	var loadBalancerID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig(vpcName, subnetName, lbName, sgName, region, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_loadbalancer.test", "name", lbName),
					resource.TestCheckResourceAttr("thalassa_loadbalancer.test", "region", region),
					resource.TestCheckResourceAttrPair("thalassa_loadbalancer.test", "subnet_id", "thalassa_subnet.test", "id"),
					resource.TestCheckResourceAttrPair("thalassa_loadbalancer.test", "vpc_id", "thalassa_vpc.test", "id"),
					testAccCheckNoSecurityGroupAttachments("thalassa_loadbalancer.test"),
					testAccCheckLoadBalancerIdentityDiffersFromSlug("thalassa_loadbalancer.test"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["thalassa_loadbalancer.test"]
						if !ok {
							return fmt.Errorf("resource thalassa_loadbalancer.test not found")
						}
						loadBalancerID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				Config: testAccLoadBalancerConfig(vpcName, subnetName, lbName, sgName, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceIDUnchanged("thalassa_loadbalancer.test", &loadBalancerID),
					testAccCheckLoadBalancerIdentityDiffersFromSlug("thalassa_loadbalancer.test"),
					resource.TestCheckResourceAttr("thalassa_loadbalancer.test", "security_group_attachments.#", "1"),
					resource.TestCheckResourceAttrPair("thalassa_loadbalancer.test", "security_group_attachments.0", "thalassa_security_group.test", "id"),
				),
			},
		},
	})
}

func testAccCheckResourceIDUnchanged(resourceName string, previousID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		if *previousID == "" {
			return fmt.Errorf("resource %s has no recorded identity", resourceName)
		}
		if rs.Primary.ID != *previousID {
			return fmt.Errorf("resource %s identity changed from %q to %q", resourceName, *previousID, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNoSecurityGroupAttachments(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		count := rs.Primary.Attributes["security_group_attachments.#"]
		if count != "" && count != "0" {
			return fmt.Errorf("expected no security group attachments on %s, got %s", resourceName, count)
		}
		return nil
	}
}

func testAccCheckLoadBalancerIdentityDiffersFromSlug(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no identity", resourceName)
		}
		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			return fmt.Errorf("resource %s has no slug", resourceName)
		}
		if rs.Primary.ID == slug {
			return fmt.Errorf("resource %s identity equals slug %q; the update would not exercise identity addressing", resourceName, slug)
		}
		return nil
	}
}

func testAccLoadBalancerConfig(vpcName, subnetName, lbName, sgName, region string, attachSecurityGroup bool) string {
	attachments := ""
	if attachSecurityGroup {
		attachments = "\n  security_group_attachments = [thalassa_security_group.test.id]\n"
	}

	return fmt.Sprintf(`
%s

%s

resource "thalassa_security_group" "test" {
  name   = %q
  vpc_id = thalassa_vpc.test.id
}

resource "thalassa_loadbalancer" "test" {
  name      = %q
  region    = %q
  subnet_id = thalassa_subnet.test.id
%s}
`, testAccProviderBlock(), testAccVpcSubnetConfigBlock(vpcName, subnetName, region), sgName, lbName, region, attachments)
}
