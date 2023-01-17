package input_set

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

func ResourceInputSet() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness InputSet.",

		ReadContext:   resourceInputSetRead,
		UpdateContext: resourceInputSetCreateOrUpdate,
		CreateContext: resourceInputSetCreateOrUpdate,
		DeleteContext: resourceInputSetDelete,
		Importer:      helpers.PipelineResourceImporter,

		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Description: "Identifier of the pipeline",
				Type:        schema.TypeString,
				Required:    true,
			},
			"yaml": {
				Description: "Input Set YAML",
				Type:        schema.TypeString,
				Required:    true,
			},
			"git_details": {
				Description: "Contains parameters related to creating an Entity for Git Experience.",
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
						"commit_message": {
							Description: "Commit message used for the merge commit.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"base_branch": {
							Description: "Name of the default branch (this checks out a new branch titled by branch_name).",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"connector_ref": {
							Description: "Identifier of the Harness Connector used for CRUD operations on the Entity.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"store_type": {
							Description:  "Specifies whether the Entity is to be stored in Git or not. Possible values: INLINE, REMOTE.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"INLINE", "REMOTE"}, false),
						},
						"repo_name": {
							Description: "Name of the repository.",
							Type:        schema.TypeString,
							Optional:    true,
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
						"parent_entity_connector_ref": {
							Description: "Connector reference for Parent Entity (Pipeline).",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"parent_entity_repo_name": {
							Description: "Repository name for Parent Entity (Pipeline).",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
	helpers.SetProjectLevelResourceSchema(resource.Schema)

	return resource
}

func resourceInputSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetClientWithContext(ctx)

	id := d.Get("identifier").(string)

	orgId := d.Get("org_id").(string)

	projectId := d.Get("project_id").(string)

	pipelineId := d.Get("pipeline_id").(string)

	var branch_name optional.String

	branch_name = helpers.BuildField(d, "git_details.0.branch_name")
	parent_entity_connector_ref := helpers.BuildField(d, "git_details.0.parent_entity_connector_ref")
	parent_entity_repo_name := helpers.BuildField(d, "git_details.0.parent_entity_repo_name")

	var store_type = helpers.BuildField(d, "git_details.0.store_type")
	var base_branch = helpers.BuildField(d, "git_details.0.base_branch")
	var commit_message = helpers.BuildField(d, "git_details.0.commit_message")
	var connector_ref = helpers.BuildField(d, "git_details.0.connector_ref")

	resp, httpResp, err := c.InputSetsApi.GetInputSet(ctx, orgId, projectId, id, pipelineId, &nextgen.InputSetsApiGetInputSetOpts{
		HarnessAccount:           optional.NewString(c.AccountId),
		BranchName:               branch_name,
		ParentEntityConnectorRef: parent_entity_connector_ref,
		ParentEntityRepoName:     parent_entity_repo_name,
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readInputSet(d, &resp, pipelineId, store_type, base_branch, commit_message, connector_ref)

	return nil

}

func resourceInputSetCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetClientWithContext(ctx)

	var err error
	var resp nextgen.InputSetResponseBody
	var httpResp *http.Response

	id := d.Id()
	orgIdentifier := d.Get("org_id").(string)
	projectIdentifier := d.Get("project_id").(string)
	pipelineIdentifier := d.Get("pipeline_id").(string)

	var store_type optional.String
	var base_branch optional.String
	var commit_message optional.String
	var connector_ref optional.String

	if id == "" {
		inputSet := buildCreateInputSet(d)
		if inputSet.GitDetails != nil {
			base_branch = optional.NewString(inputSet.GitDetails.BaseBranch)
			store_type = optional.NewString(inputSet.GitDetails.StoreType)
			commit_message = optional.NewString(inputSet.GitDetails.CommitMessage)
			connector_ref = optional.NewString(inputSet.GitDetails.ConnectorRef)
		}

		resp, httpResp, err = c.InputSetsApi.CreateInputSet(ctx, inputSet, pipelineIdentifier, orgIdentifier, projectIdentifier, &nextgen.InputSetsApiCreateInputSetOpts{
			HarnessAccount: optional.NewString(c.AccountId),
		})
	} else {
		inputSet := buildUpdateInputSet(d)
		if inputSet.GitDetails != nil {
			base_branch = optional.NewString(inputSet.GitDetails.BaseBranch)
			commit_message = optional.NewString(inputSet.GitDetails.CommitMessage)
		}
		store_type = helpers.BuildField(d, "git_details.0.store_type")
		connector_ref = helpers.BuildField(d, "git_details.0.connector_ref")

		resp, httpResp, err = c.InputSetsApi.UpdateInputSet(ctx, inputSet, pipelineIdentifier, orgIdentifier, projectIdentifier, id, &nextgen.InputSetsApiUpdateInputSetOpts{
			HarnessAccount: optional.NewString(c.AccountId),
		})
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readInputSet(d, &resp, pipelineIdentifier, store_type, base_branch, commit_message, connector_ref)

	return nil
}

func resourceInputSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetClientWithContext(ctx)
	orgIdentifier := helpers.BuildField(d, "org_id").Value()
	projectIdentifier := helpers.BuildField(d, "project_id").Value()
	pipelineIdentifier := helpers.BuildField(d, "pipeline_id").Value()

	httpResp, err := c.InputSetsApi.DeleteInputSet(ctx, orgIdentifier, projectIdentifier, d.Id(), pipelineIdentifier, &nextgen.InputSetsApiDeleteInputSetOpts{
		HarnessAccount: optional.NewString(c.AccountId),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func buildCreateInputSet(d *schema.ResourceData) nextgen.InputSetCreateRequestBody {
	inputSet := nextgen.InputSetCreateRequestBody{
		Slug:         d.Get("identifier").(string),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Tags:         helpers.ExpandTags(d.Get("tags").(*schema.Set).List()),
		InputSetYaml: d.Get("yaml").(string),
	}

	if attr, ok := d.GetOk("git_details"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		inputSet.GitDetails = &nextgen.GitCreateDetails{}
		if attr, ok := config["branch_name"]; ok {
			inputSet.GitDetails.BranchName = attr.(string)
		}
		if attr, ok := config["file_path"]; ok {
			inputSet.GitDetails.FilePath = attr.(string)
		}
		if attr, ok := config["commit_message"]; ok {
			inputSet.GitDetails.CommitMessage = attr.(string)
		}
		if attr, ok := config["base_branch"]; ok {
			inputSet.GitDetails.BaseBranch = attr.(string)
		}
		if attr, ok := config["connector_ref"]; ok {
			inputSet.GitDetails.ConnectorRef = attr.(string)
		}
		if attr, ok := config["store_type"]; ok {
			inputSet.GitDetails.StoreType = attr.(string)
		}
		if attr, ok := config["repo_name"]; ok {
			inputSet.GitDetails.RepoName = attr.(string)
		}
	}
	return inputSet
}

func buildUpdateInputSet(d *schema.ResourceData) nextgen.InputSetUpdateRequestBody {
	inputSet := nextgen.InputSetUpdateRequestBody{
		Slug:         d.Get("identifier").(string),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Tags:         helpers.ExpandTags(d.Get("tags").(*schema.Set).List()),
		InputSetYaml: d.Get("yaml").(string),
	}

	if attr, ok := d.GetOk("git_details"); ok {
		configs := attr.([]interface{})
		if len(configs) > 0 {
			config := configs[0].(map[string]interface{})

			inputSet.GitDetails = &nextgen.InputSetGitUpdateDetails{}

			if attr, ok := config["branch_name"]; ok {
				inputSet.GitDetails.BranchName = attr.(string)
			}
			if attr, ok := config["commit_message"]; ok {
				inputSet.GitDetails.CommitMessage = attr.(string)
			}
			if attr, ok := config["base_branch"]; ok {
				inputSet.GitDetails.BaseBranch = attr.(string)
			}
			if attr, ok := config["last_object_id"]; ok {
				inputSet.GitDetails.LastObjectId = attr.(string)
			}
			if attr, ok := config["last_commit_id"]; ok {
				inputSet.GitDetails.LastCommitId = attr.(string)
			}
			if attr, ok := config["parent_entity_connector_ref"]; ok {
				inputSet.GitDetails.ParentEntityConnectorRef = attr.(string)
			}
			if attr, ok := config["parent_entity_repo_name"]; ok {
				inputSet.GitDetails.ParentEntityRepoName = attr.(string)
			}
		}
	}
	return inputSet
}

func readInputSet(d *schema.ResourceData, inputSet *nextgen.InputSetResponseBody, pipelineId string, store_type optional.String, base_branch optional.String, commit_message optional.String, connector_ref optional.String) {
	d.SetId(inputSet.Slug)
	d.Set("identifier", inputSet.Slug)
	d.Set("name", inputSet.Name)
	d.Set("description", inputSet.Description)
	d.Set("tags", helpers.FlattenTags(inputSet.Tags))
	d.Set("org_id", inputSet.Org)
	d.Set("project_id", inputSet.Project)
	d.Set("pipeline_id", pipelineId)
	d.Set("yaml", inputSet.InputSetYaml)
	if inputSet.GitDetails != nil {
		d.Set("git_details", []interface{}{readGitDetails(inputSet, store_type, base_branch, commit_message, connector_ref)})
	}
}

func readGitDetails(inputSet *nextgen.InputSetResponseBody, store_type optional.String, base_branch optional.String, commit_message optional.String, connector_ref optional.String) map[string]interface{} {
	git_details := map[string]interface{}{
		"branch_name":    inputSet.GitDetails.BranchName,
		"file_path":      inputSet.GitDetails.FilePath,
		"repo_name":      inputSet.GitDetails.RepoName,
		"last_commit_id": inputSet.GitDetails.CommitId,
		"last_object_id": inputSet.GitDetails.ObjectId,
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
