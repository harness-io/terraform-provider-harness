package repository_certificates

import (
	"context"

	"github.com/antihax/optional"
	hh "github.com/harness/harness-go-sdk/harness/helpers"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGitOpsRepoCert() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a GitOps Repository Certificate. It fetches all the certificates that are added to the provided agent.",

		ReadContext: dataSourceGitopsRepoCertRead,

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Description: "Agent identifier of the GitOps repository certificate.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"account_id": {
				Description: "Account identifier of the GitOps repository certificate.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"org_id": {
				Description: "Organization identifier of the GitOps repository certificate.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"project_id": {
				Description: "Project identifier of the GitOps repository certificate.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}

	return resource
}

func dataSourceGitopsRepoCertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	ctx = context.WithValue(ctx, nextgen.ContextAccessToken, hh.EnvVars.BearerToken.Get())

	agentIdentifier := d.Get("agent_id").(string)

	resp, httpResp, err := c.RepositoryCertificatesApi.AgentCertificateServiceList(ctx, agentIdentifier, c.AccountId, &nextgen.RepositoryCertificatesApiAgentCertificateServiceListOpts{
		OrgIdentifier:     optional.NewString(d.Get("org_id").(string)),
		ProjectIdentifier: optional.NewString(d.Get("project_id").(string)),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	setGitopsRepositoriesCertificates(d, &resp, c.AccountId, agentIdentifier)
	return nil
}
