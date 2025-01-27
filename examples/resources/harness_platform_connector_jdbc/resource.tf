resource "harness_platform_connector_jdbc" "test" {
  identifier         = "identifer"
  name               = "name"
  description        = "test"
  tags               = ["foo:bar"]
  url                = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
  delegate_selectors = ["harness-delegate"]
  credentials {
    auth_type = "ServiceAccount"
    service_account {
      token_ref = "account.secret_id"
    }
  }
}

resource "harness_platform_connector_jdbc" "test" {
  identifier         = "identifer"
  name               = "name"
  description        = "test"
  tags               = ["foo:bar"]
  url                = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
  delegate_selectors = ["harness-delegate"]
  credentials {
    auth_type = "UsernamePassword"
    username_password {
      username     = "admin"
      password_ref = "account.secret_id"
    }
  }
}

resource "harness_platform_connector_jdbc" "test" {
  identifier         = "identifer"
  name               = "name"
  description        = "test"
  tags               = ["foo:bar"]
  url                = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
  delegate_selectors = ["harness-delegate"]
  credentials {
    auth_type = "UsernamePassword"
    username_password {
      username_ref = "account.user_ref"
      password_ref = "account.secret_id"
    }
  }
}

resource "harness_platform_connector_jdbc" "test" {
  identifier         = "identifer"
  name               = "name"
  description        = "test"
  tags               = ["foo:bar"]
  url                = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
  delegate_selectors = ["harness-delegate"]
  credentials {
    auth_type    = "UsernamePassword"
    username     = "admin"
    password_ref = "account.secret_id"
  }
}

resource "harness_platform_connector_jdbc" "test" {
  identifier         = "identifer"
  name               = "name"
  description        = "test"
  tags               = ["foo:bar"]
  url                = "jdbc:sqlserver://1.2.3;trustServerCertificate=true"
  delegate_selectors = ["harness-delegate"]
  credentials {
    auth_type    = "UsernamePassword"
    username_ref = "account.user_ref"
    password_ref = "account.secret_id"
  }
}