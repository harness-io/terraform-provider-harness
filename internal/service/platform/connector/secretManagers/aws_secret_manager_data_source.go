package secretManagers

import (
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceConnectorAwsSM() *schema.Resource {
	resource := &schema.Resource{
		Description: "Datasource for looking up an AWS Secret Manager connector.",
		ReadContext: resourceConnectorAwsSMRead,

		Schema: map[string]*schema.Schema{
			"secret_name_prefix": {
				Description: "A prefix to be added to all secrets.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"region": {
				Description: "The AWS region where the AWS Secret Manager is.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"use_put_secret": {
				Description: "Whether to update secret value using putSecretValue action.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"force_delete_without_recovery": {
            	Description: "Whether to force delete secret value or not.",
            	Type:        schema.TypeBool,
            	Optional:    true,
            },
            "recovery_window_in_days": {
            	Description: "Recovery duration in days in AWS Secrets Manager.",
            	Type:        schema.TypeInt,
            	Optional:    true,
            },
			"delegate_selectors": {
				Description: "Tags to filter delegates for connection.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"credentials": {
				Description: "Credentials to connect to AWS.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"manual": {
							Description: "Specify the AWS key and secret used for authenticating.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_ref": {
										Description: "The reference to the Harness secret containing the AWS access key." + secret_ref_text,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"secret_key_ref": {
										Description: "The reference to the Harness secret containing the AWS secret key." + secret_ref_text,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"inherit_from_delegate": {
							Description: "Inherit the credentials from from the delegate.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"assume_role": {
							Description: "Connect using STS assume role.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_arn": {
										Description: "The ARN of the role to assume.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"external_id": {
										Description: "If the administrator of the account to which the role belongs provided you with an external ID, then enter that value.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"duration": {
										Description: "The duration, in seconds, of the role session. The value can range from 900 seconds (15 minutes) to 3600 seconds (1 hour). By default, the value is set to 3600 seconds. An expiration can also be specified in the client request body. The minimum value is 1 hour.",
										Type:        schema.TypeInt,
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
