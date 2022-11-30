resource "harness_platform_gitops_cluster" "example" {
  identifier = "identifier"
  account_id = "account_id"
  project_id = "project_id"
  org_id     = "org_id"
  agent_id   = "agent_id"

  request {
    upsert = false
    cluster {
      server = "https://kubernetes.default.svc"
      name   = "name"
      config {
        tls_client_config {
          insecure = true
        }
        cluster_connection_type = "IN_CLUSTER"
      }

    }
  }
  lifecycle {
    ignore_changes = [
      request.0.upsert, request.0.cluster.0.config.0.bearer_token,
    ]
  }
}
