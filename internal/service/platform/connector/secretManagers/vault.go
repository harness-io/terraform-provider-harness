package secretManagers

import (
	"context"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceConnectorVault() *schema.Resource {
	resource := &schema.Resource{
		Description:   "Resource for creating a HashiCorp Vault Secret Manager connector.",
		ReadContext:   resourceConnectorVaultRead,
		CreateContext: resourceConnectorVaultCreateOrUpdate,
		UpdateContext: resourceConnectorVaultCreateOrUpdate,
		DeleteContext: resourceConnectorDelete,
		Importer:      helpers.MultiLevelResourceImporter,

		Schema: map[string]*schema.Schema{
			"auth_token": {
				Description: "Authentication token for Vault.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"base_path": {
				Description: "Location of the Vault directory where the secret will be stored.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vault_url": {
				Description: "URL of the HashiCorp Vault.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_read_only": {
				Description: "Read only or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"renewal_interval_minutes": {
				Description: "The time interval for the token renewal.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"secret_engine_manually_configured": {
				Description: "Manually entered Secret Engine.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"secret_engine_name": {
				Description: "Name of the Secret Engine.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"app_role_id": {
				Description: "ID of App Role.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"secret_id": {
				Description: "ID of the Secret.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_default": {
				Description: "Is default or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"secret_engine_version": {
				Description: "Version of Secret Engine.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"delegate_selectors": {
				Description: "List of Delegate Selectors that belong to the same Delegate and are used to connect to the Secret Manager.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"namespace": {
				Description: "Vault namespace where the Secret will be created.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"use_k8s_auth": {
				Description: "Boolean value to indicate if K8s Auth is used for authentication.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"use_jwt_auth": {
				Description: "Boolean value to indicate if JWT is used for authentication.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"sink_path": {
				Description: "The location from which the authentication token should be read.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"use_vault_agent": {
				Description: "Boolean value to indicate if Vault Agent is used for authentication.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"use_aws_iam": {
				Description: "Boolean value to indicate if AWS IAM is used for authentication.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"aws_region": {
				Description: "AWS region where the AWS IAM authentication will happen.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vault_aws_iam_role": {
				Description: "The Vault role defined to bind to aws iam account/role being accessed.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vault_jwt_auth_role": {
				Description: "The Vault role defined with JWT auth type for accessing Vault as per policies binded.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vault_jwt_auth_path": {
				Description: "Custom path at with JWT auth in enabled for Vault",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"xvault_aws_iam_server_id": {
				Description: "The AWS IAM Header Server ID that has been configured for this AWS IAM instance.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vault_k8s_auth_role": {
				Description: "The role where K8s Auth will happen.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"service_account_token_path": {
				Description: "The Service Account token path in the K8s pod where the token is mounted.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"k8s_auth_endpoint": {
				Description: "The path where Kubernetes Auth is enabled in Vault.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"renew_app_role_token": {
				Description: "Boolean value to indicate if AppRole token renewal is enabled or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"access_type": {
				Description:  "Access type.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"APP_ROLE", "TOKEN", "VAULT_AGENT", "AWS_IAM", "K8s_AUTH", "JWT"}, false),
			},
			"default": {
				Description: "Is default or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"read_only": {
				Description: "Read only.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"execute_on_delegate": {
				Description: "Execute on delegate or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourceConnectorVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := resourceConnectorReadBase(ctx, d, meta, nextgen.ConnectorTypes.Vault)
	if err != nil {
		return err
	}

	if conn == nil {
		return nil
	}

	if err := readConnectorVault(d, conn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorVaultCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := buildConnectorVault(d)

	newConn, err := resourceConnectorCreateOrUpdateBase(ctx, d, meta, conn)
	if err != nil {
		return err
	}

	if err := readConnectorVault(d, newConn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildConnectorVault(d *schema.ResourceData) *nextgen.ConnectorInfo {
	connector := &nextgen.ConnectorInfo{
		Type_: nextgen.ConnectorTypes.Vault,
		Vault: &nextgen.VaultConnector{},
	}

	if attr, ok := d.GetOk("auth_token"); ok {
		connector.Vault.AuthToken = attr.(string)
	}

	if attr, ok := d.GetOk("base_path"); ok {
		connector.Vault.BasePath = attr.(string)
	}

	if attr, ok := d.GetOk("vault_url"); ok {
		connector.Vault.VaultUrl = attr.(string)
	}

	if attr, ok := d.GetOk("is_read_only"); ok {
		connector.Vault.IsReadOnly = attr.(bool)
	}

	if attr, ok := d.GetOk("renewal_interval_minutes"); ok {
		connector.Vault.RenewalIntervalMinutes = int64(attr.(int))
	}

	if attr, ok := d.GetOk("secret_engine_manually_configured"); ok {
		connector.Vault.SecretEngineManuallyConfigured = attr.(bool)
	}

	if attr, ok := d.GetOk("secret_engine_name"); ok {
		connector.Vault.SecretEngineName = attr.(string)
	}

	if attr, ok := d.GetOk("app_role_id"); ok {
		connector.Vault.AppRoleId = attr.(string)
	}

	if attr, ok := d.GetOk("secret_id"); ok {
		connector.Vault.SecretId = attr.(string)
	}

	if attr, ok := d.GetOk("is_default"); ok {
		connector.Vault.IsDefault = attr.(bool)
	}

	if attr, ok := d.GetOk("secret_engine_version"); ok {
		connector.Vault.SecretEngineVersion = int32(attr.(int))
	}

	if attr, ok := d.GetOk("delegate_selectors"); ok {
		connector.Vault.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("namespace"); ok {
		connector.Vault.Namespace = attr.(string)
	}

	if attr, ok := d.GetOk("sink_path"); ok {
		connector.Vault.SinkPath = attr.(string)
	}

	if attr, ok := d.GetOk("use_vault_agent"); ok {
		connector.Vault.UseVaultAgent = attr.(bool)
	}

	if attr, ok := d.GetOk("use_aws_iam"); ok {
		connector.Vault.UseAwsIam = attr.(bool)
	}

	if attr, ok := d.GetOk("aws_region"); ok {
		connector.Vault.AwsRegion = attr.(string)
	}

	if attr, ok := d.GetOk("vault_aws_iam_role"); ok {
		connector.Vault.VaultAwsIamRole = attr.(string)
	}

	if attr, ok := d.GetOk("xvault_aws_iam_server_id"); ok {
		connector.Vault.XvaultAwsIamServerId = attr.(string)
	}

	if attr, ok := d.GetOk("use_k8s_auth"); ok {
		connector.Vault.UseK8sAuth = attr.(bool)
	}

	if attr, ok := d.GetOk("use_jwt_auth"); ok {
		connector.Vault.UseJwtAuth = attr.(bool)
	}

	if attr, ok := d.GetOk("vault_jwt_auth_path"); ok {
		connector.Vault.JwtAuthPath = attr.(string)
	}

	if attr, ok := d.GetOk("vault_jwt_auth_role"); ok {
		connector.Vault.JwtAuthRole = attr.(string)
	}

	if attr, ok := d.GetOk("vault_k8s_auth_role"); ok {
		connector.Vault.VaultK8sAuthRole = attr.(string)
	}

	if attr, ok := d.GetOk("service_account_token_path"); ok {
		connector.Vault.ServiceAccountTokenPath = attr.(string)
	}

	if attr, ok := d.GetOk("k8s_auth_endpoint"); ok {
		connector.Vault.K8sAuthEndpoint = attr.(string)
	}

	if attr, ok := d.GetOk("renew_app_role_token"); ok {
		connector.Vault.RenewAppRoleToken = attr.(bool)
	}

	if attr, ok := d.GetOk("access_type"); ok {
		connector.Vault.AccessType = attr.(string)
	}

	if attr, ok := d.GetOk("default"); ok {
		connector.Vault.Default_ = attr.(bool)
	}

	if attr, ok := d.GetOk("read_only"); ok {
		connector.Vault.ReadOnly = attr.(bool)
	}

	if attr, ok := d.GetOk("execute_on_delegate"); ok {
		connector.Vault.ExecuteOnDelegate = attr.(bool)
	}

	return connector
}

func readConnectorVault(d *schema.ResourceData, connector *nextgen.ConnectorInfo) error {
	d.Set("auth_token", connector.Vault.AuthToken)
	d.Set("base_path", connector.Vault.BasePath)
	d.Set("vault_url", connector.Vault.VaultUrl)
	d.Set("is_read_only", connector.Vault.IsReadOnly)
	d.Set("renewal_interval_minutes", connector.Vault.RenewalIntervalMinutes)
	d.Set("secret_engine_manually_configured", connector.Vault.SecretEngineManuallyConfigured)
	d.Set("secret_engine_name", connector.Vault.SecretEngineName)
	d.Set("app_role_id", connector.Vault.AppRoleId)
	d.Set("secret_id", connector.Vault.SecretId)
	d.Set("is_default", connector.Vault.IsDefault)
	d.Set("secret_engine_version", connector.Vault.SecretEngineVersion)
	d.Set("delegate_selectors", connector.Vault.DelegateSelectors)
	d.Set("namespace", connector.Vault.Namespace)
	d.Set("sink_path", connector.Vault.SinkPath)
	d.Set("use_vault_agent", connector.Vault.UseVaultAgent)
	d.Set("use_aws_iam", connector.Vault.UseAwsIam)
	d.Set("aws_region", connector.Vault.AwsRegion)
	d.Set("vault_aws_iam_role", connector.Vault.VaultAwsIamRole)
	d.Set("xvault_aws_iam_server_id", connector.Vault.XvaultAwsIamServerId)
	d.Set("use_k8s_auth", connector.Vault.UseK8sAuth)
	d.Set("use_jwt_auth", connector.Vault.UseJwtAuth)
	d.Set("execute_on_delegate", connector.Vault.ExecuteOnDelegate)
	d.Set("vault_k8s_auth_role", connector.Vault.VaultK8sAuthRole)
	d.Set("service_account_token_path", connector.Vault.ServiceAccountTokenPath)
	d.Set("k8s_auth_endpoint", connector.Vault.K8sAuthEndpoint)
	d.Set("renew_app_role_token", connector.Vault.RenewAppRoleToken)
	d.Set("access_type", connector.Vault.AccessType)
	d.Set("default", connector.Vault.Default_)
	d.Set("read_only", connector.Vault.ReadOnly)

	return nil
}
