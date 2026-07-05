package iaas_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa"
)

var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = thalassa.Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"thalassa": func() (*schema.Provider, error) {
			return thalassa.Provider(), nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; set TF_ACC=1 to run")
	}

	if !testAccHasAuth() {
		t.Fatal("acceptance tests require one of: THALASSA_API_TOKEN, THALASSA_ACCESS_TOKEN, or THALASSA_CLIENT_ID and THALASSA_CLIENT_SECRET")
	}

	if os.Getenv("THALASSA_ORGANISATION") == "" {
		t.Fatal("THALASSA_ORGANISATION must be set for acceptance tests")
	}
}

func testAccHasAuth() bool {
	if os.Getenv("THALASSA_API_TOKEN") != "" || os.Getenv("THALASSA_ACCESS_TOKEN") != "" {
		return true
	}

	return os.Getenv("THALASSA_CLIENT_ID") != "" && os.Getenv("THALASSA_CLIENT_SECRET") != ""
}

func testAccRegion() string {
	if region := os.Getenv("THALASSA_TEST_REGION"); region != "" {
		return region
	}

	return "nl-01"
}

func testAccVolumeType() string {
	if volumeType := os.Getenv("THALASSA_TEST_VOLUME_TYPE"); volumeType != "" {
		return volumeType
	}

	return "Block"
}
