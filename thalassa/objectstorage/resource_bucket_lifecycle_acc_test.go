package objectstorage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBucketLifecycle_basic(t *testing.T) {
	bucketName := testAccBucketName(acctest.RandomWithPrefix("tf-acc-bucket"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfig(bucketName, region, "expire-logs", 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "bucket_name", bucketName),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.#", "1"),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.id", "expire-logs"),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.prefix", "logs/"),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.status", "Enabled"),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.expiration.0.days", "30"),
				),
			},
		},
	})
}

func TestAccBucketLifecycle_update(t *testing.T) {
	bucketName := testAccBucketName(acctest.RandomWithPrefix("tf-acc-bucket"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfig(bucketName, region, "expire-logs", 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.expiration.0.days", "30"),
				),
			},
			{
				Config: testAccBucketLifecycleConfig(bucketName, region, "expire-logs", 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.expiration.0.days", "60"),
				),
			},
		},
	})
}

func TestAccBucketLifecycle_noncurrentVersion(t *testing.T) {
	bucketName := testAccBucketName(acctest.RandomWithPrefix("tf-acc-bucket"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleNoncurrentConfig(bucketName, region, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.#", "1"),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.id", "expire-noncurrent"),
					resource.TestCheckResourceAttr("thalassa_objectstorage_bucket_lifecycle.test", "rule.0.noncurrent_version_expiration.0.noncurrent_days", "7"),
				),
			},
		},
	})
}

func testAccBucketBaseConfig(name, region string, versioning bool) string {
	return fmt.Sprintf(`
resource "thalassa_objectstorage_bucket" "test" {
  name             = %q
  region           = %q
  versioning       = %t
  wait_for_ready   = true
  wait_for_deleted = true
}
`, name, region, versioning)
}

func testAccBucketLifecycleConfig(bucketName, region, ruleID string, expirationDays int) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_objectstorage_bucket_lifecycle" "test" {
  bucket_name = thalassa_objectstorage_bucket.test.name

  rule {
    id     = %q
    prefix = "logs/"
    status = "Enabled"

    expiration {
      days = %d
    }
  }
}
`, testAccProviderBlock(), testAccBucketBaseConfig(bucketName, region, false), ruleID, expirationDays)
}

func testAccBucketLifecycleNoncurrentConfig(bucketName, region string, noncurrentDays int) string {
	return fmt.Sprintf(`
%s

%s

resource "thalassa_objectstorage_bucket_lifecycle" "test" {
  bucket_name = thalassa_objectstorage_bucket.test.name

  rule {
    id     = "expire-noncurrent"
    prefix = "archive/"
    status = "Enabled"

    noncurrent_version_expiration {
      noncurrent_days = %d
    }
  }
}
`, testAccProviderBlock(), testAccBucketBaseConfig(bucketName, region, true), noncurrentDays)
}
