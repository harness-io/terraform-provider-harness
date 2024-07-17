package notification_rule

import (
	"context"

	hh "github.com/harness/harness-go-sdk/harness/helpers"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNotificationRuleService() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a Notification Rule.",

		ReadContext: dataSourceNotificationRuleRead,

		Schema: map[string]*schema.Schema{
			"org_id": {
				Description: "Identifier of the organization in which the Notification Rule is configured.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": {
				Description: "Identifier of the project in which the Notification Rule is configured.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"identifier": {
				Description: "Identifier of the Notification Rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}

	return resource
}

func dataSourceNotificationRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	ctx = context.WithValue(ctx, nextgen.ContextAccessToken, hh.EnvVars.BearerToken.Get())
	var accountIdentifier, orgIdentifier, projectIdentifier string
	identifier := d.Get("identifier").(string)
	accountIdentifier = c.AccountId
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	resp, httpResp, err := c.SrmNotificationApiService.GetSrmNotification(ctx, identifier, accountIdentifier, orgIdentifier, projectIdentifier)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readNotificationRule(d, &resp.Resource)
	return nil
}
