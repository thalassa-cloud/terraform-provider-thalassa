package containerregistry_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContainerRegistryNamespace_basic(t *testing.T) {
	namespace := testAccNamespaceName(acctest.RandomWithPrefix("tfacc"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig(namespace, region, "acceptance namespace"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace.test", "namespace", namespace),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace.test", "description", "acceptance namespace"),
					resource.TestCheckResourceAttrSet("thalassa_containerregistry_namespace.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_containerregistry_namespace.test", "created_at"),
				),
			},
		},
	})
}

func TestAccContainerRegistryNamespace_update(t *testing.T) {
	namespace := testAccNamespaceName(acctest.RandomWithPrefix("tfacc"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig(namespace, region, "initial"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace.test", "description", "initial"),
				),
			},
			{
				Config: testAccNamespaceConfig(namespace, region, "updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace.test", "description", "updated"),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace.test", "labels.environment", "test"),
				),
			},
		},
	})
}

func TestAccContainerRegistryNamespace_import(t *testing.T) {
	namespace := testAccNamespaceName(acctest.RandomWithPrefix("tfacc"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig(namespace, region, "import me"),
			},
			{
				ResourceName:            "thalassa_containerregistry_namespace.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organisation_id"},
			},
		},
	})
}

func TestAccContainerRegistryNamespaceDataSource_byName(t *testing.T) {
	namespace := testAccNamespaceName(acctest.RandomWithPrefix("tfacc"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceDataSourceConfig(namespace, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_containerregistry_namespace.test", "id", "thalassa_containerregistry_namespace.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_containerregistry_namespace.test", "namespace", namespace),
					resource.TestCheckResourceAttr("data.thalassa_containerregistry_namespace.test", "region", region),
				),
			},
		},
	})
}

func TestAccContainerRegistryNamespaceConfiguration_basic(t *testing.T) {
	namespace := testAccNamespaceName(acctest.RandomWithPrefix("tfacc"))
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfigurationConfig(namespace, region, false, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace_configuration.test", "visibility", "private"),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace_configuration.test", "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttrPair("thalassa_containerregistry_namespace_configuration.test", "namespace_id", "thalassa_containerregistry_namespace.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_containerregistry_namespace_configuration.test", "id"),
				),
			},
			{
				Config: testAccNamespaceConfigurationConfig(namespace, region, true, 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace_configuration.test", "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace_configuration.test", "retention_policy.0.delete_untagged_images", "true"),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace_configuration.test", "retention_policy.0.rules.0.days", "30"),
					resource.TestCheckResourceAttr("thalassa_containerregistry_namespace_configuration.test", "retention_policy.0.rules.0.scope", "tags"),
				),
			},
		},
	})
}

func testAccNamespaceConfig(namespace, region, description string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_containerregistry_namespace" "test" {
  region      = %q
  namespace   = %q
  description = %q

  labels = {
    environment = "test"
  }
}
`, testAccProviderBlock(), region, namespace, description)
}

func testAccNamespaceDataSourceConfig(namespace, region string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_containerregistry_namespace" "test" {
  region    = %q
  namespace = %q
}

data "thalassa_containerregistry_namespace" "test" {
  region    = thalassa_containerregistry_namespace.test.region
  namespace = thalassa_containerregistry_namespace.test.namespace
}
`, testAccProviderBlock(), region, namespace)
}

func testAccNamespaceConfigurationConfig(namespace, region string, retentionEnabled bool, days int) string {
	retentionBlock := `
  retention_policy {
    enabled                = false
    delete_untagged_images = false
  }
`
	if retentionEnabled {
		retentionBlock = fmt.Sprintf(`
  retention_policy {
    enabled                = true
    delete_untagged_images = true

    rules {
      days  = %d
      scope = "tags"
      tag_patterns = ["v*"]
    }
  }
`, days)
	}

	return fmt.Sprintf(`
%s

resource "thalassa_containerregistry_namespace" "test" {
  region    = %q
  namespace = %q
}

resource "thalassa_containerregistry_namespace_configuration" "test" {
  namespace_id = thalassa_containerregistry_namespace.test.id
  visibility   = "private"
%s
}
`, testAccProviderBlock(), region, namespace, retentionBlock)
}
