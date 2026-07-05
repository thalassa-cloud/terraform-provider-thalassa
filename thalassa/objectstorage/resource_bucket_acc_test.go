package objectstorage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBucket_basic(t *testing.T) {
	bucketName := testAccBucketName(acctest.RandomWithPrefix("tf-acc-bucket"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig(bucketName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket.test", "versioning", "false"),
					resource.TestCheckResourceAttrSet("thalassa_objectstorage_bucket.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_objectstorage_bucket.test", "status"),
					resource.TestCheckResourceAttrSet("thalassa_objectstorage_bucket.test", "endpoint"),
				),
			},
		},
	})
}

func TestAccBucket_updateVersioning(t *testing.T) {
	bucketName := testAccBucketName(acctest.RandomWithPrefix("tf-acc-bucket"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigWithVersioning(bucketName, region, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket.test", "versioning", "false"),
				),
			},
			{
				Config: testAccBucketConfigWithVersioning(bucketName, region, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket.test", "versioning", "true"),
				),
			},
		},
	})
}

func TestAccBucketDataSource_byName(t *testing.T) {
	bucketName := testAccBucketName(acctest.RandomWithPrefix("tf-acc-bucket"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketDataSourceConfig(bucketName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_objectstorage_bucket.test", "id", "thalassa_objectstorage_bucket.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_objectstorage_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("data.thalassa_objectstorage_bucket.test", "region", region),
					resource.TestCheckResourceAttrSet("data.thalassa_objectstorage_bucket.test", "status"),
					resource.TestCheckResourceAttrSet("data.thalassa_objectstorage_bucket.test", "endpoint"),
				),
			},
		},
	})
}

func testAccProviderBlock() string {
	return `provider "thalassa" {}`
}

func testAccBucketConfig(name, region string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_objectstorage_bucket" "test" {
  name             = %q
  region           = %q
  wait_for_ready   = true
  wait_for_deleted = true
}
`, testAccProviderBlock(), name, region)
}

func testAccBucketConfigWithVersioning(name, region string, versioning bool) string {
	return fmt.Sprintf(`
%s

resource "thalassa_objectstorage_bucket" "test" {
  name             = %q
  region           = %q
  versioning       = %t
  wait_for_ready   = true
  wait_for_deleted = true
}
`, testAccProviderBlock(), name, region, versioning)
}

func testAccBucketDataSourceConfig(name, region string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_objectstorage_bucket" "test" {
  name             = %q
  region           = %q
  wait_for_ready   = true
  wait_for_deleted = true
}

data "thalassa_objectstorage_bucket" "test" {
  name   = thalassa_objectstorage_bucket.test.name
  region = thalassa_objectstorage_bucket.test.region
}
`, testAccProviderBlock(), name, region)
}
