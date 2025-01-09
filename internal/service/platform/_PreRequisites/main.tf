# main.tf

terraform {
  required_providers {
    harness = {
      source = "harness/harness"
    }
  }
}

resource "harness_platform_secret_text" "azuretest" {
  identifier                = "azuretest"
  name                      = "azuretest"
  description               = "This is a test secret text"
  tags                      = ["foo:bar"]
  secret_manager_identifier = "harnessSecretManager"
  value_type                = "Inline"
  value                     = "my_secret_value"
}

resource "harness_platform_secret_text" "doNotDeleteHSM" {
  identifier                = "doNotDeleteHSM"
  name                      = "doNotDeleteHSM"
  description               = "This is a test secret text"
  tags                      = ["foo:bar"]
  secret_manager_identifier = "harnessSecretManager"
  value_type                = "Inline"
  value                     = "Harness@123"
}

resource "harness_platform_connector_vault" "my_vault_connector" {
  identifier  = "my_vault_connector_id"
  name        = "My Vault Connector"
  description = "Vault Connector example"
  tags        = ["foo:bar"]

  app_role_id                       = "570acf09-ef2a-144b-2fb0-14a42e06ffe3"
  base_path                         = "vikas-test/"
  access_type                       = "APP_ROLE"
  default                           = false
  secret_id                         = "account.${harness_platform_secret_text.azuretest.id}"
  read_only                         = true
  renewal_interval_minutes          = 60
  secret_engine_manually_configured = true
  secret_engine_name                = "harness-test"
  secret_engine_version             = 2
  use_aws_iam                       = false
  use_k8s_auth                      = false
  use_vault_agent                   = false
  delegate_selectors                = ["harness-delegate"]
  vault_url                         = "https://vaultqa.harness.io"
  use_jwt_auth                      = false

  depends_on = [time_sleep.wait_8_seconds]
}

resource "harness_platform_service_account" "my_service_account" {
  identifier  = "my_service_account_id"
  name        = "My Service Account"
  email       = "email@service.harness.io"
  description = "This is a test service account"
  tags        = ["foo:bar"]
  account_id  = "your_harness_account_id"
}

resource "harness_platform_usergroup" "my_user_group" {
  identifier = "my_user_group_id"
  name       = "My User Group"

  linked_sso_id      = "linked_sso_id"
  externally_managed = false
  users              = []

  notification_configs {
    type              = "SLACK"
    slack_webhook_url = "https://slack.webhook.url"
  }

  notification_configs {
    type                    = "EMAIL"
    group_email             = "email@domain.com"
    send_email_to_all_users = true
  }

  notification_configs {
    type                        = "MSTEAMS"
    microsoft_teams_webhook_url = "https://msteams.webhook.url"
  }

  notification_configs {
    type           = "PAGERDUTY"
    pager_duty_key = "pagerDutyKey"
  }

  linked_sso_display_name = "SSO Display Name"
  sso_group_id            = "sso_group_id"
  sso_group_name          = "sso_group_name"
  linked_sso_type         = "SAML"
  sso_linked              = true
}

resource "harness_platform_organization" "my_organization" {
  identifier  = "my_organization_id"
  name        = "My Organization"
  description = "This is a test organization"
  tags        = ["foo:bar", "baz:qux"]
}

resource "harness_platform_project" "my_project" {
  identifier = "my_project_id"
  name       = "My Project"
  org_id     = harness_platform_organization.my_organization.id
}

resource "time_sleep" "wait_8_seconds" {
  create_duration = "8s"
}

# Additional resources can be added as needed
