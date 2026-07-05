package secrets_test

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

func testAccSecretBaseConfigBlock(kmsName, region, path string) string {
	return testAccKmsKeyConfigBlock(kmsName, region) + fmt.Sprintf(`
resource "thalassa_secret" "test" {
  region     = %q
  path       = %q
  kms_key_id = thalassa_kms_key.test.id

  generate_secret {
    byte_length = 32
  }
}
`, region, path)
}
