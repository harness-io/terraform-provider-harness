package pipeline

import (
	"context"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourcePipeline() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a Harness pipeline.",

		ReadContext: dataSourcePipelineRead,

		Schema: map[string]*schema.Schema{
			"yaml": {
				Description: "YAML of the pipeline.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}

	helpers.SetProjectLevelDataSourceSchema(resource.Schema)

	return resource
}

func dataSourcePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	resp, _, err := c.PipelinesApi.GetPipeline(ctx,
		c.AccountId,
		d.Get("org_id").(string),
		d.Get("project_id").(string),
		d.Get("identifier").(string),
		&nextgen.PipelinesApiGetPipelineOpts{},
	)

	if err != nil {
		return helpers.HandleApiError(err, d)
	}

	readPipeline(d, resp.Data)

	return nil
}
