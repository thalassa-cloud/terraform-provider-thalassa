package iaas_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

var (
	testAccReservedIPOnce    sync.Once
	testAccReservedIPSkipMsg string
)

func testAccPreCheckReservedIP(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)

	testAccReservedIPOnce.Do(func() {
		ctx := context.Background()
		diags := testAccProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]any{}))
		if diags.HasError() {
			testAccReservedIPSkipMsg = "failed to configure provider for reserved IP pre-check"
			return
		}

		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"organisation_id": {Type: schema.TypeString, Optional: true},
		}, map[string]any{})

		client, err := provider.GetClient(provider.GetProvider(testAccProvider.Meta()), rd)
		if err != nil {
			testAccReservedIPSkipMsg = err.Error()
			return
		}

		_, err = client.IaaS().ListReservedIPs(ctx, &iaas.ListReservedIPsRequest{})
		if err != nil && strings.Contains(err.Error(), "status 403") {
			testAccReservedIPSkipMsg = "reserved IPs are not available for this organisation"
		} else if err != nil {
			testAccReservedIPSkipMsg = err.Error()
		}
	})

	if testAccReservedIPSkipMsg != "" {
		if strings.Contains(testAccReservedIPSkipMsg, "not available") {
			t.Skip(testAccReservedIPSkipMsg)
		}
		t.Fatal(testAccReservedIPSkipMsg)
	}
}

func TestAccReservedIP_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-rip")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckReservedIP(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedIPConfig(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "region", region),
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "status", "available"),
					resource.TestCheckResourceAttrSet("thalassa_reserved_ip.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_reserved_ip.test", "slug"),
					testAccCheckReservedIPHasPublicAddress("thalassa_reserved_ip.test"),
				),
			},
		},
	})
}

func TestAccReservedIP_withOptionalAttributes(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-rip")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckReservedIP(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedIPConfigWithOptionalAttributes(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "description", "acceptance test reserved ip"),
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "labels.environment", "test"),
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "annotations.managed_by", "terraform"),
				),
			},
		},
	})
}

func TestAccReservedIP_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-rip")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckReservedIP(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedIPConfigWithDescription(name, region, "initial description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "description", "initial description"),
				),
			},
			{
				Config: testAccReservedIPConfigWithDescription(name, region, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_reserved_ip.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccReservedIP_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-rip")
	region := testAccRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckReservedIP(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedIPConfig(name, region),
			},
			{
				ResourceName:            "thalassa_reserved_ip.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccReservedIPImportStateID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organisation_id"},
			},
		},
	})
}

func testAccReservedIPImportStateID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["thalassa_reserved_ip.test"]
	if !ok {
		return "", fmt.Errorf("resource thalassa_reserved_ip.test not found")
	}
	region := rs.Primary.Attributes["region"]
	if region == "" {
		return "", fmt.Errorf("region attribute is empty")
	}
	return region + "/" + rs.Primary.ID, nil
}

func testAccCheckReservedIPHasPublicAddress(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID", resourceName)
		}

		ipv4 := rs.Primary.Attributes["ipv4_address"]
		ipv6 := rs.Primary.Attributes["ipv6_address"]
		if ipv4 == "" && ipv6 == "" {
			return fmt.Errorf("expected %s to have ipv4_address or ipv6_address set", resourceName)
		}

		return nil
	}
}

func testAccReservedIPConfig(name, region string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_reserved_ip" "test" {
  name   = %q
  region = %q
}
`, testAccProviderBlock(), name, region)
}

func testAccReservedIPConfigWithDescription(name, region, description string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_reserved_ip" "test" {
  name        = %q
  region      = %q
  description = %q
}
`, testAccProviderBlock(), name, region, description)
}

func testAccReservedIPConfigWithOptionalAttributes(name, region string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_reserved_ip" "test" {
  name        = %q
  region      = %q
  description = "acceptance test reserved ip"

  labels = {
    environment = "test"
  }

  annotations = {
    managed_by = "terraform"
  }
}
`, testAccProviderBlock(), name, region)
}
