package acctest

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/harness/harness-go-sdk/harness/cd/graphql"
	"github.com/harness/harness-go-sdk/harness/code"
	"github.com/harness/harness-go-sdk/harness/dbops"
	"github.com/harness/harness-go-sdk/harness/har"
	"github.com/harness/harness-go-sdk/harness/helpers"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/harness-go-sdk/harness/policymgmt"
	"github.com/harness/harness-go-sdk/harness/utils"
	openapi_client_nextgen "github.com/harness/harness-openapi-go-client/nextgen"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/harness/terraform-provider-harness/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	TestAccSecretFileId      = "2WnPVgLGSZW6KbApZuxeaw"
	TestAccDefaultUsageScope = `
	usage_scope {
		environment_filter_type = "NON_PRODUCTION_ENVIRONMENTS"
	}
`
)

func TestAccConfigureProvider() {
	TestAccProviderConfigure.Do(func() {
		TestAccProvider = provider.Provider("dev")()

		config := map[string]interface{}{
			"endpoint":         helpers.EnvVars.Endpoint.GetWithDefault(utils.BaseUrl),
			"account_id":       helpers.EnvVars.AccountId.Get(),
			"api_key":          helpers.EnvVars.ApiKey.Get(),
			"platform_api_key": helpers.EnvVars.PlatformApiKey.Get(),
		}

		TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(config))
	})
}

func TestAccPreCheck(t *testing.T) {
	TestAccConfigureProvider()
}

var TestAccProvider *schema.Provider
var TestAccProviderConfigure sync.Once

func TestAccGetResource(resourceName string, state *terraform.State) *terraform.ResourceState {
	rm := state.RootModule()
	return rm.Resources[resourceName]
}

func TestAccGetApiClientFromProvider() *internal.Session {
	return TestAccProvider.Meta().(*internal.Session)
}

func TestAccGetPlatformClientWithContext() (*nextgen.APIClient, context.Context) {
	return TestAccProvider.Meta().(*internal.Session).GetPlatformClientWithContext(context.Background())
}

func TestAccGetDBOpsClientWithContext() (*dbops.APIClient, context.Context) {
	return TestAccProvider.Meta().(*internal.Session).GetDBOpsClientWithContext(context.Background())
}

func TestAccGetClientWithContext() (*openapi_client_nextgen.APIClient, context.Context) {
	return TestAccProvider.Meta().(*internal.Session).GetClientWithContext(context.Background())
}

func TestAccGetPolicyManagementClient() *policymgmt.APIClient {
	return TestAccProvider.Meta().(*internal.Session).GetPolicyManagementClient()
}

func TestAccGetCodeClientWithContext() (*code.APIClient, context.Context) {
	return TestAccProvider.Meta().(*internal.Session).GetCodeClientWithContext(context.Background())
}

func TestAccGetHarClientWithContext() (*har.APIClient, context.Context) {
	return TestAccProvider.Meta().(*internal.Session).GetHarClientWithContext(context.Background())
}

func TestAccGetApplication(resourceName string, state *terraform.State) (*graphql.Application, error) {
	r := TestAccGetResource(resourceName, state)
	c := TestAccGetApiClientFromProvider()
	id := r.Primary.ID

	return c.CDClient.ApplicationClient.GetApplicationById(id)
}

func PipelineResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		var pipelineId string
		if len(primary.Attributes["pipeline_id"]) != 0 {
			pipelineId = primary.Attributes["pipeline_id"]
		}
		if len(primary.Attributes["target_id"]) != 0 {
			pipelineId = primary.Attributes["target_id"]
		}
		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, pipelineId, id), nil
	}
}

func EnvRelatedResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		var envId string
		if len(primary.Attributes["env_id"]) != 0 {
			envId = primary.Attributes["env_id"]
		}
		if orgId == "" {
			return fmt.Sprintf("%s/%s", envId, id), nil
		}
		if projId == "" {
			return fmt.Sprintf("%s/%s/%s", orgId, envId, id), nil
		}

		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, envId, id), nil
	}
}

func DBInstanceResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		schema := primary.Attributes["schema"]

		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, schema, id), nil
	}
}

func OverridesV1ResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		var envId string
		if len(primary.Attributes["env_id"]) != 0 {
			envId = primary.Attributes["env_id"]
		}
		if orgId == "" {
			return fmt.Sprintf("%s", envId), nil
		}
		if projId == "" {
			return fmt.Sprintf("%s/%s", orgId, envId), nil
		}

		return fmt.Sprintf("%s/%s/%s", orgId, projId, envId), nil
	}
}

func RepoResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		return fmt.Sprintf("%s/%s/%s", orgId, projId, id), nil
	}
}

func ProjectResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		return fmt.Sprintf("%s/%s/%s", orgId, projId, id), nil
	}
}

func ProjectResourceImportStateIdGitFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		branch_name := primary.Attributes["git_details.0.branch_name"]
		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, id, branch_name), nil
	}
}

func UserResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		email := primary.Attributes["email"]
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		return fmt.Sprintf("%s/%s/%s", email, orgId, projId), nil
	}
}

func UserResourceImportStateIdFuncAccountLevel(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		email := primary.Attributes["email"]
		return fmt.Sprintf("%s", email), nil
	}
}

func UserResourceImportStateIdFuncOrgLevel(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		email := primary.Attributes["email"]
		orgId := primary.Attributes["org_id"]
		return fmt.Sprintf("%s/%s", email, orgId), nil
	}
}

func AccountLevelResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		return fmt.Sprintf("%s", id), nil
	}
}

func AccountFilterImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		type_ := primary.Attributes["type"]
		return fmt.Sprintf("%s/%s", id, type_), nil
	}
}
func ProjectFilterImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		type_ := primary.Attributes["type"]
		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, id, type_), nil
	}
}

func GitopsAgentProjectLevelResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID

		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		agentId := primary.Attributes["agent_id"]
		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, agentId, id), nil
	}
}

// Import of GitopsAppProjectMapping resource is always on project level
// terraform import  harness_platform_gitops_app_project_mapping.example org_id/projec_id/scope_prefixed_agent_id/argo_proj_name
func GitopsAppProjectMappingResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary

		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		agentId := primary.Attributes["agent_id"]
		argoProjName := primary.Attributes["argo_project_name"]
		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, agentId, argoProjName), nil
	}
}

func GitopsAgentOrgLevelResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		agentId := primary.Attributes["agent_id"]
		return fmt.Sprintf("%s/%s/%s", orgId, agentId, id), nil
	}
}

func GitopsProjectImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		agentId := primary.Attributes["agent_id"]
		query_name := primary.Attributes["project.0.metadata.0.name"]
		return fmt.Sprintf("%s/%s", agentId, query_name), nil
	}
}

func GitopsAgentAccountLevelResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		agentId := primary.Attributes["agent_id"]
		return fmt.Sprintf("%s/%s", agentId, id), nil
	}
}

func RepoRuleProjectLevelResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		repoIdentifier := primary.Attributes["repo_identifier"]
		return fmt.Sprintf("%s/%s/%s/%s", orgId, projId, repoIdentifier, id), nil
	}
}

func OrgResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		return fmt.Sprintf("%s/%s", orgId, id), nil
	}
}

func OrgFilterImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		id := primary.ID
		orgId := primary.Attributes["org_id"]
		type_ := primary.Attributes["type"]
		return fmt.Sprintf("%s/%s/%s", orgId, id, type_), nil
	}
}

func GitopsWebhookImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		primary := s.RootModule().Resources[resourceName].Primary
		webhook_identifier := primary.Attributes["identifier"]
		orgId := primary.Attributes["org_id"]
		projId := primary.Attributes["project_id"]
		return fmt.Sprintf("%s/%s/%s", webhook_identifier, orgId, projId), nil
	}
}

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"harness": func() (*schema.Provider, error) {
		return provider.Provider("dev")(), nil
	},
}

func TestAccResourceAwsCloudProvider(name string) string {
	return fmt.Sprintf(`
	data "harness_secret_manager" "default" {
		default = true
	}

	resource "harness_encrypted_text" "aws_access_key" {
		name = "%[1]s_access_key"
		value = "%[2]s"
		secret_manager_id = data.harness_secret_manager.default.id
	}

	resource "harness_encrypted_text" "aws_secret_key" {
		name = "%[1]s_secret_key"
		value = "%[3]s"
		secret_manager_id = data.harness_secret_manager.default.id
	}
	
	resource "harness_cloudprovider_aws" "test" {
		name = "%[1]s"
		access_key_id_secret_name = harness_encrypted_text.aws_access_key.name
		secret_access_key_secret_name = harness_encrypted_text.aws_secret_key.name
	}	
`, name, helpers.TestEnvVars.AwsAccessKeyId.Get(), helpers.TestEnvVars.AwsSecretAccessKey.Get())
}

func TestAccResourceInfraDefEnvironment(name string) string {
	return fmt.Sprintf(`
		resource "harness_application" "test" {
			name = "%[1]s"
		}

		resource "harness_environment" "test" {
			name = "%[1]s"
			app_id = harness_application.test.id
			type = "NON_PROD"
		}
`, name)
}

func TestAccResourceGitConnector_default(name string) string {

	return fmt.Sprintf(`
		data "harness_secret_manager" "test" {
			default = true
		}

		resource "harness_encrypted_text" "test" {
			name 							= "%[1]s"
			value 					  = "foo"
			secret_manager_id = data.harness_secret_manager.test.id
		}

		resource "harness_git_connector" "test" {
			name = "%[1]s"
			url = "https://github.com/micahlmartin/harness-demo"
			branch = "master"
			generate_webhook_url = true
			password_secret_id = harness_encrypted_text.test.id
			url_type = "REPO"
			username = "someuser"
		}	
`, name)
}
