package cloudprovider

import (
	"context"
	"fmt"

	sdk "github.com/harness-io/harness-go-sdk"
	"github.com/harness-io/harness-go-sdk/harness/cd/cac"
	"github.com/harness-io/terraform-provider-harness/helpers"
	"github.com/harness-io/terraform-provider-harness/internal/service/cd/usagescope"
	"github.com/harness-io/terraform-provider-harness/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCloudProviderAzure() *schema.Resource {

	providerSchema := map[string]*schema.Schema{
		"environment_type": {
			Description:  fmt.Sprintf("The type of environment. Valid options are %s", cac.AzureEnvironmentTypesSlice),
			Type:         schema.TypeString,
			Optional:     true,
			Default:      cac.AzureEnvironmentTypes.AzureGlobal.String(),
			ValidateFunc: validation.StringInSlice(cac.AzureEnvironmentTypesSlice, false),
		},
		"client_id": {
			Description: "The client id for the Azure application",
			Type:        schema.TypeString,
			Required:    true,
		},
		"tenant_id": {
			Description: "The tenant id for the Azure application",
			Type:        schema.TypeString,
			Required:    true,
		},
		"key": {
			Description: "The Name of the Harness secret containing the key for the Azure application",
			Type:        schema.TypeString,
			Required:    true,
		},
	}

	// usage_scope is not supported because the scope will always be inherited from the secret defined in `key`
	commonSchema := commonCloudProviderSchema()
	delete(commonSchema, "usage_scope")

	helpers.MergeSchemas(commonSchema, providerSchema)

	return &schema.Resource{
		Description:   utils.ConfigAsCodeDescription("Resource for creating an Azure cloud provider."),
		CreateContext: resourceCloudProviderAzureCreateOrUpdate,
		ReadContext:   resourceCloudProviderAzureRead,
		UpdateContext: resourceCloudProviderAzureCreateOrUpdate,
		DeleteContext: resourceCloudProviderDelete,

		Schema: providerSchema,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCloudProviderAzureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sdk.Session)

	cp := &cac.AzureCloudProvider{}
	if err := c.CDClient.ConfigAsCodeClient.GetCloudProviderById(d.Id(), cp); err != nil {
		return diag.FromErr(err)
	} else if cp.IsEmpty() {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}

	return readCloudProviderAzure(c, d, cp)
}

func resourceCloudProviderAzureCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sdk.Session)

	var input *cac.AzureCloudProvider
	var err error

	if d.IsNewResource() {
		input = cac.NewEntity(cac.ObjectTypes.AzureCloudProvider).(*cac.AzureCloudProvider)
	} else {
		input = &cac.AzureCloudProvider{}
		if err = c.CDClient.ConfigAsCodeClient.GetCloudProviderById(d.Id(), input); err != nil {
			return diag.FromErr(err)
		} else if input.IsEmpty() {
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
	}

	input.Name = d.Get("name").(string)
	input.AzureEnvironmentType = cac.AzureEnvironmentType(d.Get("environment_type").(string))
	input.ClientId = d.Get("client_id").(string)
	input.TenantId = d.Get("tenant_id").(string)
	input.Key = &cac.SecretRef{
		Name: d.Get("key").(string),
	}

	if input.UsageRestrictions == nil {
		input.UsageRestrictions = &cac.UsageRestrictions{}
	}

	scopes := d.Get("usage_scope")
	if scopes != nil {
		if err := usagescope.ExpandUsageRestrictions(c, scopes.(*schema.Set).List(), input.UsageRestrictions); err != nil {
			return diag.FromErr(err)
		}
	}

	cp, err := c.CDClient.ConfigAsCodeClient.UpsertAzureCloudProvider(input)

	if err != nil {
		return diag.FromErr(err)
	}

	return readCloudProviderAzure(c, d, cp)
}

func readCloudProviderAzure(c *sdk.Session, d *schema.ResourceData, cp *cac.AzureCloudProvider) diag.Diagnostics {
	d.SetId(cp.Id)
	d.Set("name", cp.Name)
	d.Set("environment_type", cp.AzureEnvironmentType)
	d.Set("client_id", cp.ClientId)
	d.Set("tenant_id", cp.TenantId)
	d.Set("key", cp.Key.Name)

	scope, err := usagescope.FlattenUsageRestrictions(c, cp.UsageRestrictions)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(scope) > 0 {
		d.Set("usage_scope", scope)
	}

	return nil
}
