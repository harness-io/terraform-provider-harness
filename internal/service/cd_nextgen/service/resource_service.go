package service

import (
	"context"
	"net/http"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceService() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness project.",

		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceCreateOrUpdate,
		DeleteContext: resourceServiceDelete,
		CreateContext: resourceServiceCreateOrUpdate,
		Importer:      helpers.MultiLevelResourceImporter,

		Schema: map[string]*schema.Schema{
			"yaml": {
				Description:      "Service YAML." + helpers.Descriptions.YamlText.String(),
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: helpers.YamlDiffSuppressFunction,
			},
			"force_delete": {
				Description: "Enable this flag for force deletion of service",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"git_details": {
				Description: "Contains parameters related to Git Experience for remote entities",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch": {
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
						"is_new_branch": {
							Description: "If the branch being created is new",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"load_from_cache": {
							Description: "Load service yaml from catch",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"load_from_fallback_branch": {
							Description: "Load service yaml from fallback branch",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"fetch_resolved_yaml": {
				Description: "to fetch resoled service yaml",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"is_force_import": {
				Description: "force import service from remote even if same file path already exist",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"import_from_git": {
				Description: "import service from git",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}

	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()

	svcParams := getSvcParams(d)
	resp, httpResp, err := c.ServicesApi.GetServiceV2(ctx, id, c.AccountId, svcParams)

	if err != nil {
		return helpers.HandleReadApiError(err, d, httpResp)
	}

	readService(d, resp.Data.Service)

	return nil
}

func resourceServiceCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	var err error
	var resp nextgen.ResponseDtoServiceResponse
	var importResp nextgen.ResponseServiceImportResponseDto
	var httpResp *http.Response
	svc := buildService(d)
	id := d.Id()

	if id == "" {
		if d.Get("import_from_git").(bool) {
			svcParams := svcImportParam(svc, d)
			importResp, httpResp, err = c.ServicesApi.ImportService(ctx, c.AccountId, &svcParams)
		} else {
			svcParams := svcCreateParam(svc, d)
			resp, httpResp, err = c.ServicesApi.CreateServiceV2(ctx, c.AccountId, &svcParams)
		}
	} else {
		svcParams := svcUpdateParam(svc, d)
		resp, httpResp, err = c.ServicesApi.UpdateServiceV2(ctx, c.AccountId, &svcParams)
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	if d.Get("import_from_git").(bool) {
		readImportRes(d, importResp.Data.Identifier)
	} else {
		readService(d, resp.Data.Service)
	}

	return nil
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	_, httpResp, err := c.ServicesApi.DeleteServiceV2(ctx, d.Id(), c.AccountId, &nextgen.ServicesApiDeleteServiceV2Opts{
		OrgIdentifier:     helpers.BuildField(d, "org_id"),
		ProjectIdentifier: helpers.BuildField(d, "project_id"),
		ForceDelete:       helpers.BuildFieldForBoolean(d, "force_delete"),
	})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func buildService(d *schema.ResourceData) *nextgen.ServiceRequest {
	return &nextgen.ServiceRequest{
		Identifier:        d.Get("identifier").(string),
		OrgIdentifier:     d.Get("org_id").(string),
		ProjectIdentifier: d.Get("project_id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Tags:              helpers.ExpandTags(d.Get("tags").(*schema.Set).List()),
		Yaml:              d.Get("yaml").(string),
	}
}

func readService(d *schema.ResourceData, project *nextgen.ServiceResponseDetails) {
	d.SetId(project.Identifier)
	d.Set("identifier", project.Identifier)
	d.Set("org_id", project.OrgIdentifier)
	d.Set("project_id", project.ProjectIdentifier)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("tags", helpers.FlattenTags(project.Tags))
	d.Set("yaml", project.Yaml)
}

func readImportRes(d *schema.ResourceData, identifier string) {
	d.SetId(identifier)
	d.Set("identifier", identifier)
}

func getSvcParams(d *schema.ResourceData) *nextgen.ServicesApiGetServiceV2Opts {
	svcParams := &nextgen.ServicesApiGetServiceV2Opts{
		OrgIdentifier:          helpers.BuildField(d, "org_id"),
		ProjectIdentifier:      helpers.BuildField(d, "project_id"),
		Deleted:                helpers.BuildFieldBool(d, "deleted"),
		FetchResolvedYaml:      helpers.BuildFieldBool(d, "fetch_resolved_yaml"),
		Branch:                 helpers.BuildField(d, "git_details.0.branch"),
		RepoName:               helpers.BuildField(d, "git_details.0.repo_name"),
		LoadFromCache:          helpers.BuildField(d, "git_details.0.load_from_cache"),
		LoadFromFallbackBranch: helpers.BuildFieldBool(d, "git_details.0.load_from_fallback_branch"),
	}
	return svcParams
}

func svcCreateParam(svc *nextgen.ServiceRequest, d *schema.ResourceData) nextgen.ServicesApiCreateServiceV2Opts {
	return nextgen.ServicesApiCreateServiceV2Opts{
		Body:              optional.NewInterface(svc),
		Branch:            helpers.BuildField(d, "git_details.0.branch"),
		FilePath:          helpers.BuildField(d, "git_details.0.file_path"),
		CommitMsg:         helpers.BuildField(d, "git_details.0.commit_message"),
		IsNewBranch:       helpers.BuildFieldBool(d, "git_details.0.is_new_branch"),
		BaseBranch:        helpers.BuildField(d, "git_details.0.base_branch"),
		ConnectorRef:      helpers.BuildField(d, "git_details.0.connector_ref"),
		StoreType:         helpers.BuildField(d, "git_details.0.store_type"),
		RepoName:          helpers.BuildField(d, "git_details.0.repo_name"),
		IsHarnessCodeRepo: helpers.BuildFieldBool(d, "git_details.0.is_harness_code_repo"),
	}
}

func svcUpdateParam(svc *nextgen.ServiceRequest, d *schema.ResourceData) nextgen.ServicesApiUpdateServiceV2Opts {
	return nextgen.ServicesApiUpdateServiceV2Opts{
		Body:              optional.NewInterface(svc),
		Branch:            helpers.BuildField(d, "git_details.0.branch"),
		FilePath:          helpers.BuildField(d, "git_details.0.file_path"),
		CommitMsg:         helpers.BuildField(d, "git_details.0.commit_message"),
		IsNewBranch:       helpers.BuildFieldBool(d, "git_details.0.is_new_branch"),
		BaseBranch:        helpers.BuildField(d, "git_details.0.base_branch"),
		ConnectorRef:      helpers.BuildField(d, "git_details.0.connector_ref"),
		StoreType:         helpers.BuildField(d, "git_details.0.store_type"),
		IfMatch:           helpers.BuildField(d, "if_match"),
		LastObjectId:      helpers.BuildField(d, "git_details.0.last_object_id"),
		LastCommitId:      helpers.BuildField(d, "git_details.0.last_commit_id"),
		IsHarnessCodeRepo: helpers.BuildFieldBool(d, "git_details.0.is_harness_code_repo"),
	}
}

func svcImportParam(svc *nextgen.ServiceRequest, d *schema.ResourceData) nextgen.ServicesApiImportServiceOpts {
	return nextgen.ServicesApiImportServiceOpts{
		OrgIdentifier:     helpers.BuildField(d, "org_id"),
		ProjectIdentifier: helpers.BuildField(d, "project_id"),
		ServiceIdentifier: helpers.BuildField(d, "identifier"),
		Branch:            helpers.BuildField(d, "git_details.0.branch"),
		FilePath:          helpers.BuildField(d, "git_details.0.file_path"),
		ConnectorRef:      helpers.BuildField(d, "git_details.0.connector_ref"),
		IsHarnessCodeRepo: helpers.BuildFieldBool(d, "git_details.0.is_harness_code_repo"),
		RepoName:          helpers.BuildField(d, "git_details.0.repo_name"),
		IsForceImport:     helpers.BuildFieldBool(d, "git_details.0.is_force_import"),
	}
}
