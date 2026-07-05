package kms_test

import "fmt"

func testAccProviderBlock() string {
	return `provider "thalassa" {}`
}

func testAccKmsKeyConfigBlock(name, region string) string {
	return fmt.Sprintf(`
resource "thalassa_kms_key" "test" {
  name     = %q
  region   = %q
  key_type = "aes256-gcm96"
}
`, name, region)
}
