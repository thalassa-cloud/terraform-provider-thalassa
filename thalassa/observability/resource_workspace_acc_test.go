package observability_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccObservabilityWorkspace_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-obs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityWorkspaceConfig(name, region, "acceptance workspace", 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "description", "acceptance workspace"),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "retention_days", "30"),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "status", "ready"),
					resource.TestCheckResourceAttrSet("thalassa_observability_workspace.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_observability_workspace.test", "remote_write_url"),
					resource.TestCheckResourceAttrSet("thalassa_observability_workspace.test", "loki_query_url"),
				),
			},
		},
	})
}

func TestAccObservabilityWorkspace_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-obs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityWorkspaceConfig(name, region, "initial", 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "description", "initial"),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "retention_days", "30"),
				),
			},
			{
				Config: testAccObservabilityWorkspaceConfig(name, region, "updated", 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "description", "updated"),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "retention_days", "60"),
					resource.TestCheckResourceAttr("thalassa_observability_workspace.test", "labels.environment", "test"),
				),
			},
		},
	})
}

func TestAccObservabilityWorkspace_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-obs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityWorkspaceConfig(name, region, "import me", 30),
			},
			{
				ResourceName:            "thalassa_observability_workspace.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccObservabilityWorkspaceImportStateID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organisation_id", "wait_for_deleted", "wait_for_deleted_timeout"},
			},
		},
	})
}

func TestAccObservabilityWorkspaceDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-obs")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityWorkspaceDataSourceConfig(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_observability_workspace.test", "id", "thalassa_observability_workspace.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_observability_workspace.test", "name", name),
					resource.TestCheckResourceAttr("data.thalassa_observability_workspace.test", "region", region),
					resource.TestCheckResourceAttrSet("data.thalassa_observability_workspace.test", "remote_write_url"),
				),
			},
		},
	})
}

func testAccObservabilityWorkspaceImportStateID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["thalassa_observability_workspace.test"]
	if !ok {
		return "", fmt.Errorf("resource thalassa_observability_workspace.test not found")
	}
	region := rs.Primary.Attributes["region"]
	if region == "" {
		return "", fmt.Errorf("region attribute is empty")
	}
	return region + "/" + rs.Primary.ID, nil
}

func testAccObservabilityWorkspaceConfig(name, region, description string, retentionDays int) string {
	return fmt.Sprintf(`
%s

resource "thalassa_observability_workspace" "test" {
  name           = %q
  region         = %q
  description    = %q
  retention_days = %d

  labels = {
    environment = "test"
  }
}
`, testAccProviderBlock(), name, region, description, retentionDays)
}

func testAccObservabilityWorkspaceDataSourceConfig(name, region string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_observability_workspace" "test" {
  name           = %q
  region         = %q
  description    = "lookup"
  retention_days = 30
}

data "thalassa_observability_workspace" "test" {
  name   = thalassa_observability_workspace.test.name
  region = thalassa_observability_workspace.test.region
}
`, testAccProviderBlock(), name, region)
}
