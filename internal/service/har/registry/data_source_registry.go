package registry

import (
	"context"
	"github.com/harness/harness-go-sdk/harness/har"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func DataSourceRegistry() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating and managing Harness Registries.",

		ReadContext: dataSourceRegistryRead,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Description: "Unique identifier of the registry",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Description of the registry",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"parent_ref": {
				Description: "Parent reference for the registry",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"space_ref": {
				Description: "Space reference for the registry",
				Type:        schema.TypeString,
				Required:    true,
			},
			"package_type": {
				Description: "Type of package (DOCKER, MAVEN, etc.)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"config": {
				Description: "Configuration for the registry",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of registry (VIRTUAL, UPSTREAM)",
							Type:        schema.TypeString,
							Required:    true,
						},
						"auth": {
							Description: "Authentication configuration for UPSTREAM type",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_type": {
										Description: "Type of authentication (UserPassword, Anonymous)",
										Type:        schema.TypeString,
										Required:    true,
									},
									"user_password": {
										Description: "User password authentication details",
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"secret_identifier": {
													Description: "Secret identifier",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"secret_space_id": {
													Description: "Secret space ID",
													Type:        schema.TypeInt,
													Optional:    true,
												},
												"secret_space_path": {
													Description: "Secret space path",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"user_name": {
													Description: "User name",
													Type:        schema.TypeString,
													Required:    true,
												},
											},
										},
									},
								},
							},
						},
						"source": {
							Description: "Source of the upstream",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"url": {
							Description: "URL of the upstream",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"url": {
				Description: "URL of the registry",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "Timestamp when the registry was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"allowed_pattern": {
				Description: "Allowed pattern for the registry",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"blocked_pattern": {
				Description: "Blocked pattern for the registry",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetHarClientWithContext(ctx)
	if c == nil {
		return diag.Errorf("Harness client is not initialized. Check provider configuration.")
	}

	var registry *har.Registry
	var err error
	var resp har.InlineResponse201
	var httpResp *http.Response

	id := d.Get("identifier").(string)
	registryRef := d.Get("space_ref").(string) + "/" + id

	if id != "" && registryRef != "" {
		resp, httpResp, err = c.RegistriesApi.GetRegistry(ctx, registryRef)
		if err != nil {
			return helpers.HandleReadApiError(err, d, httpResp)
		}

		registry = resp.Data
	} else {
		return diag.Errorf("Registry identifier and Space reference are required to read the registry.")
	}

	if registry.Identifier == "" {
		return diag.Errorf("Registry not found.")
	}

	readRegistry(d, registry)
	return nil
}