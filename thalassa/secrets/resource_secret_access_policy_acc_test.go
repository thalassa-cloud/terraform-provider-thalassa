package secrets_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretAccessPolicy_basic(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretAccessPolicyConfig(kmsName, region, path, "read", "team:platform"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "path", path),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.#", "1"),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.effect", "Allow"),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.actions.#", "1"),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.actions.0", "read"),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.principals.#", "1"),
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.principals.0", "team:platform"),
				),
			},
		},
	})
}

func TestAccSecretAccessPolicy_update(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretAccessPolicyConfig(kmsName, region, path, "read", "team:platform"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.actions.0", "read"),
				),
			},
			{
				Config: testAccSecretAccessPolicyConfig(kmsName, region, path, "read", "team:engineering"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_secret_access_policy.test", "statement.0.principals.0", "team:engineering"),
				),
			},
		},
	})
}

func TestAccSecretAccessPolicy_import(t *testing.T) {
	kmsName := acctest.RandomWithPrefix("tf-acc-kms")
	path := testAccSecretPath(acctest.RandomWithPrefix("tf_acc_secret"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretAccessPolicyConfig(kmsName, region, path, "read", "team:platform"),
			},
			{
				ResourceName:      "thalassa_secret_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"statement",
				},
			},
		},
	})
}

func testAccSecretAccessPolicyConfig(kmsName, region, path, action, principal string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_secret_access_policy" "test" {
  region = %q
  path   = thalassa_secret.test.path

  statement {
    effect     = "Allow"
    actions    = [%q]
    principals = [%q]
  }
}
`, testAccSecretBaseConfigBlock(kmsName, region, path), region, action, principal)
}
