package pipeline

import (
	"context"
	"net/http"

	"github.com/antihax/optional"
	"github.com/harness/harness-openapi-go-client/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePipeline() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness pipeline.",

		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineCreateOrUpdate,
		DeleteContext: resourcePipelineDelete,
		CreateContext: resourcePipelineCreateOrUpdate,
		Importer:      helpers.ProjectResourceImporter,

		Schema: map[string]*schema.Schema{
			"yaml": {
				Description:      "YAML of the pipeline." + helpers.Descriptions.YamlText.String(),
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: helpers.YamlDiffSuppressFunction,
			},
			"git_details": {
				Description: "Contains parameters related to creating an Entity for Git Experience.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch_name": {
							Description: "Name of the branch.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"file_path": {
							Description: "File path of the Entity in the repository.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"commit_message": {
							Description: "Commit message used for the merge commit.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"base_branch": {
							Description: "Name of the default branch (this checks out a new branch titled by branch_name).",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"connector_ref": {
							Description: "Identifier of the Harness Connector used for CRUD operations on the Entity." + helpers.Descriptions.ConnectorRefText.String(),
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"store_type": {
							Description:  "Specifies whether the Entity is to be stored in Git or not. Possible values: INLINE, REMOTE.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"INLINE", "REMOTE"}, false),
							Computed:     true,
						},
						"repo_name": {
							Description: "Name of the repository.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"last_object_id": {
							Description: "Last object identifier (for Github). To be provided only when updating Pipeline.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"last_commit_id": {
							Description: "Last commit identifier (for Git Repositories other than Github). To be provided only when updating Pipeline.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"is_harness_code_repo": {
							Description: "If the repo is harness code",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"template_applied": {
				Description: "If true, returns Pipeline YAML with Templates applied on it.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"template_applied_pipeline_yaml": {
				Description: "Pipeline YAML after resolving Templates (returned as a String).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"import_from_git": {
				Description: "Flag to set if importing from Git",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"git_import_info": {
				Description: "Contains Git Information for importing entities from Git",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch_name": {
							Description: "Name of the branch.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"file_path": {
							Description: "File path of the Entity in the repository.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"connector_ref": {
							Description: "Identifier of the Harness Connector used for importing entity from Git" + helpers.Descriptions.ConnectorRefText.String(),
							Type:        schema.TypeString,
							Optional:    true,
						},
						"repo_name": {
							Description: "Name of the repository.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"pipeline_import_request": {
				Description: "Contains parameters for importing a pipeline",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_name": {
							Description: "Name of the pipeline.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"pipeline_description": {
							Description: "Description of the pipeline.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}

	helpers.SetProjectLevelResourceSchema(resource.Schema)
	resource.Schema["tags"].Description = resource.Schema["tags"].Description + " These should match the tag value passed in the YAML; if this parameter is null or not passed, the tags specified in YAML should also be null."
	return resource
}

func resourcePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetClientWithContext(ctx)

	id := d.Id()
	org_id := d.Get("org_id").(string)
	project_id := d.Get("project_id").(string)
	template_applied := d.Get("template_applied").(bool)
	var branch_name optional.String
	branch_name = helpers.BuildField(d, "git_details.0.branch_name")
	var store_type = helpers.BuildField(d, "git_details.0.store_type")
	var base_branch = helpers.BuildField(d, "git_details.0.base_branch")
	var commit_message = helpers.BuildField(d, "git_details.0.commit_message")
	var connector_ref = helpers.BuildField(d, "git_details.0.connector_ref")
	var object_id = helpers.BuildField(d, "git_details.0.object_id")
	resp, httpResp, err := c.PipelinesApi.GetPipeline(ctx,
		org_id,
		project_id,
		id,
		&nextgen.PipelinesApiGetPipelineOpts{HarnessAccount: optional.NewString(c.AccountId), BranchName: branch_name},
	)

	if httpResp.StatusCode == 404 {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}
	print(object_id.Value())

	readPipeline(d, resp, org_id, project_id, template_applied, store_type, base_branch, commit_message, connector_ref)

	return nil
}

func resourcePipelineCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetClientWithContext(ctx)

	var err error
	var pipeline_id string
	var branch_name string
	var store_type optional.String
	var base_branch optional.String
	var commit_message optional.String
	var connector_ref optional.String
	var httpResp *http.Response
	id := d.Id()
	org_id := d.Get("org_id").(string)
	project_id := d.Get("project_id").(string)
	template_applied := d.Get("template_applied").(bool)

	if id == "" {
		if d.Get("import_from_git").(bool) {
			pipeline_id = d.Get("identifier").(string)
			pipeline_import_request_body := createImportFromGitRequest(d)
			branch_name = pipeline_import_request_body.GitImportInfo.BranchName

			_, httpResp, err = c.PipelinesApi.ImportPipelineFromGit(ctx, org_id, project_id, pipeline_id,
				&nextgen.PipelinesApiImportPipelineFromGitOpts{
					Body:           optional.NewInterface(pipeline_import_request_body),
					HarnessAccount: optional.NewString(c.AccountId)})
		} else {
			pipeline := buildCreatePipeline(d)
			if pipeline.GitDetails != nil {
				base_branch = optional.NewString(pipeline.GitDetails.BaseBranch)
				store_type = optional.NewString(pipeline.GitDetails.StoreType)
				commit_message = optional.NewString(pipeline.GitDetails.CommitMessage)
				connector_ref = optional.NewString(pipeline.GitDetails.ConnectorRef)
				branch_name = pipeline.GitDetails.BranchName
			}

			pipeline_id = pipeline.Identifier
			_, httpResp, err = c.PipelinesApi.CreatePipeline(ctx, pipeline, org_id, project_id,
				&nextgen.PipelinesApiCreatePipelineOpts{HarnessAccount: optional.NewString(c.AccountId)})
		}
	} else {
		pipeline := buildUpdatePipeline(d)
		store_type = helpers.BuildField(d, "git_details.0.store_type")
		connector_ref = helpers.BuildField(d, "git_details.0.connector_ref")
		pipeline_id = pipeline.Identifier

		if pipeline.GitDetails != nil {
			base_branch = optional.NewString(pipeline.GitDetails.BaseBranch)
			branch_name = pipeline.GitDetails.BranchName
			commit_message = optional.NewString(pipeline.GitDetails.CommitMessage)
			_, httpResp, err = c.PipelinesApi.UpdatePipelineGitMetadata(ctx, org_id, project_id, pipeline_id, &nextgen.PipelinesApiUpdatePipelineGitMetadataOpts{
				Body:           optional.NewInterface(pipeline),
				HarnessAccount: optional.NewString(c.AccountId),
			})
			if err != nil {
				return helpers.HandleApiError(err, d, httpResp)
			}
		}
		_, httpResp, err = c.PipelinesApi.UpdatePipeline(ctx, pipeline, org_id, project_id, id,
			&nextgen.PipelinesApiUpdatePipelineOpts{HarnessAccount: optional.NewString(c.AccountId)})
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	// The create/update methods don't return the yaml in the response, so we need to query for it again.
	resp, httpResp, err := c.PipelinesApi.GetPipeline(ctx, org_id, project_id, pipeline_id,
		&nextgen.PipelinesApiGetPipelineOpts{HarnessAccount: optional.NewString(c.AccountId), BranchName: optional.NewString(branch_name)})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readPipeline(d, resp, org_id, project_id, template_applied, store_type, base_branch, commit_message, connector_ref)

	return nil
}

func createImportFromGitRequest(d *schema.ResourceData) *nextgen.PipelineImportRequestBody {

	pipeline_git_import_info := &nextgen.GitImportInfo{}
	if attr, ok := d.GetOk("git_import_info"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		if attr, ok := config["branch_name"]; ok {
			pipeline_git_import_info.BranchName = attr.(string)
		}
		if attr, ok := config["file_path"]; ok {
			pipeline_git_import_info.FilePath = attr.(string)
		}
		if attr, ok := config["connector_ref"]; ok {
			pipeline_git_import_info.ConnectorRef = attr.(string)
		}
		if attr, ok := config["repo_name"]; ok {
			pipeline_git_import_info.RepoName = attr.(string)
		}
	}

	pipeline_import_request := &nextgen.PipelineImportRequestDto{}
	if attr, ok := d.GetOk("pipeline_import_request"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		if attr, ok := config["pipeline_name"]; ok {
			pipeline_import_request.PipelineName = attr.(string)
		}
		if attr, ok := config["pipeline_description"]; ok {
			pipeline_import_request.PipelineDescription = attr.(string)
		}
	}

	pipeline_import_request_body := &nextgen.PipelineImportRequestBody{}
	pipeline_import_request_body.GitImportInfo = pipeline_git_import_info
	pipeline_import_request_body.PipelineImportRequest = pipeline_import_request

	return pipeline_import_request_body
}

func resourcePipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetClientWithContext(ctx)

	id := d.Get("identifier").(string)
	org_id := d.Get("org_id").(string)
	project_id := d.Get("project_id").(string)

	httpResp, err := c.PipelinesApi.DeletePipeline(ctx, org_id, project_id, id, &nextgen.PipelinesApiDeletePipelineOpts{
		HarnessAccount: optional.NewString(c.AccountId),
	})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func buildCreatePipeline(d *schema.ResourceData) nextgen.PipelineCreateRequestBody {
	pipeline := nextgen.PipelineCreateRequestBody{
		Identifier:   d.Get("identifier").(string),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Tags:         helpers.ExpandTags(d.Get("tags").(*schema.Set).List()),
		PipelineYaml: d.Get("yaml").(string),
	}

	if attr, ok := d.GetOk("git_details"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		pipeline.GitDetails = &nextgen.GitCreateDetails{}
		if attr, ok := config["branch_name"]; ok {
			pipeline.GitDetails.BranchName = attr.(string)
		}
		if attr, ok := config["file_path"]; ok {
			pipeline.GitDetails.FilePath = attr.(string)
		}
		if attr, ok := config["commit_message"]; ok {
			pipeline.GitDetails.CommitMessage = attr.(string)
		}
		if attr, ok := config["base_branch"]; ok {
			pipeline.GitDetails.BaseBranch = attr.(string)
		}
		if attr, ok := config["connector_ref"]; ok {
			pipeline.GitDetails.ConnectorRef = attr.(string)
		}
		if attr, ok := config["store_type"]; ok {
			pipeline.GitDetails.StoreType = attr.(string)
		}
		if attr, ok := config["repo_name"]; ok {
			pipeline.GitDetails.RepoName = attr.(string)
		}
		if attr, ok := config["is_harness_code_repo"]; ok {
			pipeline.GitDetails.IsHarnessCodeRepo = attr.(bool)
		}
	}
	return pipeline
}

func buildUpdatePipeline(d *schema.ResourceData) nextgen.PipelineUpdateRequestBody {
	pipeline := nextgen.PipelineUpdateRequestBody{
		Identifier:   d.Get("identifier").(string),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Tags:         helpers.ExpandTags(d.Get("tags").(*schema.Set).List()),
		PipelineYaml: d.Get("yaml").(string),
	}

	if attr, ok := d.GetOk("git_details"); ok {
		configs := attr.([]interface{})
		if len(configs) > 0 {
			config := configs[0].(map[string]interface{})

			pipeline.GitDetails = &nextgen.GitUpdateDetails{}

			if attr, ok := config["branch_name"]; ok {
				pipeline.GitDetails.BranchName = attr.(string)
			}
			if attr, ok := config["commit_message"]; ok {
				pipeline.GitDetails.CommitMessage = attr.(string)
			}
			if attr, ok := config["base_branch"]; ok {
				pipeline.GitDetails.BaseBranch = attr.(string)
			}
			if attr, ok := config["last_object_id"]; ok {
				pipeline.GitDetails.LastObjectId = attr.(string)
			}
			if attr, ok := config["last_commit_id"]; ok {
				pipeline.GitDetails.LastCommitId = attr.(string)
			}
			if attr, ok := config["is_harness_code_repo"]; ok {
				pipeline.GitDetails.IsHarnessCodeRepo = attr.(bool)
			}
		}
	}
	return pipeline
}

// Read response from API out to the stored identifiers
func readPipeline(d *schema.ResourceData, pipeline nextgen.PipelineGetResponseBody, org_id string, project_id string, template_applied bool, store_type optional.String, base_branch optional.String, commit_message optional.String, connector_ref optional.String) {
	d.SetId(pipeline.Identifier)
	d.Set("identifier", pipeline.Identifier)
	d.Set("name", pipeline.Name)
	d.Set("org_id", org_id)
	d.Set("project_id", project_id)
	d.Set("tags", helpers.FlattenTags(pipeline.Tags))
	d.Set("yaml", pipeline.PipelineYaml)
	d.Set("description", pipeline.Description)
	d.Set("template_applied_pipeline_yaml", pipeline.TemplateAppliedPipelineYaml)
	d.Set("template_applied", template_applied)
	if pipeline.GitDetails != nil {
		d.Set("git_details", []interface{}{readGitDetails(pipeline, store_type, base_branch, commit_message, connector_ref)})
	}
}

func readGitDetails(pipeline nextgen.PipelineGetResponseBody, store_type optional.String, base_branch optional.String, commit_message optional.String, connector_ref optional.String) map[string]interface{} {
	git_details := map[string]interface{}{
		"branch_name":    pipeline.GitDetails.BranchName,
		"file_path":      pipeline.GitDetails.FilePath,
		"repo_name":      pipeline.GitDetails.RepoName,
		"last_commit_id": pipeline.GitDetails.CommitId,
		"last_object_id": pipeline.GitDetails.ObjectId,
	}
	if store_type.IsSet() {
		git_details["store_type"] = store_type.Value()
	}
	if base_branch.IsSet() {
		git_details["base_branch"] = base_branch.Value()
	}
	if commit_message.IsSet() {
		git_details["commit_message"] = commit_message.Value()
	}
	if connector_ref.IsSet() {
		git_details["connector_ref"] = connector_ref.Value()
	}
	return git_details
}
