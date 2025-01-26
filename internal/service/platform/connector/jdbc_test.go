package connector_test

import (
	"fmt"
	"testing"

	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceConnectorJDBC(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(5))
	name := id
	updatedName := fmt.Sprintf("%s_updated", name)
	resourceName := "harness_platform_connector_jdbc.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccConnectorDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorJDBC(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url", "jdbc:sqlserver://1.2.3;trustServerCertificate=true"),
					resource.TestCheckResourceAttr(resourceName, "delegate_selectors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.password_ref", fmt.Sprintf("account.%s", id)),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.auth_type", "UsernamePassword"),
				),
			},
			{
				Config: testAccResourceConnectorJDBC(id, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url", "jdbc:sqlserver://1.2.3;trustServerCertificate=true"),
					resource.TestCheckResourceAttr(resourceName, "delegate_selectors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.password_ref", fmt.Sprintf("account.%s", id)),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.auth_type", "UsernamePassword"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceConnectorJDBCDefaultAuth(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(5))
	name := id
	updatedName := fmt.Sprintf("%s_updated", name)
	resourceName := "harness_platform_connector_jdbc.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccConnectorDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorJDBCDefaultAuth(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url", "jdbc:sqlserver://1.2.3;trustServerCertificate=true"),
					resource.TestCheckResourceAttr(resourceName, "delegate_selectors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.password_ref", fmt.Sprintf("account.%s", id)),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.auth_type", "UsernamePassword"),
				),
			},
			{
				Config: testAccResourceConnectorJDBCDefaultAuth(id, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url", "jdbc:sqlserver://1.2.3;trustServerCertificate=true"),
					resource.TestCheckResourceAttr(resourceName, "delegate_selectors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.username_password.0.password_ref", fmt.Sprintf("account.%s", id)),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.auth_type", "UsernamePassword"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceConnectorJDBCServiceAccountAuth(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(5))
	name := id
	updatedName := fmt.Sprintf("%s_updated", name)
	resourceName := "harness_platform_connector_jdbc.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccConnectorDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorJDBCServiceAccountAuth(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url", "jdbc:sqlserver://1.2.3;trustServerCertificate=true"),
					resource.TestCheckResourceAttr(resourceName, "delegate_selectors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.service_account.0.token_ref", fmt.Sprintf("account.%s", id)),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.auth_type", "ServiceAccount"),
				),
			},
			{
				Config: testAccResourceConnectorJDBCServiceAccountAuth(id, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url", "jdbc:sqlserver://1.2.3;trustServerCertificate=true"),
					resource.TestCheckResourceAttr(resourceName, "delegate_selectors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.service_account.0.token_ref", fmt.Sprintf("account.%s", id)),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.auth_type", "ServiceAccount"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceConnectorJDBC(id string, name string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "secret"
	}

	resource "harness_platform_connector_jdbc" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		url = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
		delegate_selectors = ["harness-delegate"]
		credentials {
			auth_type = "UsernamePassword"
			username_password {
				username = "admin"
				password_ref = "account.${harness_platform_secret_text.test.id}"
			}
		}
		depends_on = [harness_platform_secret_text.test]
	}
`, id, name)
}

func testAccResourceConnectorJDBCDefaultAuth(id string, name string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "secret"
	}

	resource "harness_platform_connector_jdbc" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		url = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
		delegate_selectors = ["harness-delegate"]
		credentials {
			username_password {
				username = "admin"
				password_ref = "account.${harness_platform_secret_text.test.id}"
			}
		}
		depends_on = [harness_platform_secret_text.test]
	}
`, id, name)
}

func testAccResourceConnectorJDBCServiceAccountAuth(id string, name string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "secret"
	}

	resource "harness_platform_connector_jdbc" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		url = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
		delegate_selectors = ["harness-delegate"]
		credentials {
			auth_type = "ServiceAccount"
			service_account {
				token_ref = "account.${harness_platform_secret_text.test.id}"
			}
		}
		depends_on = [harness_platform_secret_text.test]
	}
`, id, name)
}
