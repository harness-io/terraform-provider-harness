package registry

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/har"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strings"
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
				Description: "Type of package (DOCKER, HELM, etc.)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config": {
				Description: "Configuration for the registry",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of registry (VIRTUAL or UPSTREAM)",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								validTypes := map[string]bool{
									"VIRTUAL":  true,
									"UPSTREAM": true,
								}
								if !validTypes[v] {
									errs = append(errs, fmt.Errorf("config type must be either 'VIRTUAL' or 'UPSTREAM', got: %s", v))
								}
								return
							},
						},
						// Virtual Config
						"upstream_proxies": {
							Description: "List of upstream proxies for VIRTUAL registry type",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},

						// Upstream Config
						"auth": {
							Description: "Authentication configuration for UPSTREAM registry type",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_type": {
										Description: "Type of authentication (UserPassword)",
										Type:        schema.TypeString,
										Required:    true,
									},
									"secret_identifier": {
										Description: "Secret identifier for UserPassword auth type",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"secret_space_path": {
										Description: "Secret space path for UserPassword auth type",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"user_name": {
										Description: "User name  for UserPassword auth type",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"auth_type": {
							Description: "Type of authentication for UPSTREAM registry type (UserPassword, Anonymous)",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"source": {
							Description: "Source of the upstream (only for UPSTREAM type)",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"url": {
							Description: "URL of the upstream (required if type=UPSTREAM & package_type=HELM)",
							Type:        schema.TypeString,
							Optional:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v, _ := val.(string)
								// Validate URL format (basic check)
								if v != "" && !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
									errs = append(errs, fmt.Errorf("invalid URL format, must start with http:// or https://"))
								}
								return
							},
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
		configList := attr.([]interface{})
		if len(configList) > 0 { // Ensure config is not empty before accessing index
			config := configList[0].(map[string]interface{})

			if t, ok := config["type"].(string); ok {
				registryType := har.RegistryType(t)
				registry.Config = &har.RegistryConfig{Type_: &registryType}

				// Handle VIRTUAL type
				if registryType == har.VIRTUAL_RegistryType {
					if proxies, ok := config["upstream_proxies"].([]interface{}); ok {
						for _, proxy := range proxies {
							registry.Config.VirtualConfig.UpstreamProxies = append(
								registry.Config.VirtualConfig.UpstreamProxies, proxy.(string))
						}
					}
				}

				// Handle UPSTREAM type
				if registryType == har.UPSTREAM_RegistryType {
					upstreamConfig := &har.UpstreamConfig{}

					if source, ok := config["source"].(string); ok {
						upstreamConfig.Source = source
					}
					if url, ok := config["url"].(string); ok {
						upstreamConfig.Url = url
					} else {
						upstreamConfig.Url = ""
					}

					// Handle authType at the top level
					if authType, ok := config["auth_type"].(string); ok {
						upstreamConfig.AuthType = (*har.AuthType)(&authType)
					}

					// Handle Authentication block
					if authAttr, ok := config["auth"].([]interface{}); ok && len(authAttr) > 0 {
						authConfig := authAttr[0].(map[string]interface{}) // Extract first element as map

						if authType, ok := authConfig["auth_type"].(string); ok {
							upstreamConfig.AuthType = (*har.AuthType)(&authType)

							if authType == "UserPassword" {
								userPassword := har.UserPassword{}

								if val, ok := authConfig["user_name"].(string); ok {
									userPassword.UserName = val
								}
								if val, ok := authConfig["secret_identifier"].(string); ok {
									userPassword.SecretIdentifier = val
								}
								if val, ok := authConfig["secret_space_path"].(string); ok {
									userPassword.SecretSpacePath = val
								}

								upstreamConfig.Auth = &har.OneOfUpstreamConfigAuth{
									UserPassword: userPassword,
									AuthType:     (*har.AuthType)(&authType),
								}
							} else if authType == "Anonymous" {
								upstreamConfig.Auth = &har.OneOfUpstreamConfigAuth{
									Anonymous: har.Anonymous{},
								}
							}
						}
					}

					registry.Config.UpstreamConfig = *upstreamConfig
				}
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
		configMap := map[string]interface{}{
			"type": registry.Config.Type_,
		}

		// Handle VIRTUAL type configuration
		if registry.Config.Type_ != nil && *registry.Config.Type_ == har.VIRTUAL_RegistryType {
			if len(registry.Config.VirtualConfig.UpstreamProxies) > 0 {
				configMap["upstream_proxies"] = registry.Config.VirtualConfig.UpstreamProxies
			}
		}

		// Handle UPSTREAM type configuration
		if registry.Config.Type_ != nil && *registry.Config.Type_ == har.UPSTREAM_RegistryType {
			if registry.Config.UpstreamConfig.Source != "" {
				configMap["source"] = registry.Config.UpstreamConfig.Source
			}
			if registry.Config.UpstreamConfig.Url != "" {
				configMap["url"] = registry.Config.UpstreamConfig.Url
			}
			if registry.Config.UpstreamConfig.AuthType != nil {
				configMap["auth_type"] = *registry.Config.UpstreamConfig.AuthType
			}

			// Handle Authentication
			if registry.Config.UpstreamConfig.Auth != nil {
				authMap := map[string]interface{}{}
				if registry.Config.UpstreamConfig.Auth.UserPassword.UserName != "" {
					authMap["auth_type"] = "UserPassword"
					authMap["user_name"] = registry.Config.UpstreamConfig.Auth.UserPassword.UserName
					authMap["secret_identifier"] = registry.Config.UpstreamConfig.Auth.UserPassword.SecretIdentifier
					authMap["secret_space_path"] = registry.Config.UpstreamConfig.Auth.UserPassword.SecretSpacePath
				} else {
					authMap["auth_type"] = "Anonymous"
				}
				configMap["auth"] = authMap
			}
		}

		// Set the updated config map to Terraform state
		d.Set("config", []interface{}{configMap})
	}
}

func convertListToString(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}
