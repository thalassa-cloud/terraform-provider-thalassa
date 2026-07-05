package iaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcPeeringConnection_basic(t *testing.T) {
	requesterName := acctest.RandomWithPrefix("tf-acc-vpc-a")
	accepterName := acctest.RandomWithPrefix("tf-acc-vpc-b")
	peeringName := acctest.RandomWithPrefix("tf-acc-peer")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionConfig(requesterName, accepterName, peeringName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc_peering_connection.test", "name", peeringName),
					resource.TestCheckResourceAttrPair("thalassa_vpc_peering_connection.test", "requester_vpc_id", "thalassa_vpc.requester", "id"),
					resource.TestCheckResourceAttrPair("thalassa_vpc_peering_connection.test", "accepter_vpc_id", "thalassa_vpc.accepter", "id"),
					resource.TestCheckResourceAttr("thalassa_vpc_peering_connection.test", "auto_accept", "true"),
					resource.TestCheckResourceAttrSet("thalassa_vpc_peering_connection.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_vpc_peering_connection.test", "status"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnection_update(t *testing.T) {
	requesterName := acctest.RandomWithPrefix("tf-acc-vpc-a")
	accepterName := acctest.RandomWithPrefix("tf-acc-vpc-b")
	peeringName := acctest.RandomWithPrefix("tf-acc-peer")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionConfigWithDescription(requesterName, accepterName, peeringName, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc_peering_connection.test", "description", "initial description"),
				),
			},
			{
				Config: testAccVpcPeeringConnectionConfigWithDescription(requesterName, accepterName, peeringName, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_vpc_peering_connection.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnection_import(t *testing.T) {
	requesterName := acctest.RandomWithPrefix("tf-acc-vpc-a")
	accepterName := acctest.RandomWithPrefix("tf-acc-vpc-b")
	peeringName := acctest.RandomWithPrefix("tf-acc-peer")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionConfig(requesterName, accepterName, peeringName, region),
			},
			{
				ResourceName:      "thalassa_vpc_peering_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"wait_for_active",
					"wait_for_active_timeout",
					"wait_for_deleted_timeout",
				},
			},
		},
	})
}

func testAccVpcPeeringConnectionConfig(requesterName, accepterName, peeringName, region string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_vpc_peering_connection" "test" {
  name             = %q
  requester_vpc_id = thalassa_vpc.requester.id
  accepter_vpc_id  = thalassa_vpc.accepter.id
  auto_accept      = true
  wait_for_active  = true
}
`, testAccProviderBlock(), testAccDualVpcConfigBlock(requesterName, accepterName, region), peeringName)
}

func testAccVpcPeeringConnectionConfigWithDescription(requesterName, accepterName, peeringName, region, description string) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_vpc_peering_connection" "test" {
  name             = %q
  description      = %q
  requester_vpc_id = thalassa_vpc.requester.id
  accepter_vpc_id  = thalassa_vpc.accepter.id
  auto_accept      = true
  wait_for_active  = true
}
`, testAccProviderBlock(), testAccDualVpcConfigBlock(requesterName, accepterName, region), peeringName, description)
}
