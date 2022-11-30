package repository_credentials

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

func DataSourceGitOpsRepoCred() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a GitOps RepoCred.",

		ReadContext: dataSourceGitopsRepoCredRead,

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Description: "agent identifier of the Repository Credentials.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_id": {
				Description: "account identifier of the Repository Credentials.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"identifier": {
				Description: "Identifier of the Repository Credentials.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"org_id": {
				Description: "Organization identifier of the Repository Credential.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"project_id": {
				Description: "Project identifier of the Repository Credential.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"upsert": {
				Description: "if the Repository credential should be upserted.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"creds": {
				Description: "credential details.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Description: "url representing this object.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"username": {
							Description: "Username for authenticating at the repo server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"password": {
							Description: "Password for authenticating at the repo server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"ssh_private_key": {
							Description: "Contains the private key data for authenticating at the repo server using SSH (only Git repos).",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"tls_client_cert_data": {
							Description: "Specifies the TLS client cert data for authenticating at the repo server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"tls_client_cert_key": {
							Description: "Specifies the TLS client cert key for authenticating at the repo server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_private_key": {
							Description: "github_app_private_key specifies the private key PEM data for authentication via GitHub app.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_id": {
							Description: "Specifies the Github App ID of the app used to access the repo for GitHub app authentication.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_installation_id": {
							Description: "Specifies the ID of the installed GitHub App for GitHub app authentication.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_enterprise_base_url": {
							Description: "Specifies the GitHub API URL for GitHub app authentication.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"enable_oci": {
							Description: "Specifies whether helm-oci support should be enabled for this repo.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"type": {
							Description: "Type specifies the type of the repoCreds.Can be either 'git' or 'helm. 'git' is assumed if empty or absent",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}

	return resource
}

func dataSourceGitopsRepoCredRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	ctx = context.WithValue(ctx, nextgen.ContextAccessToken, hh.EnvVars.BearerToken.Get())

	agentIdentifier := d.Get("agent_id").(string)
	identifier := d.Get("identifier").(string)

	resp, httpResp, err := c.RepositoryCredentialsApi.AgentRepositoryCredentialsServiceGetRepositoryCredentials(ctx, agentIdentifier, identifier, c.AccountId, &nextgen.RepositoryCredentialsApiAgentRepositoryCredentialsServiceGetRepositoryCredentialsOpts{
		OrgIdentifier:     optional.NewString(d.Get("org_id").(string)),
		ProjectIdentifier: optional.NewString(d.Get("project_id").(string)),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	if &resp == nil {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}
	setGitopsRepositoriesCredential(d, &resp)
	return nil
}
