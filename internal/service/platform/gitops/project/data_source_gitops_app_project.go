package project

import (
	"context"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGitOpsProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"agent_id": {
				Description: "Agent identifier of the agent where argo project resides(include scope prefix)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_id": {
				Description: "Account identifier of the GitOps Agent where argo project resides.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Deprecated:  "This field is deprecated and will be removed in a future release.",
			},
			"org_id": {
				Description: "Org identifier of the GitOps Agent where argo project resides.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"project_id": {
				Description: "Project identifier of the Gitops Agent where argo project resides.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"query_name": {
				Description: "Identifier for the GitOps Argo project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	var orgIdentifier, projectIdentifier, agentIdentifier, query_name, accountIdentifier string
	if attr, ok := d.GetOk("account_id"); ok {
		accountIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("agent_id"); ok {
		agentIdentifier = attr.(string)
	}

	if attr, ok := d.GetOk("query_name"); ok {
		query_name = attr.(string)
	}

	resp, httpResp, err := c.ProjectGitOpsApi.AgentProjectServiceGet(ctx, agentIdentifier, query_name, accountIdentifier, &nextgen.ProjectsApiAgentProjectServiceGetOpts{
		OrgIdentifier:     optional.NewString(orgIdentifier),
		ProjectIdentifier: optional.NewString(projectIdentifier),
	})

	if err != nil {
		return helpers.HandleReadApiError(err, d, httpResp)
	}
	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	if resp.Metadata == nil {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}
	setProjectDetails(d, accountIdentifier, &resp)

	return nil
}
