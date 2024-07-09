package dbschema

import (
	"context"
	"net/http"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/dbops"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDBSchema() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a Harness DBDevOps Schema.",

		ReadContext: dataSourceDBSchemaRead,

		Schema: map[string]*schema.Schema{
			"service": {
				Description: "The service associated with schema",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"schema_source": {
				Description: "Provides a connector and path at which to find the database schema representation",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector": {
							Description: "Connector to repository at which to find details about the database schema",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"location": {
							Description: "The path within the specified repository at which to find details about the database schema",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"repo": {
							Description: "If connector url is of account, which repository to connect to using the connector",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}

	helpers.SetProjectLevelDataSourceSchemaIdentifierRequired(resource.Schema)

	return resource
}

func dataSourceDBSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetDBOpsClientWithContext(ctx)

	var err error
	var dbSchema dbops.DbSchemaOut
	var httpResp *http.Response

	id := d.Get("identifier").(string)

	localVarOptionals := dbops.DatabaseSchemaApiV1GetProjDbSchemaOpts{
		HarnessAccount: optional.NewString(meta.(*internal.Session).AccountId),
	}
	dbSchema, httpResp, err = c.DatabaseSchemaApi.V1GetProjDbSchema(ctx, d.Get("org_id").(string), d.Get("project_id").(string), id, &localVarOptionals)

	if err != nil {
		return helpers.HandleDBOpsApiError(err, d, httpResp)
	}

	readDataSourceDBSchema(d, &dbSchema)

	return nil
}

func readDataSourceDBSchema(d *schema.ResourceData, dbSchema *dbops.DbSchemaOut) {
	d.SetId(dbSchema.Identifier)
	d.Set("identifier", dbSchema.Identifier)
	d.Set("name", dbSchema.Name)
	d.Set("tags", helpers.FlattenTags(dbSchema.Tags))
	d.Set("service", dbSchema.Service)
	d.Set("schema_source.0.location", dbSchema.Changelog.Location)
	d.Set("schema_source.0.repo", dbSchema.Changelog.Repo)
	d.Set("schema_source.0.connector", dbSchema.Changelog.Connector)
}
