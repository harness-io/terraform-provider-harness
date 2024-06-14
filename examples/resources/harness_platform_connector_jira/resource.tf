# Auth type UsernamePassword
resource "harness_platform_connector_jira" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  url                = "https://jira.com"
  delegate_selectors = ["harness-delegate"]
  auth {
    auth_type = "UsernamePassword"
    username_password {
      username     = "admin"
      password_ref = "account.secret_id"
    }
  }
}

# Auth Type PersonalAccessToken
resource "harness_platform_connector_jira" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  url                = "https://jira.com"
  delegate_selectors = ["harness-delegate"]
  auth {
    auth_type = "PersonalAccessToken"
    personal_access_token {
      pat_ref = "account.secret_id"
    }
  }
}
