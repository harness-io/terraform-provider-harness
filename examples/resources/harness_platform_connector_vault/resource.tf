
terraform {
  required_providers {
    harness = {
      source = "harness/harness"
       version = ">= 0.34.0"
    }
    
  }
}


resource "harness_platform_connector_vault" "aws_auth" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  aws_region                        = "aws_region"
  base_path                         = "base_path"
  access_type                       = "AWS_IAM"
  default                           = false
  xvault_aws_iam_server_id          = "account.${harness_platform_secret_text.test.id}"
  read_only                         = true
  renewal_interval_minutes          = 60
  secret_engine_manually_configured = true
  secret_engine_name                = "secret_engine_name"
  secret_engine_version             = 2
  vault_aws_iam_role                = "vault_aws_iam_role"
  use_aws_iam                       = true
  use_k8s_auth                      = false
  use_vault_agent                   = false
  delegate_selectors                = ["harness-delegate"]
  vault_url                         = "https://vault_url.com"
  use_jwt_auth                      = false
}

resource "harness_platform_connector_vault" "app_role" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  app_role_id                       = "app_role_id"
  base_path                         = "base_path"
  access_type                       = "APP_ROLE"
  default                           = false
  secret_id                         = "account.${harness_platform_secret_text.test.id}"
  read_only                         = true
  renewal_interval_minutes          = 60
  secret_engine_manually_configured = true
  secret_engine_name                = "secret_engine_name"
  secret_engine_version             = 2
  use_aws_iam                       = false
  use_k8s_auth                      = false
  use_vault_agent                   = false
  renew_app_role_token              = true
  delegate_selectors                = ["harness-delegate"]
  vault_url                         = "https://vault_url.com"
  use_jwt_auth                      = false
}

resource "harness_platform_connector_vault" "k8s_auth" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  auth_token                        = "account.${harness_platform_secret_text.test.id}"
  base_path                         = "base_path"
  access_type                       = "K8s_AUTH"
  default                           = false
  k8s_auth_endpoint                 = "k8s_auth_endpoint"
  namespace                         = "namespace"
  read_only                         = true
  renewal_interval_minutes          = 10
  secret_engine_manually_configured = true
  secret_engine_name                = "secret_engine_name"
  secret_engine_version             = 2
  service_account_token_path        = "service_account_token_path"
  use_aws_iam                       = false
  use_k8s_auth                      = true
  use_vault_agent                   = false
  vault_k8s_auth_role               = "vault_k8s_auth_role"
  vault_aws_iam_role                = "vault_aws_iam_role"
  delegate_selectors                = ["harness-delegate"]
  vault_url                         = "https://vault_url.com"
  use_jwt_auth                      = false
}

resource "harness_platform_connector_vault" "vault_agent" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  auth_token                        = "account.${harness_platform_secret_text.test.id}"
  base_path                         = "base_path"
  access_type                       = "VAULT_AGENT"
  default                           = false
  namespace                         = "namespace"
  read_only                         = true
  renewal_interval_minutes          = 10
  secret_engine_manually_configured = true
  secret_engine_name                = "secret_engine_name"
  secret_engine_version             = 2
  use_aws_iam                       = false
  use_k8s_auth                      = false
  use_vault_agent                   = true
  sink_path                         = "sink_path"
  delegate_selectors                = ["harness-delegate"]
  vault_url                         = "https://vault_url.com"
  use_jwt_auth                      = false
}



resource "harness_platform_connector_vault" "token" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  auth_token                        = "account.${harness_platform_secret_text.test.id}"
  base_path                         = "base_path"
  access_type                       = "TOKEN"
  default                           = false
  namespace                         = "namespace"
  read_only                         = true
  renewal_interval_minutes          = 10
  secret_engine_manually_configured = true
  secret_engine_name                = "secret_engine_name"
  secret_engine_version             = 2
  use_aws_iam                       = false
  use_k8s_auth                      = false
  vault_url                         = "https://vault_url.com"
  use_jwt_auth                      = false
}

resource "harness_platform_connector_vault" "jwt" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  base_path                         = "base_path"
  access_type                       = "JWT"
  default                           = false
  read_only                         = true
  renewal_interval_minutes          = 60
  secret_engine_manually_configured = true
  secret_engine_name                = "secret_engine_name"
  secret_engine_version             = 2
  use_aws_iam                       = false
  use_k8s_auth                      = false
  use_vault_agent                   = false
  renew_app_role_token              = false
  delegate_selectors                = ["harness-delegate"]
  vault_url                         = "https://vault_url.com"
  use_jwt_auth                      = true
  vault_jwt_auth_role               = "vault_jwt_auth_role"
  vault_jwt_auth_path               = "vault_jwt_auth_path"
  execute_on_delegate               = true
}