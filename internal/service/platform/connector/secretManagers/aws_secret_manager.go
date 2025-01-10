package secretManagers

import (
	"context"
	"fmt"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceConnectorAwsSM() *schema.Resource {
	resource := &schema.Resource{
		Description:   "Resource for creating an AWS Secret Manager connector.",
		ReadContext:   resourceConnectorAwsSMRead,
		CreateContext: resourceConnectorAwsSMCreateOrUpdate,
		UpdateContext: resourceConnectorAwsSMCreateOrUpdate,
		DeleteContext: resourceConnectorDelete,
		Importer:      helpers.MultiLevelResourceImporter,

		Schema: map[string]*schema.Schema{
			"secret_name_prefix": {
				Description: "A prefix to be added to all secrets.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"region": {
				Description: "The AWS region where the AWS Secret Manager is.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"delegate_selectors": {
				Description: "Tags to filter delegates for connection.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default": {
				Description: "Use as Default Secrets Manager.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"use_put_secret": {
				Description: "Whether to update secret value using putSecretValue action.",
				Type:        schema.TypeBool,
				Optional:    true,
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
			"credentials": {
				Description: "Credentials to connect to AWS.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"manual": {
							Description:   "Specify the AWS key and secret used for authenticating.",
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.assume_role", "credentials.0.inherit_from_delegate"},
							AtLeastOneOf:  []string{"credentials.0.manual", "credentials.0.assume_role", "credentials.0.inherit_from_delegate"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_ref": {
										Description: "The reference to the Harness secret containing the AWS access key." + secret_ref_text,
										Type:        schema.TypeString,
										Optional:    true,
									},
									"secret_key_ref": {
										Description: "The reference to the Harness secret containing the AWS secret key." + secret_ref_text,
										Type:        schema.TypeString,
										Required:    true,
									},
									"access_key_plain_text": {
										Description: "The plain text AWS access key.",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"inherit_from_delegate": {
							Description:   "Inherit the credentials from from the delegate.",
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.manual", "credentials.0.assume_role"},
							AtLeastOneOf:  []string{"credentials.0.manual", "credentials.0.assume_role", "credentials.0.inherit_from_delegate"},
							RequiredWith:  []string{"delegate_selectors"},
						},
						"assume_role": {
							Description:   "Connect using STS assume role.",
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"credentials.0.manual", "credentials.0.inherit_from_delegate"},
							AtLeastOneOf:  []string{"credentials.0.manual", "credentials.0.assume_role", "credentials.0.inherit_from_delegate"},
							RequiredWith:  []string{"delegate_selectors"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_arn": {
										Description: "The ARN of the role to assume.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"external_id": {
										Description: "If the administrator of the account to which the role belongs provided you with an external ID, then enter that value.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"duration": {
										Description: "The duration, in seconds, of the role session. The value can range from 900 seconds (15 minutes) to 3600 seconds (1 hour). By default, the value is set to 3600 seconds. An expiration can also be specified in the client request body. The minimum value is 1 hour.",
										Type:        schema.TypeInt,
										Required:    true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(int)
											if value < 900 || value > 3600 {
												errors = append(errors, fmt.Errorf("%q must be between 900 and 3600", k))
											}
											return
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourceConnectorAwsSMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := resourceConnectorReadBase(ctx, d, meta, nextgen.ConnectorTypes.AwsSecretManager)
	if err != nil {
		return err
	}

	if conn == nil {
		return nil
	}

	if err := readConnectorAwsSM(d, conn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorAwsSMCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := buildConnectorAwsSM(d)

	newConn, err := resourceConnectorCreateOrUpdateBase(ctx, d, meta, conn)
	if err != nil {
		return err
	}

	if err := readConnectorAwsSM(d, newConn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildConnectorAwsSM(d *schema.ResourceData) *nextgen.ConnectorInfo {
	connector := &nextgen.ConnectorInfo{
		Type_:            nextgen.ConnectorTypes.AwsSecretManager,
		AwsSecretManager: &nextgen.AwsSecretManager{},
	}

	if attr, ok := d.GetOk("secret_name_prefix"); ok {
		connector.AwsSecretManager.SecretNamePrefix = attr.(string)
	}

	if attr, ok := d.GetOk("region"); ok {
		connector.AwsSecretManager.Region = attr.(string)
	}

	if attr, ok := d.GetOk("delegate_selectors"); ok {
		connector.AwsSecretManager.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("default"); ok {
		connector.AwsSecretManager.Default_ = attr.(bool)
	}

	if attr, ok := d.GetOk("use_put_secret"); ok {
		connector.AwsSecretManager.UsePutSecret = attr.(bool)
	}

	if attr, ok := d.GetOk("recovery_window_in_days"); ok {
		connector.AwsSecretManager.RecoveryWindowInDays = int64(attr.(int))
	}

	if attr, ok := d.GetOk("force_delete_without_recovery"); ok {
		connector.AwsSecretManager.ForceDeleteWithoutRecovery = attr.(bool)
	}

	if attr, ok := d.GetOk("credentials"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.AwsSecretManager.Credential = &nextgen.AwsSecretManagerCredential{}

		if attr, ok := config["inherit_from_delegate"]; ok && attr.(bool) {
			connector.AwsSecretManager.Credential.Type_ = nextgen.AwsSecretManagerAuthTypes.AssumeIAMRole
			connector.AwsSecretManager.Credential.AssumeIamRole = &nextgen.AwsSmCredentialSpecAssumeIam{}
		}

		if attr := config["manual"].([]interface{}); len(attr) > 0 {
			config := attr[0].(map[string]interface{})
			connector.AwsSecretManager.Credential.Type_ = nextgen.AwsSecretManagerAuthTypes.ManualConfig
			connector.AwsSecretManager.Credential.ManualConfig = &nextgen.AwsSmCredentialSpecManualConfig{}

			if attr, ok := config["access_key_ref"]; ok {
				connector.AwsSecretManager.Credential.ManualConfig.AccessKey = attr.(string)
			}

			if attr, ok := config["secret_key_ref"]; ok {
				connector.AwsSecretManager.Credential.ManualConfig.SecretKey = attr.(string)
			}
			if attr, ok := config["access_key_plain_text"]; ok {
				connector.AwsSecretManager.Credential.ManualConfig.AccessKeyPlainText = attr.(string)
			}
		}

		if attr := config["assume_role"].([]interface{}); len(attr) > 0 {
			config := attr[0].(map[string]interface{})
			connector.AwsSecretManager.Credential.Type_ = nextgen.AwsSecretManagerAuthTypes.AssumeSTSRole
			connector.AwsSecretManager.Credential.AssumeStsRole = &nextgen.AwsSmCredentialSpecAssumeSts{}

			if attr, ok := config["role_arn"]; ok {
				connector.AwsSecretManager.Credential.AssumeStsRole.RoleArn = attr.(string)
			}

			if attr, ok := config["external_id"]; ok {
				connector.AwsSecretManager.Credential.AssumeStsRole.ExternalId = attr.(string)
			}

			if attr, ok := config["duration"]; ok {
				connector.AwsSecretManager.Credential.AssumeStsRole.AssumeStsRoleDuration = int32(attr.(int))
			}
		}
	}

	return connector
}

func readConnectorAwsSM(d *schema.ResourceData, connector *nextgen.ConnectorInfo) error {
	d.Set("secret_name_prefix", connector.AwsSecretManager.SecretNamePrefix)
	d.Set("region", connector.AwsSecretManager.Region)
	d.Set("delegate_selectors", connector.AwsSecretManager.DelegateSelectors)
	d.Set("default", connector.AwsSecretManager.Default_)
	d.Set("use_put_secret", connector.AwsSecretManager.UsePutSecret)
	d.Set("recovery_window_in_days", connector.AwsSecretManager.RecoveryWindowInDays)
	d.Set("force_delete_without_recovery", connector.AwsSecretManager.ForceDeleteWithoutRecovery)


	switch connector.AwsSecretManager.Credential.Type_ {
	case nextgen.AwsSecretManagerAuthTypes.AssumeIAMRole:
		d.Set("credentials", []interface{}{
			map[string]interface{}{
				"inherit_from_delegate": true,
			},
		})
	case nextgen.AwsSecretManagerAuthTypes.ManualConfig:
		d.Set("credentials", []interface{}{
			map[string]interface{}{
				"manual": []interface{}{
					map[string]interface{}{
						"access_key_ref":        connector.AwsSecretManager.Credential.ManualConfig.AccessKey,
						"secret_key_ref":        connector.AwsSecretManager.Credential.ManualConfig.SecretKey,
						"access_key_plain_text": connector.AwsSecretManager.Credential.ManualConfig.AccessKeyPlainText,
					},
				},
			},
		})
	case nextgen.AwsSecretManagerAuthTypes.AssumeSTSRole:
		d.Set("credentials", []interface{}{
			map[string]interface{}{
				"assume_role": []interface{}{
					map[string]interface{}{
						"role_arn":    connector.AwsSecretManager.Credential.AssumeStsRole.RoleArn,
						"external_id": connector.AwsSecretManager.Credential.AssumeStsRole.ExternalId,
						"duration":    connector.AwsSecretManager.Credential.AssumeStsRole.AssumeStsRoleDuration,
					},
				},
			},
		})
	default:
		return fmt.Errorf("unsupported aws kms auth type: %s", connector.AwsSecretManager.Credential.Type_)
	}

	return nil
}
