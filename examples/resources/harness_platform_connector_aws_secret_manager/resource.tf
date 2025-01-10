# Credentials inherit_from_delegate
resource "harness_platform_connector_aws_secret_manager" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  secret_name_prefix = "test"
  region             = "us-east-1"
  delegate_selectors = ["harness-delegate"]
  use_put_secret     = false
  credentials {
    inherit_from_delegate = true
  }
}

# Credentials manual
resource "harness_platform_connector_aws_secret_manager" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  secret_name_prefix = "test"
  region             = "us-east-1"
  delegate_selectors = ["harness-delegate"]
  use_put_secret     = false
  credentials {
    manual {
      secret_key_ref = "account.secret_id"
      access_key_ref = "account.secret_id"
    }
  }
}

# Credentials assume_role
resource "harness_platform_connector_aws_secret_manager" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  secret_name_prefix = "test"
  region             = "us-east-1"
  delegate_selectors = ["harness-delegate"]
  default            = true
  use_put_secret     = false
  credentials {
    assume_role {
      role_arn    = "somerolearn"
      external_id = "externalid"
      duration    = 900
    }
  }
}


# Force delete true
resource "harness_platform_connector_aws_secret_manager" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  secret_name_prefix = "test"
  region             = "us-east-1"
  delegate_selectors = ["harness-delegate"]
  default            = true
  force_delete_without_recovery     = true
  credentials {
    assume_role {
      role_arn    = "somerolearn"
      external_id = "externalid"
      duration    = 900
    }
  }
}

# With recovery duration of 15 days
resource "harness_platform_connector_aws_secret_manager" "test" {
  identifier  = "identifier"
  name        = "name"
  description = "test"
  tags        = ["foo:bar"]

  secret_name_prefix = "test"
  region             = "us-east-1"
  delegate_selectors = ["harness-delegate"]
  default            = true
  recovery_window_in_days     = 15
  credentials {
    assume_role {
      role_arn    = "somerolearn"
      external_id = "externalid"
      duration    = 900
    }
  }
}