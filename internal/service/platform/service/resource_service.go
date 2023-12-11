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
		},
	}

	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()

	resp, httpResp, err := c.ServicesApi.GetServiceV2(ctx, id, c.AccountId, &nextgen.ServicesApiGetServiceV2Opts{
		OrgIdentifier:     helpers.BuildField(d, "org_id"),
		ProjectIdentifier: helpers.BuildField(d, "project_id"),
	})

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
	var httpResp *http.Response
	svc := buildService(d)
	id := d.Id()

	if id == "" {
		resp, httpResp, err = c.ServicesApi.CreateServiceV2(ctx, c.AccountId, &nextgen.ServicesApiCreateServiceV2Opts{
			Body: optional.NewInterface(svc),
		})
	} else {
		resp, httpResp, err = c.ServicesApi.UpdateServiceV2(ctx, c.AccountId, &nextgen.ServicesApiUpdateServiceV2Opts{
			Body: optional.NewInterface(svc),
		})
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readService(d, resp.Data.Service)

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
