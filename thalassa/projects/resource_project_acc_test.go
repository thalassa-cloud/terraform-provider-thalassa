package projects_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProject_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-project")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig(name, "acceptance project"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_project.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_project.test", "description", "acceptance project"),
					resource.TestCheckResourceAttrSet("thalassa_project.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_project.test", "slug"),
					resource.TestCheckResourceAttrSet("thalassa_project.test", "created_at"),
				),
			},
		},
	})
}

func TestAccProject_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-project")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig(name, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_project.test", "description", "initial description"),
				),
			},
			{
				Config: testAccProjectConfig(name, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_project.test", "description", "updated description"),
					resource.TestCheckResourceAttr("thalassa_project.test", "labels.environment", "test"),
				),
			},
		},
	})
}

func TestAccProject_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-project")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig(name, "import me"),
			},
			{
				ResourceName:      "thalassa_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-project")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_project.test", "id", "thalassa_project.test", "id"),
					resource.TestCheckResourceAttrPair("data.thalassa_project.test", "slug", "thalassa_project.test", "slug"),
					resource.TestCheckResourceAttr("data.thalassa_project.test", "name", name),
				),
			},
		},
	})
}

func testAccProjectConfig(name, description string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_project" "test" {
  name        = %q
  description = %q

  labels = {
    environment = "test"
  }
}
`, testAccProviderBlock(), name, description)
}

func testAccProjectDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_project" "test" {
  name        = %q
  description = "lookup"
}

data "thalassa_project" "test" {
  name = thalassa_project.test.name
}
`, testAccProviderBlock(), name)
}
