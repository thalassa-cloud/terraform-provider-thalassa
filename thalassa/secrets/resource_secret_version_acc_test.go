package secrets_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretVersion_basic(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretVersionConfig(kmsName, region, path, "version-one"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret_version.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_secret_version.test", "path", path),
					resource.TestCheckResourceAttrSet("thalassa_secret_version.test", "version"),
					resource.TestCheckResourceAttrSet("thalassa_secret_version.test", "id"),
				),
			},
		},
	})
}

func TestAccSecretVersion_addSecondVersion(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretVersionConfig(kmsName, region, path, "version-one"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret_version.test", "version", "2"),
				),
			},
			{
				Config: testAccSecretVersionConfigWithSecond(kmsName, region, path, "version-one", "version-two"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret_version.test", "version", "2"),
					resource.TestCheckResourceAttr("thalassa_secret_version.second", "version", "3"),
				),
			},
		},
	})
}

func TestAccSecretVersion_import(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretVersionConfig(kmsName, region, path, "import-version"),
			},
			{
				ResourceName:            "thalassa_secret_version.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_string", "generate_secret", "secret_key_values"},
			},
		},
	})
}

func testAccSecretVersionConfig(kmsName, region, path, value string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_secret_version" "test" {
  region        = %q
  path          = thalassa_secret.test.path
  secret_string = %q
}
`, testAccSecretBaseConfigBlock(kmsName, region, path), region, value)
}

func testAccSecretVersionConfigWithSecond(kmsName, region, path, firstValue, secondValue string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_secret_version" "test" {
  region        = %q
  path          = thalassa_secret.test.path
  secret_string = %q
}

resource "thalassa_secret_version" "second" {
  region        = %q
  path          = thalassa_secret.test.path
  secret_string = %q
}
`, testAccSecretBaseConfigBlock(kmsName, region, path), region, firstValue, region, secondValue)
}
