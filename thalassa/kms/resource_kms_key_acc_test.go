package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccKmsKey_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-kms")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyConfig(name, region, false, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "key_type", "aes256-gcm96"),
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "status", "active"),
					resource.TestCheckResourceAttrSet("thalassa_kms_key.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_kms_key.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_kms_key.test", "latest_version"),
				),
			},
		},
	})
}

func TestAccKmsKey_updateRotation(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-kms")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyConfig(name, region, false, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "key_rotation_enabled", "false"),
				),
			},
			{
				Config: testAccKmsKeyConfig(name, region, true, 90),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "key_rotation_enabled", "true"),
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "rotation_period_in_days", "90"),
				),
			},
		},
	})
}

func TestAccKmsKey_updateStatus(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-kms")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyConfigWithStatus(name, region, "active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "status", "active"),
				),
			},
			{
				Config: testAccKmsKeyConfigWithStatus(name, region, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "status", "disabled"),
				),
			},
			{
				Config: testAccKmsKeyConfigWithStatus(name, region, "active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_kms_key.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccKmsKey_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-kms")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyConfig(name, region, false, 0),
			},
			{
				ResourceName:      "thalassa_kms_key.test",
				ImportState:       true,
				ImportStateIdFunc: testAccKmsKeyImportStateIDFunc(region),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKmsKeyDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-kms")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSourceConfigByName(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_kms_key.test", "id", "thalassa_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_kms_key.test", "slug", "thalassa_kms_key.test", "slug"),
					resource.TestCheckResourceAttr("data.thalassa_kms_key.test", "name", name),
					resource.TestCheckResourceAttr("data.thalassa_kms_key.test", "region", region),
					resource.TestCheckResourceAttr("data.thalassa_kms_key.test", "key_type", "aes256-gcm96"),
					resource.TestCheckResourceAttrSet("data.thalassa_kms_key.test", "status"),
				),
			},
		},
	})
}

func testAccKmsKeyConfig(name, region string, rotationEnabled bool, rotationDays int) string {
	rotationBlock := ""
	if rotationEnabled {
		rotationBlock = fmt.Sprintf(`
  key_rotation_enabled    = true
  rotation_period_in_days = %d
`, rotationDays)
	}

	return fmt.Sprintf(`
%s

resource "thalassa_kms_key" "test" {
  name     = %q
  region   = %q
  key_type = "aes256-gcm96"
%s
}
`, testAccProviderBlock(), name, region, rotationBlock)
}

func testAccKmsKeyConfigWithStatus(name, region, status string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_kms_key" "test" {
  name     = %q
  region   = %q
  key_type = "aes256-gcm96"
  status   = %q
}
`, testAccProviderBlock(), name, region, status)
}

func testAccKmsKeyDataSourceConfigByName(name, region string) string {
	return fmt.Sprintf(`
%s

%s

data "thalassa_kms_key" "test" {
  name   = thalassa_kms_key.test.name
  region = thalassa_kms_key.test.region
}
`, testAccProviderBlock(), testAccKmsKeyConfigBlock(name, region))
}

func testAccKmsKeyImportStateIDFunc(region string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources["thalassa_kms_key.test"]
		if !ok {
			return "", fmt.Errorf("resource thalassa_kms_key.test not found in state")
		}
		if rs.Primary == nil || rs.Primary.ID == "" {
			return "", fmt.Errorf("resource thalassa_kms_key.test has no primary ID")
		}

		return fmt.Sprintf("%s/%s", region, rs.Primary.ID), nil
	}
}
