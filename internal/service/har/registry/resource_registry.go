package registry

import (
	"context"
	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/har"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func ResourceRegistry() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating and managing Harness Registries.",

		ReadContext:   resourceRegistryRead,
		CreateContext: resourceRegistryCreateOrUpdate,
		UpdateContext: resourceRegistryCreateOrUpdate,
		DeleteContext: resourceRegistryDelete,

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
				Optional:    true,
			},
			"package_type": {
				Description: "Type of package (DOCKER, MAVEN, etc.)",
				Type:        schema.TypeString,
				Required:    true,
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

func resourceRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetHarClientWithContext(ctx)

	registryRef := d.Get("parent_ref").(string) + "/" + d.Get("identifier").(string)
	resp, httpResp, err := c.RegistriesApi.GetRegistry(ctx, registryRef)

	if err != nil {
		return helpers.HandleReadApiError(err, d, httpResp)
	}

	readRegistry(d, resp.Data)
	return nil
}

func resourceRegistryCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetHarClientWithContext(ctx)
	if c == nil {
		return diag.Errorf("Harness client is not initialized. Check provider configuration.")
	}

	var err error
	var resp har.InlineResponse201
	var httpResp *http.Response

	registry := buildRegistry(d)
	spaceRef := d.Get("space_ref").(string)

	if d.Id() == "" {
		resp, httpResp, err = c.RegistriesApi.CreateRegistry(ctx, &har.RegistriesApiCreateRegistryOpts{
			Body: optional.NewInterface(registry), SpaceRef: optional.NewString(spaceRef),
		})
	} else {
		resp, httpResp, err = c.RegistriesApi.ModifyRegistry(ctx, d.Id(), &har.RegistriesApiModifyRegistryOpts{
			Body: optional.NewInterface(registry),
		})
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readRegistry(d, resp.Data)
	return nil
}

func resourceRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetHarClientWithContext(ctx)

	registryRef := d.Get("parent_ref").(string) + "/" + d.Get("identifier").(string)

	_, httpResp, err := c.RegistriesApi.DeleteRegistry(ctx, registryRef)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func buildRegistry(d *schema.ResourceData) *har.RegistryRequest {
	registry := &har.RegistryRequest{}

	if attr, ok := d.GetOk("identifier"); ok {
		registry.Identifier = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		registry.Description = attr.(string)
	}

	if attr, ok := d.GetOk("parent_ref"); ok {
		registry.ParentRef = attr.(string)
	}

	if attr, ok := d.GetOk("package_type"); ok {
		pt := har.PackageType(attr.(string))
		registry.PackageType = &pt
	}

	if attr, ok := d.GetOk("config"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		if t, ok := config["type"]; ok {
			registry.Config = &har.RegistryConfig{
				Type_: func(t string) *har.RegistryType { r := har.RegistryType(t); return &r }(t.(string)),
			}
		}
	}

	if attr, ok := d.GetOk("allowed_pattern"); ok {
		registry.AllowedPattern = convertListToString(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("blocked_pattern"); ok {
		registry.BlockedPattern = convertListToString(attr.([]interface{}))
	}

	return registry
}

func readRegistry(d *schema.ResourceData, registry *har.Registry) {
	d.SetId(registry.Identifier)
	d.Set("identifier", registry.Identifier)
	d.Set("description", registry.Description)
	d.Set("url", registry.Url)
	d.Set("package_type", registry.PackageType)
	d.Set("created_at", registry.CreatedAt)
	d.Set("allowed_pattern", registry.AllowedPattern)
	d.Set("blocked_pattern", registry.BlockedPattern)

	if registry.Config != nil {
		d.Set("config", []interface{}{
			map[string]interface{}{
				"type": registry.Config.Type_,
			},
		})
	}
}

func convertListToString(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}