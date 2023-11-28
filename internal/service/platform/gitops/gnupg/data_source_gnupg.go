package gnupg

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

func DataSourceGitopsGnupg() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for fetching a Harness GitOps GPG public key.",

		ReadContext: dataSourceGitopsGnupgRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account Identifier for the GnuPG Key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"org_id": {
				Description: "Organization Identifier for the GnuPG Key.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"project_id": {
				Description: "Project Identifier for the GnuPG Key.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"agent_id": {
				Description: "Agent identifier for the GnuPG Key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"identifier": {
				Description: "Identifier for the GnuPG Key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"request": {
				Description: "GnuPGPublicKey is a representation of a GnuPG public key",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upsert": {
							Description: "Indicates if the GnuPG Key should be inserted if not present or updated if present.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"publickey": {
							Description: "Public key details.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Description: "KeyID specifies the key ID, in hexadecimal string format.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"fingerprint": {
										Description: "Fingerprint is the fingerprint of the key",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"owner": {
										Description: "Owner holds the owner identification, e.g. a name and e-mail address",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"trust": {
										Description: "Trust holds the level of trust assigned to this key",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"sub_type": {
										Description: "SubType holds the key's sub type",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"key_data": {
										Description: "KeyData holds the raw key data, in base64 encoded format.",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return resource
}

func dataSourceGitopsGnupgRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	ctx = context.WithValue(ctx, nextgen.ContextAccessToken, hh.EnvVars.BearerToken.Get())
	var agentIdentifier, orgIdentifier, projectIdentifier string
	keyId := d.Get("identifier").(string)
	if attr, ok := d.GetOk("agent_id"); ok {
		agentIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}

	resp, httpResp, err := c.GnuPGPKeysApi.AgentGPGKeyServiceGet(ctx, agentIdentifier, keyId, c.AccountId, &nextgen.GnuPGPKeysApiAgentGPGKeyServiceGetOpts{
		OrgIdentifier:     optional.NewString(orgIdentifier),
		ProjectIdentifier: optional.NewString(projectIdentifier),
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
	readGnupgKey(d, &resp, c.AccountId, agentIdentifier, orgIdentifier, projectIdentifier)
	return nil
}
