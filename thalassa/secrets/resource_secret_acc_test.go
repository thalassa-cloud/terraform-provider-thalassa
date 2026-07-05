package secrets_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecret_basic(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretConfigWithString(kmsName, region, path, "initial-value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_secret.test", "path", path),
					resource.TestCheckResourceAttrSet("thalassa_secret.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_secret.test", "current_version"),
				),
			},
		},
	})
}

func TestAccSecret_withGenerateSecret(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretConfigWithGenerate(kmsName, region, path, 32),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret.test", "path", path),
					resource.TestCheckResourceAttrSet("thalassa_secret.test", "current_version"),
				),
			},
		},
	})
}

func TestAccSecret_import(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretConfigWithString(kmsName, region, path, "import-test"),
			},
			{
				ResourceName:            "thalassa_secret.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_string", "generate_secret", "secret_key_values"},
			},
		},
	})
}

func testAccSecretConfigWithString(kmsName, region, path, value string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_secret" "test" {
  region        = %q
  path          = %q
  kms_key_id    = thalassa_kms_key.test.id
  secret_string = %q
}
`, testAccKmsKeyConfigBlock(kmsName, region), region, path, value)
}

func testAccSecretConfigWithGenerate(kmsName, region, path string, byteLength int) string {
	return fmt.Sprintf(`
%s

resource "thalassa_secret" "test" {
  region     = %q
  path       = %q
  kms_key_id = thalassa_kms_key.test.id

  generate_secret {
    byte_length = %d
  }
}
`, testAccKmsKeyConfigBlock(kmsName, region), region, path, byteLength)
}
