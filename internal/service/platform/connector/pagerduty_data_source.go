package connector

import (
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceConnectorPagerDuty() *schema.Resource {
	resource := &schema.Resource{
		Description: "Datasource for looking up a PagerDuty connector.",
		ReadContext: resourceConnectorPagerDutyRead,

		Schema: map[string]*schema.Schema{
			"api_token_ref": {
				Description: "Reference to the Harness secret containing the api token." + secret_ref_text,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"delegate_selectors": {
				Description: "Tags to filter delegates for connection.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	helpers.SetMultiLevelDatasourceSchema(resource.Schema)

	return resource
}
