package provider

import (
	"context"
	"fmt"
	"github.com/harness/harness-go-sdk/harness/har"
	"log"

	"github.com/harness/harness-go-sdk/harness/chaos"
	cdng_service "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/service"
	"github.com/harness/terraform-provider-harness/internal/service/chaos/infrastructure"
	har_registry "github.com/harness/terraform-provider-harness/internal/service/har/registry"
	"github.com/harness/terraform-provider-harness/internal/service/platform/module_registry"
	"github.com/harness/terraform-provider-harness/internal/service/platform/service_account"

	pipeline_gitx "github.com/harness/terraform-provider-harness/internal/service/pipeline/gitx/webhook"
	"github.com/harness/terraform-provider-harness/internal/service/platform/cluster_orchestrator"
	dbinstance "github.com/harness/terraform-provider-harness/internal/service/platform/db_instance"
	dbschema "github.com/harness/terraform-provider-harness/internal/service/platform/db_schema"
	governance_enforcement "github.com/harness/terraform-provider-harness/internal/service/platform/governance/enforcement"
	governance_rule "github.com/harness/terraform-provider-harness/internal/service/platform/governance/rule"
	governance_rule_set "github.com/harness/terraform-provider-harness/internal/service/platform/governance/rule_set"
	"github.com/harness/terraform-provider-harness/internal/service/platform/notification_rule"

	cdng_manual_freeze "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/manual_freeze"
	"github.com/harness/terraform-provider-harness/internal/service/platform/feature_flag"
	"github.com/harness/terraform-provider-harness/internal/service/platform/feature_flag_target"
	feature_flag_target_group "github.com/harness/terraform-provider-harness/internal/service/platform/feature_flag_target_group"
	"github.com/harness/terraform-provider-harness/internal/service/platform/ff_api_key"
	"github.com/harness/terraform-provider-harness/internal/service/platform/gitops/agent_yaml"
	"github.com/harness/terraform-provider-harness/internal/service/platform/iacm"
	"github.com/harness/terraform-provider-harness/internal/service/platform/policy"
	"github.com/harness/terraform-provider-harness/internal/service/platform/policyset"
	"github.com/harness/terraform-provider-harness/internal/service/platform/repo_rule_branch"
	"github.com/harness/terraform-provider-harness/internal/service/platform/repo_webhook"
	"github.com/harness/terraform-provider-harness/internal/service/platform/workspace"
	"github.com/sirupsen/logrus"

	"github.com/harness/harness-go-sdk/harness"
	"github.com/harness/harness-go-sdk/harness/cd"
	"github.com/harness/harness-go-sdk/harness/code"
	"github.com/harness/harness-go-sdk/harness/dbops"
	"github.com/harness/harness-go-sdk/harness/helpers"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/harness-go-sdk/harness/utils"
	openapi_client_nextgen "github.com/harness/harness-openapi-go-client/nextgen"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/harness/terraform-provider-harness/internal/service/cd/account"
	"github.com/harness/terraform-provider-harness/internal/service/cd/application"
	"github.com/harness/terraform-provider-harness/internal/service/cd/cloudprovider"
	cd_connector "github.com/harness/terraform-provider-harness/internal/service/cd/connector"
	"github.com/harness/terraform-provider-harness/internal/service/cd/delegate"
	"github.com/harness/terraform-provider-harness/internal/service/cd/environment"
	"github.com/harness/terraform-provider-harness/internal/service/cd/secrets"
	"github.com/harness/terraform-provider-harness/internal/service/cd/service"
	"github.com/harness/terraform-provider-harness/internal/service/cd/sso"
	cd_trigger "github.com/harness/terraform-provider-harness/internal/service/cd/trigger"
	"github.com/harness/terraform-provider-harness/internal/service/cd/user"
	"github.com/harness/terraform-provider-harness/internal/service/cd/yamlconfig"
	cdng_connector_artifactRepositories "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/connector/artifactRepositories"
	cdng_connector_cloudProviders "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/connector/cloudProviders"
	cdng_connector_codeRepositories "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/connector/codeRepositories"
	cdng_environment "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/environment"
	cdng_environment_group "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/environment_group"
	cdng_environment_service_overrides "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/environment_service_overrides"
	cdng_file_store "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/file_store"
	cdng_filters "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/filters"
	cdng_infrastructure "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/infrastructure"
	pl_account "github.com/harness/terraform-provider-harness/internal/service/platform/account"
	pl_apikey "github.com/harness/terraform-provider-harness/internal/service/platform/api_key"
	"github.com/harness/terraform-provider-harness/internal/service/platform/autostopping/load_balancer"
	as_rule "github.com/harness/terraform-provider-harness/internal/service/platform/autostopping/rule"
	"github.com/harness/terraform-provider-harness/internal/service/platform/autostopping/schedule"
	"github.com/harness/terraform-provider-harness/internal/service/platform/ccm_filters"
	"github.com/harness/terraform-provider-harness/internal/service/platform/connector"
	pl_secretManagers "github.com/harness/terraform-provider-harness/internal/service/platform/connector/secretManagers"
	pl_current_user "github.com/harness/terraform-provider-harness/internal/service/platform/current_user"
	pl_delegatetoken "github.com/harness/terraform-provider-harness/internal/service/platform/delegate_token"
	pl_environment_clusters_mapping "github.com/harness/terraform-provider-harness/internal/service/platform/environment_clusters_mapping"
	gitops_agent "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/agent"
	gitops_project_mapping "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/app_project"
	gitops_applications "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/applications"
	gitops_cluster "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/cluster"
	gitops_gnupg "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/gnupg"
	gitops_project "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/project"
	gitops_repository "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/repository"
	gitops_repo_cert "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/repository_certificates"
	gitops_repo_cred "github.com/harness/terraform-provider-harness/internal/service/platform/gitops/repository_credentials"

	cdng_overrides "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/overrides"
	cdng_service_overrides_v2 "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/service_overrides_v2"
	cdng_variables "github.com/harness/terraform-provider-harness/internal/service/cd_nextgen/variables"
	pipeline_input_set "github.com/harness/terraform-provider-harness/internal/service/pipeline/input_set"
	"github.com/harness/terraform-provider-harness/internal/service/pipeline/pipeline"
	"github.com/harness/terraform-provider-harness/internal/service/pipeline/pipeline_filters"
	pipeline_template "github.com/harness/terraform-provider-harness/internal/service/pipeline/template"
	pipeline_template_filters "github.com/harness/terraform-provider-harness/internal/service/pipeline/template_filters"
	pipeline_triggers "github.com/harness/terraform-provider-harness/internal/service/pipeline/triggers"
	"github.com/harness/terraform-provider-harness/internal/service/platform/monitored_service"
	"github.com/harness/terraform-provider-harness/internal/service/platform/organization"
	pl_permissions "github.com/harness/terraform-provider-harness/internal/service/platform/permissions"
	"github.com/harness/terraform-provider-harness/internal/service/platform/project"
	pl_provider "github.com/harness/terraform-provider-harness/internal/service/platform/provider"
	"github.com/harness/terraform-provider-harness/internal/service/platform/repo"
	"github.com/harness/terraform-provider-harness/internal/service/platform/resource_group"
	"github.com/harness/terraform-provider-harness/internal/service/platform/role_assignments"
	"github.com/harness/terraform-provider-harness/internal/service/platform/roles"
	"github.com/harness/terraform-provider-harness/internal/service/platform/secret"
	"github.com/harness/terraform-provider-harness/internal/service/platform/slo"
	pl_token "github.com/harness/terraform-provider-harness/internal/service/platform/token"
	pl_user "github.com/harness/terraform-provider-harness/internal/service/platform/user"
	"github.com/harness/terraform-provider-harness/internal/service/platform/usergroup"

	"github.com/harness/harness-go-sdk/logging"
	openapi_client_logging "github.com/harness/harness-openapi-go-client/logging"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func Provider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"endpoint": {
					Description: fmt.Sprintf("The URL of the Harness API endpoint. The default is `https://app.harness.io/gateway`. This can also be set using the `%s` environment variable.", helpers.EnvVars.Endpoint.String()),
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(helpers.EnvVars.Endpoint.String(), utils.BaseUrl),
				},
				"account_id": {
					Description: fmt.Sprintf("The Harness account id. This can also be set using the `%s` environment variable.", helpers.EnvVars.AccountId.String()),
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(helpers.EnvVars.AccountId.String(), nil),
				},
				"api_key": {
					Description: fmt.Sprintf("The Harness API key. This can also be set using the `%s` environment variable. For more information to create an API key in FirstGen, see https://docs.harness.io/article/smloyragsm-api-keys#create_an_api_key.", helpers.EnvVars.ApiKey.String()),
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(helpers.EnvVars.ApiKey.String(), nil),
				},
				"platform_api_key": {
					Description: fmt.Sprintf("The API key for the Harness next gen platform. This can also be set using the `%s` environment variable. For more information to create an API key in NextGen, see https://docs.harness.io/article/tdoad7xrh9-add-and-manage-api-keys.", helpers.EnvVars.PlatformApiKey.String()),
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(helpers.EnvVars.PlatformApiKey.String(), nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"harness_platform_template":                        pipeline_template.DataSourceTemplate(),
				"harness_platform_connector_azure_key_vault":       pl_secretManagers.DataSourceConnectorAzureKeyVault(),
				"harness_platform_connector_gcp_cloud_cost":        connector.DataSourceConnectorGCPCloudCost(),
				"harness_platform_connector_kubernetes_cloud_cost": connector.DatasourceConnectorKubernetesCloudCost(),
				"harness_platform_connector_azure_cloud_cost":      connector.DataSourceConnectorAzureCloudCost(),
				"harness_platform_connector_appdynamics":           connector.DatasourceConnectorAppDynamics(),
				"harness_platform_connector_elasticsearch":         connector.DatasourceConnectorElasticSearch(),
				"harness_platform_connector_artifactory":           cdng_connector_artifactRepositories.DatasourceConnectorArtifactory(),
				"harness_platform_connector_aws_secret_manager":    pl_secretManagers.DatasourceConnectorAwsSM(),
				"harness_platform_connector_aws":                   cdng_connector_cloudProviders.DatasourceConnectorAws(),
				"harness_platform_connector_awscc":                 connector.DatasourceConnectorAwsCC(),
				"harness_platform_connector_awskms":                pl_secretManagers.DatasourceConnectorAwsKms(),
				"harness_platform_connector_bitbucket":             cdng_connector_codeRepositories.DatasourceConnectorBitbucket(),
				"harness_platform_connector_customhealthsource":    connector.DatasourceConnectorCustomHealthSource(),
				"harness_platform_connector_datadog":               connector.DatasourceConnectorDatadog(),
				"harness_platform_connector_docker":                cdng_connector_artifactRepositories.DatasourceConnectorDocker(),
				"harness_platform_connector_jdbc":                  connector.DatasourceConnectorJDBC(),
				"harness_platform_connector_dynatrace":             connector.DatasourceConnectorDynatrace(),
				"harness_platform_connector_gcp":                   cdng_connector_cloudProviders.DatasourceConnectorGcp(),
				"harness_platform_connector_gcp_secret_manager":    pl_secretManagers.DatasourceConnectorGcpSM(),
				"harness_platform_connector_git":                   cdng_connector_codeRepositories.DatasourceConnectorGit(),
				"harness_platform_connector_github":                cdng_connector_codeRepositories.DatasourceConnectorGithub(),
				"harness_platform_connector_gitlab":                cdng_connector_codeRepositories.DatasourceConnectorGitlab(),
				"harness_platform_connector_helm":                  cdng_connector_artifactRepositories.DatasourceConnectorHelm(),
				"harness_platform_connector_oci_helm":              cdng_connector_artifactRepositories.DatasourceConnectorOciHelm(),
				"harness_platform_connector_jira":                  connector.DatasourceConnectorJira(),
				"harness_platform_connector_jenkins":               cdng_connector_artifactRepositories.DataSourceConnectorJenkins(),
				"harness_platform_connector_kubernetes":            cdng_connector_cloudProviders.DatasourceConnectorKubernetes(),
				"harness_platform_connector_nexus":                 cdng_connector_artifactRepositories.DatasourceConnectorNexus(),
				"harness_platform_connector_pagerduty":             connector.DatasourceConnectorPagerDuty(),
				"harness_platform_connector_prometheus":            connector.DatasourceConnectorPrometheus(),
				"harness_platform_connector_rancher":               cdng_connector_cloudProviders.DatasourceConnectorRancher(),
				"harness_platform_connector_splunk":                connector.DatasourceConnectorSplunk(),
				"harness_platform_connector_spot":                  cdng_connector_cloudProviders.DatasourceConnectorSpot(),
				"harness_platform_connector_terraform_cloud":       cdng_connector_cloudProviders.DatasourceConnectorTerraformCloud(),
				"harness_platform_connector_sumologic":             connector.DatasourceConnectorSumologic(),
				"harness_platform_connector_pdc":                   cdng_connector_cloudProviders.DatasourceConnectorPdc(),
				"harness_platform_connector_custom_secret_manager": pl_secretManagers.DatasourceConnectorCustomSM(),
				"harness_platform_current_account":                 pl_account.DataSourceCurrentAccount(),
				"harness_platform_current_user":                    pl_current_user.DataSourceCurrentUser(),
				"harness_platform_user":                            pl_user.DataSourceUser(),
				"harness_platform_environment":                     cdng_environment.DataSourceEnvironment(),
				"harness_platform_db_schema":                       dbschema.DataSourceDBSchema(),
				"harness_platform_db_instance":                     dbinstance.DataSourceDBInstance(),
				"harness_platform_environment_list":                cdng_environment.DataSourceEnvironmentList(),
				"harness_platform_environment_group":               cdng_environment_group.DataSourceEnvironmentGroup(),
				"harness_platform_environment_clusters_mapping":    pl_environment_clusters_mapping.DataSourceEnvironmentClustersMapping(),
				"harness_platform_environment_service_overrides":   cdng_environment_service_overrides.DataSourceEnvironmentServiceOverrides(),
				"harness_platform_service_overrides_v2":            cdng_service_overrides_v2.DataSourceServiceOverrides(),
				"harness_platform_provider":                        pl_provider.DataSourceProvider(),
				"harness_platform_overrides":                       cdng_overrides.DataSourceOverrides(),
				"harness_platform_gitops_agent":                    gitops_agent.DataSourceGitopsAgent(),
				"harness_platform_gitops_agent_deploy_yaml":        agent_yaml.DataSourceGitopsAgentDeployYaml(),
				"harness_platform_gitops_applications":             gitops_applications.DataSourceGitopsApplications(),
				"harness_platform_gitops_cluster":                  gitops_cluster.DataSourceGitopsCluster(),
				"harness_platform_gitops_gnupg":                    gitops_gnupg.DataSourceGitopsGnupg(),
				"harness_platform_gitops_app_project_mapping":      gitops_project_mapping.DatasourceGitopsAppProjectMapping(),
				"harness_platform_gitops_repository":               gitops_repository.DataSourceGitopsRepository(),
				"harness_platform_gitops_repo_cert":                gitops_repo_cert.DataSourceGitOpsRepoCert(),
				"harness_platform_gitops_repo_cred":                gitops_repo_cred.DataSourceGitOpsRepoCred(),
				"harness_platform_infrastructure":                  cdng_infrastructure.DataSourceInfrastructure(),
				"harness_platform_input_set":                       pipeline_input_set.DataSourceInputSet(),
				"harness_platform_monitored_service":               monitored_service.DataSourceMonitoredService(),
				"harness_platform_organization":                    organization.DataSourceOrganization(),
				"harness_platform_pipeline":                        pipeline.DataSourcePipeline(),
				"harness_platform_pipeline_list":                   pipeline.DataSourcePipelineList(),
				"harness_platform_permissions":                     pl_permissions.DataSourcePermissions(),
				"harness_platform_project":                         project.DataSourceProject(),
				"harness_platform_project_list":                    project.DataSourceProjectList(),
				"harness_platform_service":                         cdng_service.DataSourceService(),
				"harness_platform_service_list":                    cdng_service.DataSourceServiceList(),
				"harness_platform_usergroup":                       usergroup.DataSourceUserGroup(),
				"harness_platform_secret_text":                     secret.DataSourceSecretText(),
				"harness_platform_secret_file":                     secret.DataSourceSecretFile(),
				"harness_platform_secret_sshkey":                   secret.DataSourceSecretSSHKey(),
				"harness_platform_roles":                           roles.DataSourceRoles(),
				"harness_platform_resource_group":                  resource_group.DataSourceResourceGroup(),
				"harness_platform_service_account":                 service_account.DataSourceServiceAccount(),
				"harness_platform_triggers":                        pipeline_triggers.DataSourceTriggers(),
				"harness_platform_role_assignments":                role_assignments.DataSourceRoleAssignments(),
				"harness_platform_variables":                       cdng_variables.DataSourceVariables(),
				"harness_platform_connector_vault":                 pl_secretManagers.DataSourceConnectorVault(),
				"harness_platform_filters":                         cdng_filters.DataSourceFilters(),
				"harness_platform_pipeline_filters":                pipeline_filters.DataSourcePipelineFilters(),
				"harness_platform_ccm_filters":                     ccm_filters.DataSourceCCMFilters(),
				"harness_platform_template_filters":                pipeline_template_filters.DataSourceTemplateFilters(),
				"harness_application":                              application.DataSourceApplication(),
				"harness_current_account":                          account.DataSourceCurrentAccountConnector(),
				"harness_delegate":                                 delegate.DataSourceDelegate(),
				"harness_delegate_ids":                             delegate.DataSourceDelegateIds(),
				"harness_encrypted_text":                           secrets.DataSourceEncryptedText(),
				"harness_environment":                              environment.DataSourceEnvironment(),
				"harness_git_connector":                            cd_connector.DataSourceGitConnector(),
				"harness_secret_manager":                           secrets.DataSourceSecretManager(),
				"harness_service":                                  service.DataSourceService(),
				"harness_platform_slo":                             slo.DataSourceSloService(),
				"harness_platform_notification_rule":               notification_rule.DataSourceNotificationRuleService(),
				"harness_ssh_credential":                           secrets.DataSourceSshCredential(),
				"harness_sso_provider":                             sso.DataSourceSSOProvider(),
				"harness_user_group":                               user.DataSourceUserGroup(),
				"harness_user":                                     user.DataSourceUser(),
				"harness_yaml_config":                              yamlconfig.DataSourceYamlConfig(),
				"harness_platform_connector_azure_cloud_provider":  cdng_connector_cloudProviders.DataSourceConnectorAzureCloudProvider(),
				"harness_platform_connector_tas":                   cdng_connector_cloudProviders.DataSourceConnectorTas(),
				"harness_trigger":                                  cd_trigger.DataSourceTrigger(),
				"harness_platform_policy":                          policy.DataSourcePolicy(),
				"harness_platform_policyset":                       policyset.DataSourcePolicyset(),
				"harness_platform_manual_freeze":                   cdng_manual_freeze.DataSourceManualFreeze(),
				"harness_platform_connector_service_now":           connector.DataSourceConnectorSerivceNow(),
				"harness_platform_apikey":                          pl_apikey.DataSourceApiKey(),
				"harness_platform_token":                           pl_token.DataSourceToken(),
				"harness_autostopping_rule_vm":                     as_rule.DataSourceVMRule(),
				"harness_autostopping_rule_rds":                    as_rule.DataSourceRDSRule(),
				"harness_autostopping_rule_ecs":                    as_rule.DataSourceECSRule(),
				"harness_platform_file_store_file":                 cdng_file_store.DataSourceFileStoreNodeFile(),
				"harness_platform_file_store_folder":               cdng_file_store.DataSourceFileStoreNodeFolder(),
				"harness_autostopping_azure_proxy":                 load_balancer.DataSourceAzureProxy(),
				"harness_autostopping_aws_proxy":                   load_balancer.DataSourceAWSProxy(),
				"harness_autostopping_gcp_proxy":                   load_balancer.DataSourceGCPProxy(),
				"harness_autostopping_aws_alb":                     load_balancer.DataSourceAwsALB(),
				"harness_autostopping_azure_gateway":               load_balancer.DataSourceAzureGateway(),
				"harness_autostopping_schedule":                    schedule.DataSourceFixedSchedule(),
				"harness_platform_delegatetoken":                   pl_delegatetoken.DataSourceDelegateToken(),
				"harness_platform_workspace":                       workspace.DataSourceWorkspace(),
				"harness_platform_workspace_output":                workspace.DataSourceWorkspaceOutput(),
				"harness_platform_iacm_default_pipeline":           iacm.DataSourceIacmDefaultPipeline(),
				"harness_platform_repo":                            repo.DataSourceRepo(),
				"harness_platform_repo_rule_branch":                repo_rule_branch.DataSourceRepoBranchRule(),
				"harness_platform_repo_webhook":                    repo_webhook.DataSourceRepoWebhook(),
				"harness_platform_gitops_app_project":              gitops_project.DataSourceGitOpsProject(),
				"harness_platform_gitx_webhook":                    pipeline_gitx.DataSourceWebhook(),
				"harness_governance_rule_enforcement":              governance_enforcement.DatasourceRuleEnforcement(),
				"harness_governance_rule":                          governance_rule.DatasourceRule(),
				"harness_governance_rule_set":                      governance_rule_set.DatasourceRuleSet(),
				"harness_cluster_orchestrator":                     cluster_orchestrator.DataSourceClusterOrchestrator(),
				"harness_platform_infra_module":                    module_registry.DataSourceInfraModule(),
				"harness_chaos_infrastructure":                     infrastructure.DataSourceChaosInfrastructureService(),
				"harness_platform_har_registry": 					har_registry.DataSourceRegistry(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"harness_platform_template":                        pipeline_template.ResourceTemplate(),
				"harness_platform_connector_azure_key_vault":       pl_secretManagers.ResourceConnectorAzureKeyVault(),
				"harness_platform_connector_gcp_cloud_cost":        connector.ResourceConnectorGCPCloudCost(),
				"harness_platform_connector_kubernetes_cloud_cost": connector.ResourceConnectorKubernetesCloudCost(),
				"harness_platform_connector_azure_cloud_cost":      connector.ResourceConnectorAzureCloudCost(),
				"harness_platform_connector_appdynamics":           connector.ResourceConnectorAppDynamics(),
				"harness_platform_connector_elasticsearch":         connector.ResourceConnectorElasticSearch(),
				"harness_platform_connector_artifactory":           cdng_connector_artifactRepositories.ResourceConnectorArtifactory(),
				"harness_platform_connector_azure_artifacts":       cdng_connector_artifactRepositories.ResourceConnectorAzureArtifacts(),
				"harness_platform_connector_aws_secret_manager":    pl_secretManagers.ResourceConnectorAwsSM(),
				"harness_platform_connector_aws":                   cdng_connector_cloudProviders.ResourceConnectorAws(),
				"harness_platform_connector_awscc":                 connector.ResourceConnectorAwsCC(),
				"harness_platform_connector_awskms":                pl_secretManagers.ResourceConnectorAwsKms(),
				"harness_platform_connector_bitbucket":             cdng_connector_codeRepositories.ResourceConnectorBitbucket(),
				"harness_platform_connector_customhealthsource":    connector.ResourceConnectorCustomHealthSource(),
				"harness_platform_connector_datadog":               connector.ResourceConnectorDatadog(),
				"harness_platform_connector_docker":                cdng_connector_artifactRepositories.ResourceConnectorDocker(),
				"harness_platform_connector_jdbc":                  connector.ResourceConnectorJDBC(),
				"harness_platform_connector_dynatrace":             connector.ResourceConnectorDynatrace(),
				"harness_platform_connector_gcp":                   cdng_connector_cloudProviders.ResourceConnectorGcp(),
				"harness_platform_connector_gcp_secret_manager":    pl_secretManagers.ResourceConnectorGCPSecretManager(),
				"harness_platform_connector_git":                   cdng_connector_codeRepositories.ResourceConnectorGit(),
				"harness_platform_connector_github":                cdng_connector_codeRepositories.ResourceConnectorGithub(),
				"harness_platform_connector_gitlab":                cdng_connector_codeRepositories.ResourceConnectorGitlab(),
				"harness_platform_connector_helm":                  cdng_connector_artifactRepositories.ResourceConnectorHelm(),
				"harness_platform_connector_oci_helm":              cdng_connector_artifactRepositories.ResourceConnectorOciHelm(),
				"harness_platform_connector_jira":                  connector.ResourceConnectorJira(),
				"harness_platform_connector_jenkins":               cdng_connector_artifactRepositories.ResourceConnectorJenkins(),
				"harness_platform_connector_kubernetes":            cdng_connector_cloudProviders.ResourceConnectorK8s(),
				"harness_platform_connector_newrelic":              connector.ResourceConnectorNewRelic(),
				"harness_platform_connector_nexus":                 cdng_connector_artifactRepositories.ResourceConnectorNexus(),
				"harness_platform_connector_pagerduty":             connector.ResourceConnectorPagerDuty(),
				"harness_platform_connector_prometheus":            connector.ResourceConnectorPrometheus(),
				"harness_platform_connector_rancher":               cdng_connector_cloudProviders.ResourceConnectorK8sRancher(),
				"harness_platform_connector_splunk":                connector.ResourceConnectorSplunk(),
				"harness_platform_connector_spot":                  cdng_connector_cloudProviders.ResourceConnectorSpot(),
				"harness_platform_connector_terraform_cloud":       cdng_connector_cloudProviders.ResourceConnectorTerraformCloud(),
				"harness_platform_connector_sumologic":             connector.ResourceConnectorSumologic(),
				"harness_platform_connector_pdc":                   cdng_connector_cloudProviders.ResourceConnectorPdc(),
				"harness_platform_environment":                     cdng_environment.ResourceEnvironment(),
				"harness_platform_db_schema":                       dbschema.ResourceDBSchema(),
				"harness_platform_db_instance":                     dbinstance.ResourceDBInstance(),
				"harness_platform_environment_group":               cdng_environment_group.ResourceEnvironmentGroup(),
				"harness_platform_environment_clusters_mapping":    pl_environment_clusters_mapping.ResourceEnvironmentClustersMapping(),
				"harness_platform_environment_service_overrides":   cdng_environment_service_overrides.ResourceEnvironmentServiceOverrides(),
				"harness_platform_feature_flag":                    feature_flag.ResourceFeatureFlag(),
				"harness_platform_feature_flag_target_group":       feature_flag_target_group.ResourceFeatureFlagTargetGroup(),
				"harness_platform_feature_flag_target":             feature_flag_target.ResourceFeatureFlagTarget(),
				"harness_platform_service_overrides_v2":            cdng_service_overrides_v2.ResourceServiceOverrides(),
				"harness_platform_provider":                        pl_provider.ResourceProvider(),
				"harness_platform_overrides":                       cdng_overrides.ResourceOverrides(),
				"harness_platform_ff_api_key":                      ff_api_key.ResourceFFApiKey(),
				"harness_platform_gitops_agent":                    gitops_agent.ResourceGitopsAgent(),
				"harness_platform_gitops_applications":             gitops_applications.ResourceGitopsApplication(),
				"harness_platform_gitops_cluster":                  gitops_cluster.ResourceGitopsCluster(),
				"harness_platform_gitops_gnupg":                    gitops_gnupg.ResourceGitopsGnupg(),
				"harness_platform_gitops_app_project_mapping":      gitops_project_mapping.ResourceGitopsAppProjectMapping(),
				"harness_platform_gitops_repository":               gitops_repository.ResourceGitopsRepositories(),
				"harness_platform_gitops_app_project":              gitops_project.ResourceProject(),
				"harness_platform_gitops_repo_cert":                gitops_repo_cert.ResourceGitopsRepoCerts(),
				"harness_platform_gitops_repo_cred":                gitops_repo_cred.ResourceGitopsRepoCred(),
				"harness_platform_infrastructure":                  cdng_infrastructure.ResourceInfrastructure(),
				"harness_platform_input_set":                       pipeline_input_set.ResourceInputSet(),
				"harness_platform_monitored_service":               monitored_service.ResourceMonitoredService(),
				"harness_platform_organization":                    organization.ResourceOrganization(),
				"harness_platform_pipeline":                        pipeline.ResourcePipeline(),
				"harness_platform_project":                         project.ResourceProject(),
				"harness_platform_service":                         cdng_service.ResourceService(),
				"harness_platform_user":                            pl_user.ResourceUser(),
				"harness_platform_usergroup":                       usergroup.ResourceUserGroup(),
				"harness_platform_secret_text":                     secret.ResourceSecretText(),
				"harness_platform_secret_file":                     secret.ResourceSecretFile(),
				"harness_platform_secret_sshkey":                   secret.ResourceSecretSSHKey(),
				"harness_platform_roles":                           roles.ResourceRoles(),
				"harness_platform_resource_group":                  resource_group.ResourceResourceGroup(),
				"harness_platform_service_account":                 service_account.ResourceServiceAccount(),
				"harness_platform_triggers":                        pipeline_triggers.ResourceTriggers(),
				"harness_platform_role_assignments":                role_assignments.ResourceRoleAssignments(),
				"harness_platform_variables":                       cdng_variables.ResourceVariables(),
				"harness_platform_connector_vault":                 pl_secretManagers.ResourceConnectorVault(),
				"harness_platform_filters":                         cdng_filters.ResourceFilters(),
				"harness_platform_pipeline_filters":                pipeline_filters.ResourcePipelineFilters(),
				"harness_platform_ccm_filters":                     ccm_filters.ResourceCCMFilters(),
				"harness_platform_template_filters":                pipeline_template_filters.ResourceTemplateFilters(),
				"harness_add_user_to_group":                        user.ResourceAddUserToGroup(),
				"harness_application_gitsync":                      application.ResourceApplicationGitSync(),
				"harness_application":                              application.ResourceApplication(),
				"harness_delegate_approval":                        delegate.ResourceDelegateApproval(),
				"harness_cloudprovider_aws":                        cloudprovider.ResourceCloudProviderAws(),
				"harness_cloudprovider_azure":                      cloudprovider.ResourceCloudProviderAzure(),
				"harness_cloudprovider_datacenter":                 cloudprovider.ResourceCloudProviderDataCenter(),
				"harness_cloudprovider_gcp":                        cloudprovider.ResourceCloudProviderGcp(),
				"harness_cloudprovider_kubernetes":                 cloudprovider.ResourceCloudProviderK8s(),
				"harness_cloudprovider_spot":                       cloudprovider.ResourceCloudProviderSpot(),
				"harness_cloudprovider_tanzu":                      cloudprovider.ResourceCloudProviderTanzu(),
				"harness_encrypted_text":                           secrets.ResourceEncryptedText(),
				"harness_environment":                              environment.ResourceEnvironment(),
				"harness_git_connector":                            cd_connector.ResourceGitConnector(),
				"harness_infrastructure_definition":                environment.ResourceInfraDefinition(),
				"harness_service_ami":                              service.ResourceAMIService(),
				"harness_service_aws_codedeploy":                   service.ResourceAWSCodeDeployService(),
				"harness_service_aws_lambda":                       service.ResourceAWSLambdaService(),
				"harness_service_ecs":                              service.ResourceECSService(),
				"harness_service_helm":                             service.ResourceHelmService(),
				"harness_service_kubernetes":                       service.ResourceKubernetesService(),
				"harness_service_ssh":                              service.ResourceSSHService(),
				"harness_service_tanzu":                            service.ResourcePCFService(),
				"harness_service_winrm":                            service.ResourceWinRMService(),
				"harness_platform_slo":                             slo.ResourceSloService(),
				"harness_platform_notification_rule":               notification_rule.ResourceNotificationRuleService(),
				"harness_ssh_credential":                           secrets.ResourceSSHCredential(),
				"harness_user_group":                               user.ResourceUserGroup(),
				"harness_user_group_permissions":                   user.ResourceUserGroupPermissions(),
				"harness_user":                                     user.ResourceUser(),
				"harness_yaml_config":                              yamlconfig.ResourceYamlConfig(),
				"harness_platform_connector_azure_cloud_provider":  cdng_connector_cloudProviders.ResourceConnectorAzureCloudProvider(),
				"harness_platform_connector_tas":                   cdng_connector_cloudProviders.ResourceConnectorTas(),
				"harness_platform_policy":                          policy.ResourcePolicy(),
				"harness_platform_policyset":                       policyset.ResourcePolicyset(),
				"harness_platform_manual_freeze":                   cdng_manual_freeze.ResourceManualFreeze(),
				"harness_platform_connector_service_now":           connector.ResourceConnectorServiceNow(),
				"harness_platform_apikey":                          pl_apikey.ResourceApiKey(),
				"harness_platform_token":                           pl_token.ResourceToken(),
				"harness_autostopping_rule_vm":                     as_rule.ResourceVMRule(),
				"harness_autostopping_rule_rds":                    as_rule.ResourceRDSRule(),
				"harness_autostopping_rule_ecs":                    as_rule.ResourceECSRule(),
				"harness_platform_file_store_file":                 cdng_file_store.ResourceFileStoreNodeFile(),
				"harness_platform_file_store_folder":               cdng_file_store.ResourceFileStoreNodeFolder(),
				"harness_autostopping_azure_proxy":                 load_balancer.ResourceAzureProxy(),
				"harness_autostopping_aws_proxy":                   load_balancer.ResourceAWSProxy(),
				"harness_autostopping_gcp_proxy":                   load_balancer.ResourceGCPProxy(),
				"harness_autostopping_aws_alb":                     load_balancer.ResourceAwsALB(),
				"harness_autostopping_azure_gateway":               load_balancer.ResourceAzureGateway(),
				"harness_autostopping_schedule":                    schedule.ResourceVMRule(),
				"harness_platform_delegatetoken":                   pl_delegatetoken.ResourceDelegateToken(),
				"harness_platform_workspace":                       workspace.ResourceWorkspace(),
				"harness_platform_iacm_default_pipeline":           iacm.ResourceIacmDefaultPipeline(),
				"harness_platform_repo":                            repo.ResourceRepo(),
				"harness_platform_repo_rule_branch":                repo_rule_branch.ResourceRepoBranchRule(),
				"harness_platform_repo_webhook":                    repo_webhook.ResourceRepoWebhook(),
				"harness_platform_connector_custom_secret_manager": pl_secretManagers.ResourceConnectorCSM(),
				"harness_platform_gitx_webhook":                    pipeline_gitx.ResourceWebhook(),
				"harness_governance_rule_enforcement":              governance_enforcement.ResourceRuleEnforcement(),
				"harness_governance_rule":                          governance_rule.ResourceRule(),
				"harness_governance_rule_set":                      governance_rule_set.ResourceRuleSet(),
				"harness_cluster_orchestrator":                     cluster_orchestrator.ResourceClusterOrchestrator(),
				"harness_platform_infra_module":                    module_registry.ResourceInfraModule(),
				"harness_chaos_infrastructure":                     infrastructure.ResourceChaosInfrastructure(),
				"harness_platform_har_registry": 					har_registry.ResourceRegistry(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func getHttpClient(logger *logrus.Logger) *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Transport = logging.NewTransport(harness.SDKName, logger, cleanhttp.DefaultPooledClient().Transport)
	httpClient.RetryMax = 10
	return httpClient
}

func getOpenApiHttpClient(logger *logrus.Logger) *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Transport = openapi_client_logging.NewTransport(harness.SDKName, logger, cleanhttp.DefaultPooledClient().Transport)
	httpClient.RetryMax = 10
	return httpClient
}

func getCDClient(d *schema.ResourceData, version string) *cd.ApiClient {
	cfg := cd.DefaultConfig()
	cfg.AccountId = d.Get("account_id").(string)
	cfg.Endpoint = d.Get("endpoint").(string)
	cfg.APIKey = d.Get("api_key").(string)
	cfg.UserAgent = fmt.Sprintf("terraform-provider-harness-%s", version)
	cfg.HTTPClient = getHttpClient(cfg.Logger)
	cfg.DebugLogging = logging.IsDebugOrHigher(cfg.Logger)

	client, err := cd.NewClient(cfg)

	if err != nil {
		log.Printf("[WARN] error creating CD client: %s", err)
	}

	return client
}

func getPLClient(d *schema.ResourceData, version string) *nextgen.APIClient {
	cfg := nextgen.NewConfiguration()
	client := nextgen.NewAPIClient(&nextgen.Configuration{
		AccountId:    d.Get("account_id").(string),
		BasePath:     d.Get("endpoint").(string),
		ApiKey:       d.Get("platform_api_key").(string),
		UserAgent:    fmt.Sprintf("terraform-provider-harness-platform-%s", version),
		HTTPClient:   getHttpClient(cfg.Logger),
		DebugLogging: logging.IsDebugOrHigher(cfg.Logger),
	})

	return client
}

func getDBOpsClient(d *schema.ResourceData, version string) *dbops.APIClient {
	client := dbops.NewAPIClient(&dbops.Configuration{
		AccountId: d.Get("account_id").(string),
		BasePath:  d.Get("endpoint").(string),
		ApiKey:    d.Get("platform_api_key").(string),
		UserAgent: fmt.Sprintf("terraform-provider-harness-platform-%s", version),
	})

	return client
}

func getClient(d *schema.ResourceData, version string) *openapi_client_nextgen.APIClient {
	cfg := openapi_client_nextgen.NewConfiguration()
	client := openapi_client_nextgen.NewAPIClient(&openapi_client_nextgen.Configuration{
		AccountId:    d.Get("account_id").(string),
		BasePath:     d.Get("endpoint").(string),
		ApiKey:       d.Get("platform_api_key").(string),
		UserAgent:    fmt.Sprintf("terraform-provider-harness-platform-%s", version),
		HTTPClient:   getOpenApiHttpClient(cfg.Logger),
		DebugLogging: openapi_client_logging.IsDebugOrHigher(cfg.Logger),
	})

	return client
}

func getCodeClient(d *schema.ResourceData, version string) *code.APIClient {
	cfg := code.NewConfiguration()
	client := code.NewAPIClient(&code.Configuration{
		AccountId:     d.Get("account_id").(string),
		BasePath:      d.Get("endpoint").(string) + "/code/api/v1", // todo: this should be fixed in go sdk later
		ApiKey:        d.Get("platform_api_key").(string),
		UserAgent:     fmt.Sprintf("terraform-provider-harness-platform-%s", version),
		HTTPClient:    getOpenApiHttpClient(cfg.Logger),
		DefaultHeader: map[string]string{"X-Api-Key": d.Get("platform_api_key").(string)}, // todo: this should be fixed in go sdk later
		DebugLogging:  openapi_client_logging.IsDebugOrHigher(cfg.Logger),
	})
	return client
}

func getChaosClient(d *schema.ResourceData, version string) *chaos.APIClient {
	client := chaos.NewAPIClient(&chaos.Configuration{
		AccountId:     d.Get("account_id").(string),
		BasePath:      d.Get("endpoint").(string) + "/chaos/manager/api", // check if this can be taken from sdk
		ApiKey:        d.Get("platform_api_key").(string),
		UserAgent:     fmt.Sprintf("terraform-provider-harness-platform-%s", version),
		DefaultHeader: map[string]string{"X-Api-Key": d.Get("platform_api_key").(string)},
	})
	return client
}

func getHarClient(d *schema.ResourceData, version string) *har.APIClient {
	cfg := har.NewConfiguration()
	client := har.NewAPIClient(&har.Configuration{
		AccountId:     d.Get("account_id").(string),
		BasePath:      d.Get("endpoint").(string) + "/har/api/v1",
		ApiKey:        d.Get("platform_api_key").(string),
		UserAgent:     fmt.Sprintf("terraform-provider-harness-platform-%s", version),
		HTTPClient:    getOpenApiHttpClient(cfg.Logger),
		DefaultHeader: map[string]string{"X-Api-Key": d.Get("platform_api_key").(string)},
		DebugLogging:  openapi_client_logging.IsDebugOrHigher(cfg.Logger),
	})
	return client
}

// Setup the client for interacting with the Harness API
func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &internal.Session{
			AccountId:   d.Get("account_id").(string),
			Endpoint:    d.Get("endpoint").(string),
			CDClient:    getCDClient(d, version),
			PLClient:    getPLClient(d, version),
			Client:      getClient(d, version),
			CodeClient:  getCodeClient(d, version),
			DBOpsClient: getDBOpsClient(d, version),
			ChaosClient: getChaosClient(d, version),
			HARClient:   getHarClient(d, version),
		}, nil
	}
}
