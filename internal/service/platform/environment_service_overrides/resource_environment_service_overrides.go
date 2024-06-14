package environment_service_overrides

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

func ResourceEnvironmentServiceOverrides() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness environment service overrides.",

		ReadContext:   resourceEnvironmentServiceOverridesRead,
		UpdateContext: resourceEnvironmentServiceOverridesCreateOrUpdate,
		DeleteContext: resourceEnvironmentServiceOverridesDelete,
		CreateContext: resourceEnvironmentServiceOverridesCreateOrUpdate,
		Importer:      helpers.ServiceOverrideResourceImporter,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Description: "identifier of the service overrides.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"service_id": {
				Description: "The service ID to which the overrides applies.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"env_id": {
				Description: "The env ID to which the overrides associated.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"yaml": {
				Description:      "Environment Service Overrides YAML." + helpers.Descriptions.YamlText.String(),
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: helpers.YamlDiffSuppressFunction,
			},
		},
	}

	SetScopedResourceSchemaForServiceOverride(resource.Schema)

	return resource
}

func resourceEnvironmentServiceOverridesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	envId := d.Get("env_id").(string)

	resp, httpResp, err := c.EnvironmentsApi.GetServiceOverridesList(ctx, c.AccountId, envId,
		&nextgen.EnvironmentsApiGetServiceOverridesListOpts{
			ServiceIdentifier: helpers.BuildField(d, "service_id"),
			OrgIdentifier:     helpers.BuildField(d, "org_id"),
			ProjectIdentifier: helpers.BuildField(d, "project_id"),
		})

	if err != nil {
		return helpers.HandleReadApiError(err, d, httpResp)
	}

	if resp.Data == nil || len(resp.Data.Content) == 0 {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}

	readEnvironmentServiceOverridesList(d, resp.Data)

	return nil
}

func resourceEnvironmentServiceOverridesCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	var err error
	var resp nextgen.ResponseDtoServiceOverrideResponse
	var httpResp *http.Response
	env := buildEnvironmentServiceOverride(d)

	resp, httpResp, err = c.EnvironmentsApi.UpsertServiceOverride(ctx, c.AccountId, &nextgen.EnvironmentsApiUpsertServiceOverrideOpts{
		Body: optional.NewInterface(env),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	if &resp == nil || resp.Data == nil {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}

	readEnvironmentServiceOverrides(d, resp.Data)

	return nil
}

func resourceEnvironmentServiceOverridesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	serviceRef := d.Get("service_id").(string)
	envId := d.Get("env_id").(string)

	_, httpResp, err := c.EnvironmentsApi.DeleteServiceOverride(ctx, c.AccountId, &nextgen.EnvironmentsApiDeleteServiceOverrideOpts{
		EnvironmentIdentifier: optional.NewString(envId),
		ServiceIdentifier:     optional.NewString(serviceRef),
		OrgIdentifier:         helpers.BuildField(d, "org_id"),
		ProjectIdentifier:     helpers.BuildField(d, "project_id"),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func buildEnvironmentServiceOverride(d *schema.ResourceData) *nextgen.ServiceOverrideRequest {
	return &nextgen.ServiceOverrideRequest{
		OrgIdentifier:         d.Get("org_id").(string),
		ProjectIdentifier:     d.Get("project_id").(string),
		EnvironmentIdentifier: d.Get("env_id").(string),
		ServiceIdentifier:     d.Get("service_id").(string),
		Yaml:                  d.Get("yaml").(string),
	}
}

func readEnvironmentServiceOverridesList(d *schema.ResourceData, env *nextgen.PageResponseServiceOverrideResponse) {
	ServiceOverrideList := env.Content
	for _, value := range ServiceOverrideList {
		readEnvironmentServiceOverrides(d, &value)
	}
}

func readEnvironmentServiceOverrides(d *schema.ResourceData, so *nextgen.ServiceOverrideResponse) {
	serviceOverrideID := so.ServiceRef + "-" + so.EnvironmentRef
	d.SetId(serviceOverrideID)
	d.Set("identifier", serviceOverrideID)
	d.Set("org_id", so.OrgIdentifier)
	d.Set("project_id", so.ProjectIdentifier)
	d.Set("env_id", so.EnvironmentRef)
	d.Set("service_id", so.ServiceRef)
	d.Set("yaml", so.Yaml)
}

func SetScopedResourceSchemaForServiceOverride(s map[string]*schema.Schema) {
	s["project_id"] = helpers.GetProjectIdSchema(helpers.SchemaFlagTypes.Optional)
	s["org_id"] = helpers.GetOrgIdSchema(helpers.SchemaFlagTypes.Optional)
}
