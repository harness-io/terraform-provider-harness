package secretManagers_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// VAULT TOKEN
func TestAccDataSourceConnectorVault(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "harness"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "TOKEN"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVaultProjectLevel(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVaultProjectLevel(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "harness"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "TOKEN"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVaultOrgLevel(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVaultOrgLevel(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "harness"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "TOKEN"),
				),
			},
		},
	})
}

// Vault Agent
func TestAccDataSourceConnectorVault_VaultAgent(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_VaultAgent(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "base_path"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "VAULT_AGENT"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_VaultAgentProjectLevel(t *testing.T) {
	var (
		name          = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName  = "data.harness_platform_connector_vault.test"
		connectorName = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(5))
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_VaultAgentProjectLevel(name, connectorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "base_path"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "VAULT_AGENT"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_VaultAgentOrgLevel(t *testing.T) {
	var (
		name          = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName  = "data.harness_platform_connector_vault.test"
		connectorName = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(5))
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_VaultAgentOrgLevel(name, connectorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", connectorName),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "base_path"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "VAULT_AGENT"),
				),
			},
		},
	})
}

// k8sAuth
func TestAccDataSourceConnectorVault_k8sAuth(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_k8sAuth(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "base_path"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "K8s_AUTH"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_k8sAuthProjectLevel(t *testing.T) {
	var (
		name          = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		connectorName = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName  = "data.harness_platform_connector_vault.test"
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_k8sAuthProjectLevel(name, connectorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "base_path"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "K8s_AUTH"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_k8sAuthOrgLevel(t *testing.T) {
	var (
		name          = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName  = "data.harness_platform_connector_vault.test"
		connectorName = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(5))
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_k8sAuthOrgLevel(name, connectorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "base_path"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "K8s_AUTH"),
				),
			},
		},
	})
}

// App Role
func TestAccDataSourceConnectorVault_AppRole(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_AppRole(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "vikas-test/"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "APP_ROLE"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_AppRoleProjectLevel(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_AppRoleProjectLevel(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "vikas-test/"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "APP_ROLE"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_AppRoleOrgLevel(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_AppRoleOrgLevel(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "vikas-test/"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "APP_ROLE"),
				),
			},
		},
	})
}

func TestAccDataSourceConnectorVault_AWSAuth(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_AWSAuth(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "vikas-test/"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_vault_agent", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "AWS_IAM"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_AWSAuthProjectLevel(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_AWSAuthProjectLevel(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "vikas-test/"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_vault_agent", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "AWS_IAM"),
				),
			},
		},
	})
}
func TestAccDataSourceConnectorVault_AWSAuthOrgLevel(t *testing.T) {
	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "data.harness_platform_connector_vault.test"
		vaultToken   = os.Getenv("HARNESS_TEST_VAULT_SECRET")
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectorVault_AWSAuthProjectLevel(name, vaultToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "base_path", "vikas-test/"),
					resource.TestCheckResourceAttr(resourceName, "is_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "renewal_interval_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "secret_engine_manually_configured", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_vault_agent", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "AWS_IAM"),
				),
			},
		},
	})
}

func testAccDataSourceConnectorVault(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		auth_token = "account.${harness_platform_secret_text.test.id}"
		base_path = "harness"
		access_type = "TOKEN"
		default = false
		read_only = true
		renewal_interval_minutes = 0
		secret_engine_manually_configured = true
		secret_engine_name = "QA_Secrets"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false

		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		destroy_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id

	}
`, name, vaultToken)
}
func testAccDataSourceConnectorVaultProjectLevel(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	
	resource "harness_platform_project" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		org_id = harness_platform_organization.test.id
		color = "#472848"
	}

	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_project.test]
		create_duration = "4s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		auth_token = "${harness_platform_secret_text.test.id}"
		base_path = "harness"
		access_type = "TOKEN"
		default = false
		read_only = true
		renewal_interval_minutes = 0
		secret_engine_manually_configured = true
		secret_engine_name = "QA_Secrets"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false

		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "3s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
		project_id = harness_platform_connector_vault.test.project_id
	}
`, name, vaultToken)
}
func testAccDataSourceConnectorVaultOrgLevel(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_organization.test]
		create_duration = "3s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		auth_token = "org.${harness_platform_secret_text.test.id}"
		base_path = "harness"
		access_type = "TOKEN"
		default = false
		read_only = true
		renewal_interval_minutes = 0
		secret_engine_manually_configured = true
		secret_engine_name = "QA_Secrets"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
	}
`, name, vaultToken)
}

func testAccDataSourceConnectorVault_VaultAgent(name string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Reference"
		value = "secret"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		auth_token = "account.${harness_platform_secret_text.test.id}"
		base_path = "base_path"
		access_type = "VAULT_AGENT"
		default = false
		namespace = "namespace"
		read_only = true
		renewal_interval_minutes = 10
		secret_engine_manually_configured = true
		secret_engine_name = "secret_engine_name"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		use_vault_agent = true
		sink_path = "sink_path"
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vault_url.com"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id

	}
`, name)
}
func testAccDataSourceConnectorVault_VaultAgentProjectLevel(name string, connectorName string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	
	resource "harness_platform_project" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		org_id = harness_platform_organization.test.id
		color = "#472848"
	}

	resource "harness_platform_connector_azure_key_vault" "test" {
		identifier = "%[2]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		client_id = "38fca8d7-4dda-41d5-b106-e5d8712b733a"
		secret_key = "account.azuretest"
		tenant_id = "b229b2bb-5f33-4d22-bce0-730f6474e906"
		vault_name = "Aman-test"
		subscription = "20d6a917-99fa-4b1b-9b2e-a3d624e9dcf0"
		is_default = false

		azure_environment_type = "AZURE"
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_project.test]
		create_duration = "3s"
	}

	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Reference"
		value = "secret"
		depends_on = [time_sleep.wait_5_seconds]
	}

	resource "time_sleep" "wait_5_seconds" {
		depends_on = [harness_platform_connector_azure_key_vault.test]
		create_duration = "5s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		auth_token = "${harness_platform_secret_text.test.id}"
		base_path = "base_path"
		access_type = "VAULT_AGENT"
		default = false
		namespace = "namespace"
		read_only = true
		renewal_interval_minutes = 10
		secret_engine_manually_configured = true
		secret_engine_name = "secret_engine_name"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		use_vault_agent = true
		sink_path = "sink_path"
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vault_url.com"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
		project_id = harness_platform_connector_vault.test.project_id
	}
`, name, connectorName)
}
func testAccDataSourceConnectorVault_VaultAgentOrgLevel(name string, connectorName string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	resource "harness_platform_connector_azure_key_vault" "test" {
		identifier = "%[2]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		client_id = "38fca8d7-4dda-41d5-b106-e5d8712b733a"
		secret_key = "account.azuretest"
		tenant_id = "b229b2bb-5f33-4d22-bce0-730f6474e906"
		vault_name = "Aman-test"
		subscription = "20d6a917-99fa-4b1b-9b2e-a3d624e9dcf0"
		is_default = false

		azure_environment_type = "AZURE"
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_organization.test]
		create_duration = "3s"
	}

	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id = harness_platform_organization.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Reference"
		value = "secret"
		depends_on = [time_sleep.wait_2_seconds]
	}

	resource "time_sleep" "wait_2_seconds" {
		depends_on = [harness_platform_connector_azure_key_vault.test]
		create_duration = "2s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		auth_token = "org.${harness_platform_secret_text.test.id}"
		base_path = "base_path"
		access_type = "VAULT_AGENT"
		default = false
		namespace = "namespace"
		read_only = true
		renewal_interval_minutes = 10
		secret_engine_manually_configured = true
		secret_engine_name = "secret_engine_name"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		use_vault_agent = true
		sink_path = "sink_path"
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vault_url.com"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
	}
`, name, connectorName)
}

func testAccDataSourceConnectorVault_k8sAuth(name string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Reference"
		value = "secret"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		auth_token = "account.${harness_platform_secret_text.test.id}"
		base_path = "base_path"
		access_type = "K8s_AUTH"
		default = false
		k8s_auth_endpoint = "k8s_auth_endpoint"
		namespace = "namespace"
		read_only = true
		renewal_interval_minutes = 10
		secret_engine_manually_configured = true
		secret_engine_name = "secret_engine_name"
		secret_engine_version = 2
		service_account_token_path = "service_account_token_path"
		use_aws_iam = false
		use_k8s_auth = true
		use_vault_agent = false
		vault_k8s_auth_role = "vault_k8s_auth_role"
		vault_aws_iam_role = "vault_aws_iam_role"
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vault_url.com"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id

	}
`, name)
}
func testAccDataSourceConnectorVault_k8sAuthProjectLevel(name string, connectorName string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	
	resource "harness_platform_project" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		org_id = harness_platform_organization.test.id
		color = "#472848"
	}

	resource "harness_platform_connector_azure_key_vault" "test" {
		identifier = "%[2]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		client_id = "38fca8d7-4dda-41d5-b106-e5d8712b733a"
		secret_key = "account.azuretest"
		tenant_id = "b229b2bb-5f33-4d22-bce0-730f6474e906"
		vault_name = "Aman-test"
		subscription = "20d6a917-99fa-4b1b-9b2e-a3d624e9dcf0"
		is_default = false

		azure_environment_type = "AZURE"
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_project.test]
		create_duration = "3s"
	}
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Reference"
		value = "secret"
		depends_on = [time_sleep.wait_5_seconds]
	}

	resource "time_sleep" "wait_5_seconds" {
		depends_on = [harness_platform_connector_azure_key_vault.test]
		create_duration = "5s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		project_id=harness_platform_project.test.id
		org_id= harness_platform_organization.test.id
		auth_token = "%[1]s"
		base_path = "base_path"
		access_type = "K8s_AUTH"
		default = false
		k8s_auth_endpoint = "k8s_auth_endpoint"
		namespace = "namespace"
		read_only = true
		renewal_interval_minutes = 10
		secret_engine_manually_configured = true
		secret_engine_name = "secret_engine_name"
		secret_engine_version = 2
		service_account_token_path = "service_account_token_path"
		use_aws_iam = false
		use_k8s_auth = true
		use_vault_agent = false
		vault_k8s_auth_role = "vault_k8s_auth_role"
		vault_aws_iam_role = "vault_aws_iam_role"
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vault_url.com"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
		project_id = harness_platform_connector_vault.test.project_id
	}
`, name, connectorName)
}
func testAccDataSourceConnectorVault_k8sAuthOrgLevel(name string, connectorName string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	resource "harness_platform_connector_azure_key_vault" "test" {
		identifier = "%[2]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		client_id = "38fca8d7-4dda-41d5-b106-e5d8712b733a"
		secret_key = "account.azuretest"
		tenant_id = "b229b2bb-5f33-4d22-bce0-730f6474e906"
		vault_name = "Aman-test"
		subscription = "20d6a917-99fa-4b1b-9b2e-a3d624e9dcf0"
		is_default = false

		azure_environment_type = "AZURE"
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_organization.test]
		create_duration = "3s"
	}
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Reference"
		value = "secret"
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_connector_azure_key_vault.test]
		create_duration = "4s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		auth_token = "org.${harness_platform_secret_text.test.id}"
		base_path = "base_path"
		access_type = "K8s_AUTH"
		default = false
		k8s_auth_endpoint = "k8s_auth_endpoint"
		namespace = "namespace"
		read_only = true
		renewal_interval_minutes = 10
		secret_engine_manually_configured = true
		secret_engine_name = "secret_engine_name"
		secret_engine_version = 2
		service_account_token_path = "service_account_token_path"
		use_aws_iam = false
		use_k8s_auth = true
		use_vault_agent = false
		vault_k8s_auth_role = "vault_k8s_auth_role"
		vault_aws_iam_role = "vault_aws_iam_role"
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vault_url.com"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_5_seconds]
	}

	resource "time_sleep" "wait_5_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "5s"
	}
	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
	
	}
`, name, connectorName)
}

func testAccDataSourceConnectorVault_AppRole(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		app_role_id = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
		base_path = "vikas-test/"
		access_type = "APP_ROLE"
		default = false
		secret_id = "account.${harness_platform_secret_text.test.id}"
		read_only = true
		renewal_interval_minutes = 60
		secret_engine_manually_configured = true
		secret_engine_name = "harness-test"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		use_vault_agent = false
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "15s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_connector_vault.test]
		create_duration = "3s"
	}
`, name, vaultToken)
}
func testAccDataSourceConnectorVault_AppRoleProjectLevel(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	
	resource "harness_platform_project" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		org_id = harness_platform_organization.test.id
		color = "#472848"
		depends_on = [time_sleep.wait_1_seconds]
	}

	resource "time_sleep" "wait_1_seconds" {
		depends_on = [harness_platform_organization.test]
		create_duration = "3s"
	}

	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_project.test]
		create_duration = "10s"
	}
	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		app_role_id = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
		base_path = "vikas-test/"
		access_type = "APP_ROLE"
		default = false
		secret_id = "${harness_platform_secret_text.test.id}"
		read_only = true
		renewal_interval_minutes = 60
		secret_engine_manually_configured = true
		secret_engine_name = "harness-test"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		use_vault_agent = false
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "15s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
		project_id = harness_platform_connector_vault.test.project_id
		depends_on = [time_sleep.wait_10_seconds]
	}

	resource "time_sleep" "wait_10_seconds" {
		depends_on = [harness_platform_connector_vault.test]
		create_duration = "3s"
	}
`, name, vaultToken)
}
func testAccDataSourceConnectorVault_AppRoleOrgLevel(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}

	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_organization.test]
		create_duration = "15s"
	}
	

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		app_role_id = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
		base_path = "vikas-test/"
		access_type = "APP_ROLE"
		default = false
		secret_id = "org.${harness_platform_secret_text.test.id}"
		read_only = true
		renewal_interval_minutes = 60
		secret_engine_manually_configured = true
		secret_engine_name = "harness-test"
		secret_engine_version = 2
		use_aws_iam = false
		use_k8s_auth = false
		use_vault_agent = false
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "15s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
		depends_on = [time_sleep.wait_6_seconds]
	}

	resource "time_sleep" "wait_6_seconds" {
		depends_on = [harness_platform_connector_vault.test]
		create_duration = "6s"
	}
`, name, vaultToken)
}

func testAccDataSourceConnectorVault_AWSAuth(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]

		aws_region = "us-east-2"
		base_path = "vikas-test/"
		access_type = "AWS_IAM"
		default = false
		xvault_aws_iam_server_id = "account.${harness_platform_secret_text.test.id}"
		read_only = true
		renewal_interval_minutes = 60
		secret_engine_manually_configured = true
		secret_engine_name = "harness-test"
		secret_engine_version = 2
		vault_aws_iam_role = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
		use_aws_iam = true
		use_k8s_auth = false
		use_vault_agent = false
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}
	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id

	}
`, name, vaultToken)
}
func testAccDataSourceConnectorVault_AWSAuthProjectLevel(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}
	resource "harness_platform_project" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		org_id = harness_platform_organization.test.id
		color = "#472848"
	}


	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_project.test]
		create_duration = "4s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		project_id=harness_platform_project.test.id
		aws_region = "us-east-2"
		base_path = "vikas-test/"
		access_type = "AWS_IAM"
		default = false
		xvault_aws_iam_server_id = "${harness_platform_secret_text.test.id}"
		read_only = true
		renewal_interval_minutes = 60
		secret_engine_manually_configured = true
		secret_engine_name = "harness-test"
		secret_engine_version = 2
		vault_aws_iam_role = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
		use_aws_iam = true
		use_k8s_auth = false
		use_vault_agent = false
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "3s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
		project_id = harness_platform_connector_vault.test.project_id
	}
`, name, vaultToken)
}
func testAccDataSourceConnectorVault_AWSAuthOrgLevel(name string, vaultToken string) string {
	return fmt.Sprintf(`
	resource "harness_platform_organization" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
	}

	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[2]s"
		depends_on = [time_sleep.wait_3_seconds]
	}

	resource "time_sleep" "wait_3_seconds" {
		depends_on = [harness_platform_organization.test]
		create_duration = "3s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[1]s"
		description = "test"
		tags = ["foo:bar"]
		org_id= harness_platform_organization.test.id
		aws_region = "us-east-2"
		base_path = "vikas-test/"
		access_type = "AWS_IAM"
		default = false
		xvault_aws_iam_server_id = "org.${harness_platform_secret_text.test.id}"
		read_only = true
		renewal_interval_minutes = 60
		secret_engine_manually_configured = true
		secret_engine_name = "harness-test"
		secret_engine_version = 2
		vault_aws_iam_role = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
		use_aws_iam = true
		use_k8s_auth = false
		use_vault_agent = false
		delegate_selectors = ["harness-delegate"]
		vault_url = "https://vaultqa.harness.io"
		use_jwt_auth = false
		depends_on = [time_sleep.wait_4_seconds]
	}

	resource "time_sleep" "wait_4_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "4s"
	}

	data "harness_platform_connector_vault" "test" {
		identifier = harness_platform_connector_vault.test.id
		org_id = harness_platform_connector_vault.test.org_id
	}
`, name, vaultToken)
}
