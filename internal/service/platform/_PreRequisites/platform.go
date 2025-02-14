package _PreRequisites

import (
	"fmt"
	"path/filepath"
)

// Helper functions for Creation of Resources
func createConnectorVault_app_role(id string, name string, vault_secret string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_text" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]

		secret_manager_identifier = "harnessSecretManager"
		value_type = "Inline"
		value = "%[3]s"
	}

	resource "harness_platform_connector_vault" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
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

		depends_on = [time_sleep.wait_8_seconds]
	}

	resource "time_sleep" "wait_8_seconds" {
		depends_on = [harness_platform_secret_text.test]
		create_duration = "8s"
	}
	`, id, name, vault_secret)
}

func createSecretFile(id string, name string) string {
	return fmt.Sprintf(`
	resource "harness_platform_secret_file" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		description = "test"
		tags = ["foo:bar"]
		file_path = "%[3]s"
		secret_manager_identifier = "harnessSecretManager"
	}
		`, id, name, getAbsFilePath("../../../acctest/secret_files/secret.txt"))
}

func getAbsFilePath(file_path string) string {
	absPath, _ := filepath.Abs(file_path)
	return absPath
}

func createServiceAccount(id string, name string, accountId string) string {
	return fmt.Sprintf(`
	resource "harness_platform_service_account" "test" {
		identifier = "%[1]s"
		name = "%[2]s"
		email = "email@service.harness.io"
		description = "test"
		tags = ["foo:bar"]
		account_id = "%[3]s"
	}
	`, id, name, accountId)
}

func createUserGroup(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_usergroup" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			linked_sso_id = "linked_sso_id"
			externally_managed = false
			users = []
			notification_configs {
				type = "SLACK"
				slack_webhook_url = "https://google.com"
			}
			notification_configs {
				type = "EMAIL"
				group_email = "email@email.com"
				send_email_to_all_users = true
			}
			notification_configs {
				type = "MSTEAMS"
				microsoft_teams_webhook_url = "https://google.com"
			}
			notification_configs {
				type = "PAGERDUTY"
				pager_duty_key = "pagerDutyKey"
			}
			linked_sso_display_name = "linked_sso_display_name"
			sso_group_id = "sso_group_id"
			sso_group_name = "sso_group_name"
			linked_sso_type = "SAML"
			sso_linked = true
		}
`, id, name)
}

func createProject(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
		}

		resource "harness_platform_project" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			org_id = harness_platform_organization.test.id
		}
`, id, name)
}

func createOrganization(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			description = "test"
			tags = ["foo:bar", "baz:qux"]
		}
`, id, name)
}

func createSecretText_inline(id string, name string, secretValue string) string {
	return fmt.Sprintf(`
		resource "harness_platform_secret_text" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			description = "test"
			tags = ["foo:bar"]
			secret_manager_identifier = "harnessSecretManager"
			value_type = "Inline"
			value = "%[3]s"
		}
`, id, name, secretValue)
}

// Add more similar functions for other resources
