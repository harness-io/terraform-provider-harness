package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/harness-openapi-go-client/nextgen"
	openapi_client_nextgen "github.com/harness/harness-openapi-go-client/nextgen"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccResourcePipeline(t *testing.T) {
	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	updatedName := fmt.Sprintf("%s_updated", id)

	resourceName := "harness_platform_pipeline.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccPipelineDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipeline(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccResourcePipeline(id, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       acctest.ProjectResourceImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{"git_details.0.commit_message", "git_details.0.connector_ref", "git_details.0.store_type"},
			},
		},
	})
}

func TestAccResourcePipelineInline(t *testing.T) {
	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	updatedName := fmt.Sprintf("%s_updated", id)

	resourceName := "harness_platform_pipeline.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccPipelineDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipelineInline(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccResourcePipelineInline(id, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: acctest.ProjectResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccResourcePipelineImportFromGit(t *testing.T) {
	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id

	resourceName := "harness_platform_pipeline.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccPipelineDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipelineImportFromGit(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "gitx"),
					resource.TestCheckResourceAttr(resourceName, "name", "gitx"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       acctest.ProjectResourceImportStateIdGitFunc(resourceName),
				ImportStateVerifyIgnore: []string{"git_import_info.0.branch_name", "git_import_info.0.connector_ref", "git_import_info.0.file_path", "git_import_info.0.repo_name", "import_from_git", "pipeline_import_request.0.pipeline_description", "pipeline_import_request.0.pipeline_name", "git_import_info.#", "git_import_info.0.%", "pipeline_import_request.#", "pipeline_import_request.0.%"},
			},
		},
	})
}

func TestAccResourcePipeline_DeleteUnderlyingResource(t *testing.T) {
	name := t.Name()
	id := fmt.Sprintf("%s_%s", name, utils.RandStringBytes(5))
	project_id := id
	org_id := id
	resourceName := "harness_platform_pipeline.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipelineInline(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					c, ctx := acctest.TestAccGetClientWithContext()
					_, err := c.PipelinesApi.DeletePipeline(ctx, org_id, project_id, id, &nextgen.PipelinesApiDeletePipelineOpts{})
					require.NoError(t, err)
				},
				Config:             testAccResourcePipelineInline(id, name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccGetPipeline(resourceName string, state *terraform.State) (*openapi_client_nextgen.PipelineGetResponseBody, error) {
	r := acctest.TestAccGetResource(resourceName, state)
	c, ctx := acctest.TestAccGetClientWithContext()
	id := r.Primary.ID
	orgId := r.Primary.Attributes["org_id"]
	projId := r.Primary.Attributes["project_id"]
	branch_name := r.Primary.Attributes["git_details.0.branch_name"]

	resp, _, err := c.PipelinesApi.GetPipeline(ctx, orgId, projId, id, &openapi_client_nextgen.PipelinesApiGetPipelineOpts{HarnessAccount: optional.NewString(c.AccountId), BranchName: optional.NewString(branch_name)})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func testAccPipelineDestroy(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		pipeline, _ := testAccGetPipeline(resourceName, state)
		if pipeline != nil {
			return fmt.Errorf("Found pipeline: %s", pipeline.PipelineYaml)
		}

		return nil
	}
}

func testAccResourcePipeline(id string, name string) string {
	return fmt.Sprintf(`
				resource "harness_platform_organization" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
				}
				resource "harness_platform_project" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
					org_id = harness_platform_organization.test.id
					color = "#472848"
				}
        resource "harness_platform_pipeline" "test" {
                        identifier = "%[1]s"
                        org_id = harness_platform_project.test.org_id
						project_id = harness_platform_project.test.id
                        name = "%[2]s"
                        tags = ["foo:bar", "bar:foo"]
                        git_details {
                            branch_name = "main"
                            commit_message = "Commit"
                            file_path = ".harness/GitEnabledPipeline%[1]s.yaml"
                            connector_ref = "account.TF_Jajoo_github_connector"
                            store_type = "REMOTE"
                            repo_name = "jajoo_git"
                        }
            yaml = <<-EOT
                pipeline:
                    name: %[2]s
                    identifier: %[1]s   
                    allowStageExecutions: false
                    projectIdentifier: ${harness_platform_project.test.id}
                    orgIdentifier: ${harness_platform_project.test.org_id}
                    tags:
                      foo: bar
                      bar: foo
                    stages:
                        - stage:
                            name: dep
                            identifier: dep
                            description: ""
                            type: Deployment
                            spec:
                                serviceConfig:
                                    serviceRef: service
                                    serviceDefinition:
                                        type: Kubernetes
                                        spec:
                                            variables: []
                                infrastructure:
                                    environmentRef: testenv
                                    infrastructureDefinition:
                                        type: KubernetesDirect
                                        spec:
                                            connectorRef: testconf
                                            namespace: test
                                            releaseName: release-<+INFRA_KEY>
                                    allowSimultaneousDeployments: false
                                execution:
                                    steps:
                                        - stepGroup:
                                                name: Canary Deployment
                                                identifier: canaryDepoyment
                                                steps:
                                                    - step:
                                                        name: Canary Deployment
                                                        identifier: canaryDeployment
                                                        type: K8sCanaryDeploy
                                                        timeout: 10m
                                                        spec:
                                                            instanceSelection:
                                                                type: Count
                                                                spec:
                                                                    count: 1
                                                            skipDryRun: false
                                                    - step:
                                                        name: Canary Delete
                                                        identifier: canaryDelete
                                                        type: K8sCanaryDelete
                                                        timeout: 10m
                                                        spec: {}
                                                rollbackSteps:
                                                    - step:
                                                        name: Canary Delete
                                                        identifier: rollbackCanaryDelete
                                                        type: K8sCanaryDelete
                                                        timeout: 10m
                                                        spec: {}
                                        - stepGroup:
                                                name: Primary Deployment
                                                identifier: primaryDepoyment
                                                steps:
                                                    - step:
                                                        name: Rolling Deployment
                                                        identifier: rollingDeployment
                                                        type: K8sRollingDeploy
                                                        timeout: 10m
                                                        spec:
                                                            skipDryRun: false
                                                rollbackSteps:
                                                    - step:
                                                        name: Rolling Rollback
                                                        identifier: rollingRollback
                                                        type: K8sRollingRollback
                                                        timeout: 10m
                                                        spec: {}
                                    rollbackSteps: []
                            tags: {}
                            failureStrategies:
                                - onFailure:
                                        errors:
                                            - AllErrors
                                        action:
                                            type: StageRollback
            EOT
        }
        `, id, name)
}

func testAccResourcePipelineInline(id string, name string) string {
	return fmt.Sprintf(`
				resource "harness_platform_organization" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
				}
				resource "harness_platform_project" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
					org_id = harness_platform_organization.test.id
					color = "#472848"
				}
        resource "harness_platform_pipeline" "test" {
                        identifier = "%[1]s"
                        org_id = harness_platform_project.test.org_id
						project_id = harness_platform_project.test.id
                        name = "%[2]s"
                        tags = ["foo:bar", "bar:foo"]
            yaml = <<-EOT
                pipeline:
                    name: %[2]s
                    identifier: %[1]s
                    allowStageExecutions: false
                    projectIdentifier: ${harness_platform_project.test.id}
                    orgIdentifier: ${harness_platform_project.test.org_id}
                    tags:
                      foo: bar
                      bar: foo
                    stages:
                        - stage:
                            name: dep
                            identifier: dep
                            description: ""
                            type: Deployment
                            spec:
                                serviceConfig:
                                    serviceRef: service
                                    serviceDefinition:
                                        type: Kubernetes
                                        spec:
                                            variables: []
                                infrastructure:
                                    environmentRef: testenv
                                    infrastructureDefinition:
                                        type: KubernetesDirect
                                        spec:
                                            connectorRef: testconf
                                            namespace: test
                                            releaseName: release-<+INFRA_KEY>
                                    allowSimultaneousDeployments: false
                                execution:
                                    steps:
                                        - stepGroup:
                                                name: Canary Deployment
                                                identifier: canaryDepoyment
                                                steps:
                                                    - step:
                                                        name: Canary Deployment
                                                        identifier: canaryDeployment
                                                        type: K8sCanaryDeploy
                                                        timeout: 10m
                                                        spec:
                                                            instanceSelection:
                                                                type: Count
                                                                spec:
                                                                    count: 1
                                                            skipDryRun: false
                                                    - step:
                                                        name: Canary Delete
                                                        identifier: canaryDelete
                                                        type: K8sCanaryDelete
                                                        timeout: 10m
                                                        spec: {}
                                                rollbackSteps:
                                                    - step:
                                                        name: Canary Delete
                                                        identifier: rollbackCanaryDelete
                                                        type: K8sCanaryDelete
                                                        timeout: 10m
                                                        spec: {}
                                        - stepGroup:
                                                name: Primary Deployment
                                                identifier: primaryDepoyment
                                                steps:
                                                    - step:
                                                        name: Rolling Deployment
                                                        identifier: rollingDeployment
                                                        type: K8sRollingDeploy
                                                        timeout: 10m
                                                        spec:
                                                            skipDryRun: false
                                                rollbackSteps:
                                                    - step:
                                                        name: Rolling Rollback
                                                        identifier: rollingRollback
                                                        type: K8sRollingRollback
                                                        timeout: 10m
                                                        spec: {}
                                    rollbackSteps: []
                            tags: {}
                            failureStrategies:
                                - onFailure:
                                        errors:
                                            - AllErrors
                                        action:
                                            type: StageRollback
            EOT
        }
        `, id, name)
}

func testAccResourcePipelineImportFromGit(id string, name string) string {
	return fmt.Sprintf(`
        resource "harness_platform_organization" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
				}
		resource "harness_platform_project" "Project_Test" {
				identifier = "TF_GitX_Pipeline_Test"
				name = "TF_GitX_Pipeline_Test"
				color = "#0063F7"
				org_id = "default"
		}
        resource "harness_platform_pipeline" "test" {
                        identifier = "gitx"
                        org_id = "default"
						project_id = "TF_GitX_Pipeline_Test"
                        name = "gitx"
                        import_from_git = true
                        git_import_info {
                            branch_name = "main"
                            file_path = ".harness/gitx.yaml"
                            connector_ref = "account.TF_open_repo_github_connector"
                            repo_name = "open-repo"
                        }
                        pipeline_import_request {
                            pipeline_name = "gitx"
                            pipeline_description = "Pipeline Description"
                        }
                }
        `, id, name)
}

func testAccResourcePipelineHC(id string, name string) string {
	return fmt.Sprintf(`
				resource "harness_platform_organization" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
				}
				resource "harness_platform_project" "test" {
					identifier = "%[1]s"
					name = "%[2]s"
					org_id = harness_platform_organization.test.id
					color = "#472848"
				}
				resource "harness_platform_repo" "test" {
				  identifier     = "%[1]s"
				  org_id         = harness_platform_organization.test.id
				  project_id     = harness_platform_project.test.id
				  default_branch = "main"
				  source {
					repo = "octocat/hello-worId"
					type = "github"
				  }
				}
        resource "harness_platform_pipeline" "test" {
                        identifier = "%[1]s"
                        org_id = harness_platform_project.test.org_id
						project_id = harness_platform_project.test.id
                        name = "%[2]s"
                        tags = ["foo:bar", "bar:foo"]
                        git_details {
                            branch_name = "main"
                            commit_message = "Commit"
                            file_path = ".harness/GitEnabledPipeline%[1]s.yaml"
                            store_type = "REMOTE"
                            repo_name = harness_platform_repo.test.name
							is_harness_code_repo = true
                        }
            yaml = <<-EOT
                pipeline:
                    name: %[2]s
                    identifier: %[1]s   
                    allowStageExecutions: false
                    projectIdentifier: ${harness_platform_project.test.id}
                    orgIdentifier: ${harness_platform_project.test.org_id}
                    tags:
                      foo: bar
                      bar: foo
                    stages:
                        - stage:
                            name: dep
                            identifier: dep
                            description: ""
                            type: Deployment
                            spec:
                                serviceConfig:
                                    serviceRef: service
                                    serviceDefinition:
                                        type: Kubernetes
                                        spec:
                                            variables: []
                                infrastructure:
                                    environmentRef: testenv
                                    infrastructureDefinition:
                                        type: KubernetesDirect
                                        spec:
                                            connectorRef: testconf
                                            namespace: test
                                            releaseName: release-<+INFRA_KEY>
                                    allowSimultaneousDeployments: false
                                execution:
                                    steps:
                                        - stepGroup:
                                                name: Canary Deployment
                                                identifier: canaryDepoyment
                                                steps:
                                                    - step:
                                                        name: Canary Deployment
                                                        identifier: canaryDeployment
                                                        type: K8sCanaryDeploy
                                                        timeout: 10m
                                                        spec:
                                                            instanceSelection:
                                                                type: Count
                                                                spec:
                                                                    count: 1
                                                            skipDryRun: false
                                                    - step:
                                                        name: Canary Delete
                                                        identifier: canaryDelete
                                                        type: K8sCanaryDelete
                                                        timeout: 10m
                                                        spec: {}
                                                rollbackSteps:
                                                    - step:
                                                        name: Canary Delete
                                                        identifier: rollbackCanaryDelete
                                                        type: K8sCanaryDelete
                                                        timeout: 10m
                                                        spec: {}
                                        - stepGroup:
                                                name: Primary Deployment
                                                identifier: primaryDepoyment
                                                steps:
                                                    - step:
                                                        name: Rolling Deployment
                                                        identifier: rollingDeployment
                                                        type: K8sRollingDeploy
                                                        timeout: 10m
                                                        spec:
                                                            skipDryRun: false
                                                rollbackSteps:
                                                    - step:
                                                        name: Rolling Rollback
                                                        identifier: rollingRollback
                                                        type: K8sRollingRollback
                                                        timeout: 10m
                                                        spec: {}
                                    rollbackSteps: []
                            tags: {}
                            failureStrategies:
                                - onFailure:
                                        errors:
                                            - AllErrors
                                        action:
                                            type: StageRollback
            EOT
        }
        `, id, name)
}

func TestAccResourcePipelineHarnessCode(t *testing.T) {
	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	updatedName := fmt.Sprintf("%s_updated", id)

	resourceName := "harness_platform_pipeline.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccPipelineDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipelineHC(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccResourcePipelineHC(id, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       acctest.ProjectResourceImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{"git_details.0.commit_message", "git_details.0.connector_ref", "git_details.0.store_type"},
			},
		},
	})
}
