package connector

import (
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceConnectorJDBC() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a Harness JDBC Connector.",
		ReadContext: resourceConnectorJDBCRead,

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL of the database server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"delegate_selectors": {
				Description: "Tags to filter delegates for connection.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"credentials": {
				Description: "The credentials to use for the database server.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Description: "Authentication types for JDBC connector",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username": {
							Description: "The username to use for the database server.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username_ref": {
							Description: "The reference to the Harness secret containing the username to use for the database server." + secret_ref_text,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"password_ref": {
							Description: "The reference to the Harness secret containing the password to use for the database server." + secret_ref_text,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username_password": {
							Description: "Authenticate using username password.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Description: "Username to use for authentication.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"username_ref": {
										Description: "Reference to a secret containing the username to use for authentication." + secret_ref_text,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"password_ref": {
										Description: "Reference to a secret containing the password to use for authentication." + secret_ref_text,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"service_account": {
							Description: "Authenticate using service account.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token_ref": {
										Description: "Reference to a secret containing the token to use for authentication." + secret_ref_text,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	helpers.SetMultiLevelDatasourceSchemaIdentifierRequired(resource.Schema)

	return resource
}
