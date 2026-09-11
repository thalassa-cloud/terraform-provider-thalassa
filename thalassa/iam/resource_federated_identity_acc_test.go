package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFederatedIdentityProvider_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-idp")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityProviderConfig(name, issuer, "active", "acceptance idp"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "name", name),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "provider_issuer", issuer),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "status", "active"),
					resource.TestCheckResourceAttrSet("thalassa_iam_federated_identity_provider.test", "id"),
					resource.TestCheckResourceAttrSet("thalassa_iam_federated_identity_provider.test", "created_at"),
				),
			},
		},
	})
}

func TestAccFederatedIdentityProvider_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-idp")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityProviderConfig(name, issuer, "active", "initial"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "description", "initial"),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "status", "active"),
				),
			},
			{
				Config: testAccFederatedIdentityProviderConfig(name, issuer, "inactive", "updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "description", "updated"),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity_provider.test", "status", "inactive"),
				),
			},
		},
	})
}

func TestAccFederatedIdentityProvider_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-idp")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityProviderConfig(name, issuer, "active", "import me"),
			},
			{
				ResourceName:            "thalassa_iam_federated_identity_provider.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organisation_id", "local_jwks"},
			},
		},
	})
}

func TestAccFederatedIdentityProviderDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-idp")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityProviderDataSourceConfig(name, issuer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_iam_federated_identity_provider.test", "id", "thalassa_iam_federated_identity_provider.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_iam_federated_identity_provider.test", "name", name),
					resource.TestCheckResourceAttr("data.thalassa_iam_federated_identity_provider.test", "provider_issuer", issuer),
				),
			},
		},
	})
}

func TestAccFederatedIdentity_basic(t *testing.T) {
	saName := acctest.RandomWithPrefix("tf-acc-sa")
	idpName := acctest.RandomWithPrefix("tf-acc-idp")
	identityName := acctest.RandomWithPrefix("tf-acc-fi")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))
	subject := fmt.Sprintf("repo:thalassa-cloud/tf-acc:ref:refs/heads/%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityConfig(saName, idpName, identityName, issuer, subject, "active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "name", identityName),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "provider_subject", subject),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "audience_match_mode", "any"),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "status", "active"),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "allowed_scopes.#", "2"),
					resource.TestCheckResourceAttrSet("thalassa_iam_federated_identity.test", "id"),
					resource.TestCheckResourceAttrPair("thalassa_iam_federated_identity.test", "service_account_id", "thalassa_iam_service_account.test", "id"),
					resource.TestCheckResourceAttrPair("thalassa_iam_federated_identity.test", "provider_id", "thalassa_iam_federated_identity_provider.test", "id"),
				),
			},
		},
	})
}

func TestAccFederatedIdentity_update(t *testing.T) {
	saName := acctest.RandomWithPrefix("tf-acc-sa")
	idpName := acctest.RandomWithPrefix("tf-acc-idp")
	identityName := acctest.RandomWithPrefix("tf-acc-fi")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))
	subject := fmt.Sprintf("repo:thalassa-cloud/tf-acc:ref:refs/heads/%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityConfig(saName, idpName, identityName, issuer, subject, "active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "status", "active"),
				),
			},
			{
				Config: testAccFederatedIdentityConfig(saName, idpName, identityName, issuer, subject, "inactive"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "status", "inactive"),
					resource.TestCheckResourceAttr("thalassa_iam_federated_identity.test", "description", "federated identity for acceptance tests"),
				),
			},
		},
	})
}

func TestAccFederatedIdentityDataSource_byName(t *testing.T) {
	saName := acctest.RandomWithPrefix("tf-acc-sa")
	idpName := acctest.RandomWithPrefix("tf-acc-idp")
	identityName := acctest.RandomWithPrefix("tf-acc-fi")
	issuer := fmt.Sprintf("https://token.actions.githubusercontent.com/%s", acctest.RandString(8))
	subject := fmt.Sprintf("repo:thalassa-cloud/tf-acc:ref:refs/heads/%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedIdentityDataSourceConfig(saName, idpName, identityName, issuer, subject),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.thalassa_iam_federated_identity.test", "id", "thalassa_iam_federated_identity.test", "id"),
					resource.TestCheckResourceAttr("data.thalassa_iam_federated_identity.test", "name", identityName),
					resource.TestCheckResourceAttr("data.thalassa_iam_federated_identity.test", "provider_subject", subject),
				),
			},
		},
	})
}

func testAccFederatedIdentityProviderConfig(name, issuer, status, description string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_iam_federated_identity_provider" "test" {
  name            = %q
  description     = %q
  provider_issuer = %q
  status          = %q

  labels = {
    environment = "test"
  }
}
`, testAccProviderBlock(), name, description, issuer, status)
}

func testAccFederatedIdentityProviderDataSourceConfig(name, issuer string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_iam_federated_identity_provider" "test" {
  name            = %q
  provider_issuer = %q
  status          = "active"
}

data "thalassa_iam_federated_identity_provider" "test" {
  name = thalassa_iam_federated_identity_provider.test.name
}
`, testAccProviderBlock(), name, issuer)
}

func testAccFederatedIdentityConfig(saName, idpName, identityName, issuer, subject, status string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_iam_service_account" "test" {
  name        = %q
  description = "acceptance service account"
}

resource "thalassa_iam_federated_identity_provider" "test" {
  name            = %q
  provider_issuer = %q
  status          = "active"
}

resource "thalassa_iam_federated_identity" "test" {
  name                     = %q
  description              = "federated identity for acceptance tests"
  service_account_id = thalassa_iam_service_account.test.id
  provider_id        = thalassa_iam_federated_identity_provider.test.id
  provider_subject         = %q
  trusted_audiences        = ["https://api.thalassa.cloud"]
  audience_match_mode      = "any"
  allowed_scopes = [
    "api:read",
    "api:write",
  ]
  status = %q
}
`, testAccProviderBlock(), saName, idpName, issuer, identityName, subject, status)
}

func testAccFederatedIdentityDataSourceConfig(saName, idpName, identityName, issuer, subject string) string {
	return fmt.Sprintf(`
%s

resource "thalassa_iam_service_account" "test" {
  name = %q
}

resource "thalassa_iam_federated_identity_provider" "test" {
  name            = %q
  provider_issuer = %q
  status          = "active"
}

resource "thalassa_iam_federated_identity" "test" {
  name                     = %q
  service_account_id = thalassa_iam_service_account.test.id
  provider_id        = thalassa_iam_federated_identity_provider.test.id
  provider_subject         = %q
  trusted_audiences        = ["https://api.thalassa.cloud"]
  audience_match_mode      = "any"
  allowed_scopes           = ["api:read"]
  status                   = "active"
}

data "thalassa_iam_federated_identity" "test" {
  name = thalassa_iam_federated_identity.test.name
}
`, testAccProviderBlock(), saName, idpName, issuer, identityName, subject)
}
