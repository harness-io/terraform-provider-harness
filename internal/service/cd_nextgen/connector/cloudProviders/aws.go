package cloudProviders

import (
	"context"
	"fmt"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceConnectorAws() *schema.Resource {
	resource := &schema.Resource{
		Description:   "Resource for creating an AWS connector.",
		ReadContext:   resourceConnectorAwsRead,
		CreateContext: resourceConnectorAwsCreateOrUpdate,
		UpdateContext: resourceConnectorAwsCreateOrUpdate,
		DeleteContext: resourceConnectorDelete,
		Importer:      helpers.MultiLevelResourceImporter,

		Schema: map[string]*schema.Schema{
			"manual": {
				Description: "Use IAM role for service accounts.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"irsa",
					"inherit_from_delegate",
					"oidc_authentication",
				},
				ExactlyOneOf: []string{
					"manual",
					"irsa",
					"inherit_from_delegate",
					"oidc_authentication",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key": {
							Description:   "AWS access key.",
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"manual.0.access_key_ref"},
							AtLeastOneOf:  []string{"manual.0.access_key", "manual.0.access_key_ref"},
						},
						"access_key_ref": {
							Description:   "Reference to the Harness secret containing the aws access key." + secret_ref_text,
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"manual.0.access_key"},
							AtLeastOneOf:  []string{"manual.0.access_key", "manual.0.access_key_ref"},
						},
						"secret_key_ref": {
							Description: "Reference to the Harness secret containing the aws secret key." + secret_ref_text,
							Type:        schema.TypeString,
							Required:    true,
						},
						"session_token_ref": {
							Description: "Reference to the Harness secret containing the aws session token." + secret_ref_text,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"delegate_selectors": {
							Description: "Connect only use delegates with these tags.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"execute_on_delegate": {
							Description: "Execute on delegate or not.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"region": {
							Description: "Test Region to perform Connection test of AWS Connector" + secret_ref_text,
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"irsa": {
				Description: "Use IAM role for service accounts.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"manual",
					"inherit_from_delegate",
					"oidc_authentication",
				},
				ExactlyOneOf: []string{
					"manual",
					"irsa",
					"inherit_from_delegate",
					"oidc_authentication",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delegate_selectors": {
							Description: "The delegates to inherit the credentials from.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"region": {
							Description: "Test Region to perform Connection test of AWS Connector" + secret_ref_text,
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"inherit_from_delegate": {
				Description: "Inherit credentials from the delegate.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"irsa",
					"manual",
					"oidc_authentication",
				},
				ExactlyOneOf: []string{
					"manual",
					"irsa",
					"inherit_from_delegate",
					"oidc_authentication",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delegate_selectors": {
							Description: "The delegates to inherit the credentials from.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"region": {
							Description: "Test Region to perform Connection test of AWS Connector" + secret_ref_text,
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"oidc_authentication": {
				Description: "Authentication using harness oidc.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"irsa",
					"manual",
					"inherit_from_delegate",
				},
				ExactlyOneOf: []string{
					"manual",
					"irsa",
					"inherit_from_delegate",
					"oidc_authentication",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_role_arn": {
							Description: "The IAM Role to assume the credentials from.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"delegate_selectors": {
							Description: "The delegates to inherit the credentials from.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"execute_on_delegate": {
							Description: "Execute on delegate or not.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"region": {
							Description: "Test Region to perform Connection test of AWS Connector." + secret_ref_text,
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"cross_account_access": {
				Description: "Select this option if you want to use one AWS account for the connection, but you want to deploy or build in a different AWS account. In this scenario, the AWS account used for AWS access in Credentials will assume the IAM role you specify in Cross-account role ARN setting. This option uses the AWS Security Token Service (STS) feature.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Description: "The Amazon Resource Name (ARN) of the role that you want to assume. This is an IAM role in the target AWS account.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"external_id": {
							Description: "If the administrator of the account to which the role belongs provided you with an external ID, then enter that value.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"equal_jitter_backoff_strategy": {
				Description: "Equal Jitter BackOff Strategy.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"full_jitter_backoff_strategy",
					"fixed_delay_backoff_strategy",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_delay": {
							Description: "Base delay.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"max_backoff_time": {
							Description: "Max BackOff Time.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"retry_count": {
							Description: "Retry Count.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"full_jitter_backoff_strategy": {
				Description: "Full Jitter BackOff Strategy.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"equal_jitter_backoff_strategy",
					"fixed_delay_backoff_strategy",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_delay": {
							Description: "Base delay.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"max_backoff_time": {
							Description: "Max BackOff Time.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"retry_count": {
							Description: "Retry Count.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"fixed_delay_backoff_strategy": {
				Description: "Fixed Delay BackOff Strategy.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"full_jitter_backoff_strategy",
					"equal_jitter_backoff_strategy",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fixed_backoff": {
							Description: "Fixed Backoff.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"retry_count": {
							Description: "Retry Count.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"force_delete": {
				Description: "Enable this flag for force deletion of connector",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}

	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourceConnectorAwsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := resourceConnectorReadBase(ctx, d, meta, nextgen.ConnectorTypes.Aws)
	if err != nil {
		return err
	}

	if conn == nil {
		return nil
	}

	if err := readConnectorAws(d, conn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorAwsCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := buildConnectorAws(d)

	newConn, err := resourceConnectorCreateOrUpdateBase(ctx, d, meta, conn)
	if err != nil {
		return err
	}

	if err := readConnectorAws(d, newConn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildConnectorAws(d *schema.ResourceData) *nextgen.ConnectorInfo {
	connector := &nextgen.ConnectorInfo{
		Type_: nextgen.ConnectorTypes.Aws,
		Aws: &nextgen.AwsConnector{
			Credential: &nextgen.AwsCredential{},
		},
	}

	if attr, ok := d.GetOk("manual"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.Credential.Type_ = nextgen.AwsAuthTypes.ManualConfig
		connector.Aws.Credential.ManualConfig = &nextgen.AwsManualConfigSpec{}

		if attr := config["access_key"].(string); attr != "" {
			connector.Aws.Credential.ManualConfig.AccessKey = attr
		}

		if attr := config["access_key_ref"].(string); attr != "" {
			connector.Aws.Credential.ManualConfig.AccessKeyRef = attr
		}

		if attr := config["secret_key_ref"].(string); attr != "" {
			connector.Aws.Credential.ManualConfig.SecretKeyRef = attr
		}

		if attr := config["session_token_ref"].(string); attr != "" {
			connector.Aws.Credential.ManualConfig.SessionTokenRef = attr
		}

		if attr := config["delegate_selectors"].(*schema.Set).List(); len(attr) > 0 {
			connector.Aws.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr)
		}

		if attr, ok := d.GetOk("execute_on_delegate"); ok {
			connector.Aws.ExecuteOnDelegate = attr.(bool)
		}

		if attr := config["region"].(string); attr != "" {
			connector.Aws.Credential.Region = attr
		}
	}

	if attr, ok := d.GetOk("irsa"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.Credential.Type_ = nextgen.AwsAuthTypes.Irsa

		if attr := config["delegate_selectors"].(*schema.Set).List(); len(attr) > 0 {
			connector.Aws.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr)
		}

		if attr := config["region"].(string); attr != "" {
			connector.Aws.Credential.Region = attr
		}
	}

	if attr, ok := d.GetOk("inherit_from_delegate"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.Credential.Type_ = nextgen.AwsAuthTypes.InheritFromDelegate

		if attr := config["delegate_selectors"].(*schema.Set).List(); len(attr) > 0 {
			connector.Aws.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr)
		}

		if attr := config["region"].(string); attr != "" {
			connector.Aws.Credential.Region = attr
		}
	}

	if attr, ok := d.GetOk("oidc_authentication"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.Credential.Type_ = nextgen.AwsAuthTypes.OidcAuthentication
		connector.Aws.Credential.OidcConfig = &nextgen.AwsOidcConfigSpec{}

		if attr := config["iam_role_arn"].(string); attr != "" {
			connector.Aws.Credential.OidcConfig.IamRoleArn = attr
		}

		if attr := config["delegate_selectors"].(*schema.Set).List(); len(attr) > 0 {
			connector.Aws.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr)
		}

		if attr, ok := d.GetOk("execute_on_delegate"); ok {
			connector.Aws.ExecuteOnDelegate = attr.(bool)
		}

		if attr := config["region"].(string); attr != "" {
			connector.Aws.Credential.Region = attr
		}
	}

	if attr, ok := d.GetOk("cross_account_access"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.Credential.CrossAccountAccess = &nextgen.CrossAccountAccess{}

		if attr := config["role_arn"].(string); attr != "" {
			connector.Aws.Credential.CrossAccountAccess.CrossAccountRoleArn = attr
		}

		if attr := config["external_id"].(string); attr != "" {
			connector.Aws.Credential.CrossAccountAccess.ExternalId = attr
		}
	}

	if attr, ok := d.GetOk("equal_jitter_backoff_strategy"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.AwsSdkClientBackOffStrategyOverride = &nextgen.AwsSdkClientBackoffStrategy{}
		connector.Aws.AwsSdkClientBackOffStrategyOverride.Type_ = nextgen.AwsSdkClientBackOffStrategyTypes.EqualJitterBackoffStrategy
		connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter = &nextgen.AwsEqualJitterBackoffStrategy{}

		if val, ok := config["retry_count"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter.RetryCount = int32(val.(int))
		}
		if val, ok := config["max_backoff_time"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter.MaxBackoffTime = int64(val.(int))
		}
		if val, ok := config["base_delay"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter.BaseDelay = int64(val.(int))
		}
	}

	if attr, okk := d.GetOk("full_jitter_backoff_strategy"); okk {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.AwsSdkClientBackOffStrategyOverride = &nextgen.AwsSdkClientBackoffStrategy{}
		connector.Aws.AwsSdkClientBackOffStrategyOverride.Type_ = nextgen.AwsSdkClientBackOffStrategyTypes.FullJitterBackoffStrategy
		connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter = &nextgen.AwsFullJitterBackoffStrategy{}

		if val, ok := config["retry_count"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter.RetryCount = int32(val.(int))
		}
		if val, ok := config["max_backoff_time"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter.MaxBackoffTime = int64(val.(int))
		}
		if val, ok := config["base_delay"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter.BaseDelay = int64(val.(int))
		}
	}

	if attr, okk := d.GetOk("fixed_delay_backoff_strategy"); okk {
		config := attr.([]interface{})[0].(map[string]interface{})
		connector.Aws.AwsSdkClientBackOffStrategyOverride = &nextgen.AwsSdkClientBackoffStrategy{}
		connector.Aws.AwsSdkClientBackOffStrategyOverride.Type_ = nextgen.AwsSdkClientBackOffStrategyTypes.FixedDelayBackoffStrategy
		connector.Aws.AwsSdkClientBackOffStrategyOverride.FixedDelay = &nextgen.AwsFixedDelayBackoffStrategy{}

		if val, ok := config["retry_count"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.FixedDelay.RetryCount = int32(val.(int))
		}
		if val, ok := config["fixed_backoff"]; ok {
			connector.Aws.AwsSdkClientBackOffStrategyOverride.FixedDelay.FixedBackoff = int64(val.(int))
		}
	}

	return connector
}

func readConnectorAws(d *schema.ResourceData, connector *nextgen.ConnectorInfo) error {

	d.Set("execute_on_delegate", connector.Aws.ExecuteOnDelegate)
	switch connector.Aws.Credential.Type_ {
	case nextgen.AwsAuthTypes.ManualConfig:
		d.Set("manual", []map[string]interface{}{
			{
				"access_key":          connector.Aws.Credential.ManualConfig.AccessKey,
				"access_key_ref":      connector.Aws.Credential.ManualConfig.AccessKeyRef,
				"secret_key_ref":      connector.Aws.Credential.ManualConfig.SecretKeyRef,
				"session_token_ref":   connector.Aws.Credential.ManualConfig.SessionTokenRef,
				"delegate_selectors":  connector.Aws.DelegateSelectors,
				"execute_on_delegate": connector.Aws.ExecuteOnDelegate,
				"region":              connector.Aws.Credential.Region,
			},
		})
	case nextgen.AwsAuthTypes.Irsa:
		d.Set("irsa", []map[string]interface{}{
			{
				"delegate_selectors": connector.Aws.DelegateSelectors,
				"region":             connector.Aws.Credential.Region,
			},
		})
	case nextgen.AwsAuthTypes.InheritFromDelegate:
		d.Set("inherit_from_delegate", []map[string]interface{}{
			{
				"delegate_selectors": connector.Aws.DelegateSelectors,
				"region":             connector.Aws.Credential.Region,
			},
		})
	case nextgen.AwsAuthTypes.OidcAuthentication:
		d.Set("oidc_authentication", []map[string]interface{}{
			{
				"iam_role_arn":        connector.Aws.Credential.OidcConfig.IamRoleArn,
				"delegate_selectors":  connector.Aws.DelegateSelectors,
				"execute_on_delegate": connector.Aws.ExecuteOnDelegate,
				"region":              connector.Aws.Credential.Region,
			},
		})
	default:
		return fmt.Errorf("unsupported aws credential type: %s", connector.Aws.Credential.Type_)
	}
	if connector.Aws.Credential.CrossAccountAccess != nil {
		d.Set("cross_account_access", []map[string]interface{}{
			{
				"role_arn":    connector.Aws.Credential.CrossAccountAccess.CrossAccountRoleArn,
				"external_id": connector.Aws.Credential.CrossAccountAccess.ExternalId,
			},
		})
	}
	if connector.Aws.AwsSdkClientBackOffStrategyOverride != nil {
		switch connector.Aws.AwsSdkClientBackOffStrategyOverride.Type_ {
		case nextgen.AwsSdkClientBackOffStrategyTypes.EqualJitterBackoffStrategy:
			d.Set("equal_jitter_backoff_strategy", []map[string]interface{}{
				{
					"base_delay":       connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter.BaseDelay,
					"max_backoff_time": connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter.MaxBackoffTime,
					"retry_count":      connector.Aws.AwsSdkClientBackOffStrategyOverride.EqualJitter.RetryCount,
				},
			})
		case nextgen.AwsSdkClientBackOffStrategyTypes.FullJitterBackoffStrategy:
			d.Set("full_jitter_backoff_strategy", []map[string]interface{}{
				{
					"base_delay":       connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter.BaseDelay,
					"max_backoff_time": connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter.MaxBackoffTime,
					"retry_count":      connector.Aws.AwsSdkClientBackOffStrategyOverride.FullJitter.RetryCount,
				},
			})
		case nextgen.AwsSdkClientBackOffStrategyTypes.FixedDelayBackoffStrategy:
			d.Set("fixed_delay_backoff_strategy", []map[string]interface{}{
				{
					"fixed_backoff": connector.Aws.AwsSdkClientBackOffStrategyOverride.FixedDelay.FixedBackoff,
					"retry_count":   connector.Aws.AwsSdkClientBackOffStrategyOverride.FixedDelay.RetryCount,
				},
			})
		default:
			return fmt.Errorf("unsupported aws credential type: %s", connector.Aws.AwsSdkClientBackOffStrategyOverride.Type_)
		}

	}
	return nil
}
