package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsRecord_basic(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordConfig(zoneName, "www", "A", 300, []string{"192.0.2.1"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("thalassa_dns_record.test", "zone_id", "thalassa_dns_zone.test", "id"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "name", "www"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "ttl", "300"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.#", "1"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.0", "192.0.2.1"),
					resource.TestCheckResourceAttrSet("thalassa_dns_record.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_dns_record.test", "created_at"),
				),
			},
		},
	})
}

func TestAccDnsRecord_update(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordConfig(zoneName, "www", "A", 300, []string{"192.0.2.1"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "ttl", "300"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.#", "1"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.0", "192.0.2.1"),
				),
			},
			{
				Config: testAccDnsRecordConfig(zoneName, "www", "A", 600, []string{"192.0.2.1", "192.0.2.2"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "ttl", "600"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.#", "2"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.0", "192.0.2.1"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.1", "192.0.2.2"),
				),
			},
		},
	})
}

func TestAccDnsRecord_txt(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordConfig(zoneName, "@", "TXT", 3600, []string{"v=spf1 -all"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "name", "@"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "type", "TXT"),
					resource.TestCheckResourceAttr("thalassa_dns_record.test", "values.0", "v=spf1 -all"),
					resource.TestCheckResourceAttrSet("thalassa_dns_record.test", "id"),
				),
			},
		},
	})
}

func TestAccDnsRecord_import(t *testing.T) {
	zoneName := testAccDnsZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordConfig(zoneName, "www", "A", 300, []string{"192.0.2.1"}),
			},
			{
				ResourceName:      "thalassa_dns_record.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Keep the prior working directory so zone_id can resolve from
				// configuration during post-import refresh (record Read needs zone_id).
				ImportStatePersist: true,
			},
		},
	})
}

func testAccDnsRecordConfig(zoneName, name, recordType string, ttl int, values []string) string {
	valuesBlock := ""
	for _, value := range values {
		if valuesBlock != "" {
			valuesBlock += ", "
		}
		valuesBlock += fmt.Sprintf("%q", value)
	}

	return fmt.Sprintf(`
%s

%s

resource "thalassa_dns_record" "test" {
  zone_id = thalassa_dns_zone.test.id
  name    = %q
  type    = %q
  ttl     = %d
  values  = [%s]
}
`, testAccProviderBlock(), testAccDnsZoneConfigBlock(zoneName, "acceptance dns zone for records"), name, recordType, ttl, valuesBlock)
}
