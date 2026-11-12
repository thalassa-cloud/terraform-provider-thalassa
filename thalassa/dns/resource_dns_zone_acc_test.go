package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsZone_basic(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsZoneConfig(zoneName, "acceptance dns zone"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "zone_name", zoneName),
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "description", "acceptance dns zone"),
					resource.TestCheckResourceAttrSet("thalassa_dns_zone.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_dns_zone.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_dns_zone.test", "created_at"),
				),
			},
		},
	})
}

func TestAccDnsZone_update(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsZoneConfigWithMeta(zoneName, "initial description", "initial", "terraform"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "description", "initial description"),
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "labels.environment", "initial"),
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "annotations.managed_by", "terraform"),
				),
			},
			{
				Config: testAccDnsZoneConfigWithMeta(zoneName, "updated description", "updated", "terraform-acc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "description", "updated description"),
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "labels.environment", "updated"),
					resource.TestCheckResourceAttr("thalassa_dns_zone.test", "annotations.managed_by", "terraform-acc"),
				),
			},
		},
	})
}

func TestAccDnsZone_import(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsZoneConfig(zoneName, "import me"),
			},
			{
				ResourceName:      "thalassa_dns_zone.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsZoneConfig(zoneName, description string) string {
	return fmt.Sprintf(`
%s

%s
`, testAccProviderBlock(), testAccDnsZoneConfigBlock(zoneName, description))
}

func testAccDnsZoneConfigWithMeta(zoneName, description, labelEnv, annotationManagedBy string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_dns_zone" "test" {
  zone_name   = %q
  description = %q

  labels = {
    environment = %q
  }

  annotations = {
    managed_by = %q
  }
}
`, testAccProviderBlock(), zoneName, description, labelEnv, annotationManagedBy)
}
