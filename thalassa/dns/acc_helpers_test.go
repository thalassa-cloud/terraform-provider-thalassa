package dns_test

import "fmt"

func testAccProviderBlock() string {
	return `provider "thalassa" {}`
}

func testAccDnsZoneConfigBlock(zoneName, description string) string {
	return fmt.Sprintf(`
resource "thalassa_dns_zone" "test" {
  zone_name   = %q
  description = %q
}
`, zoneName, description)
}
