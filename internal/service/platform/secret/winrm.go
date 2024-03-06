package secret

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecretWinRM() *schema.Resource {
	resource := &schema.Resource{
		Description:   "Resource for creating a WinRM secret.",
		ReadContext:   resourceSecretWinRMRead,
		CreateContext: resourceSecretWinRMCreateOrUpdate,
		UpdateContext: resourceSecretWinRMCreateOrUpdate,
		DeleteContext: resourceSecretDelete,
		Importer:      helpers.MultiLevelResourceImporter,

		Schema: map[string]*schema.Schema{
			"port": {
				Description: "WinRM port",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5985,
			},
			"kerberos": {
				Description: "Kerberos authentication configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal": {
							Description: "Username to use for Kerberos authentication.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"realm": {
							Description: "Kerberos realm to use for authentication.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"tgt_generation_method": {
							Description:  "Method for generating TGT (Ticket Granting Ticket)",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "auto",
							ValidateFunc: validation.StringInSlice([]string{"auto", "manual"}, false),
						},
						"tgt_key_tab_file_path_spec": {
							Description: "Keytab file path used for Kerberos TGT generation",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_path": {
										Description: "Path to the keytab file.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
						"tgt_password_spec": {
							Description: "Authenticate to App Dynamics using username and password.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Description: "password. To reference a password at the organization scope, prefix 'org' to the expression: org.{identifier}. To reference a password at the account scope, prefix 'account` to the expression: account.{identifier}",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"ntlm": {
				Description: "NTLM authentication configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Description: "NTLM authentication type",
							Type:        schema.TypeString,
							Required:    true,
						},
						"domain": {
							Description: "Domain for NTLM authentication.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"username": {
							Description: "Username for NTLM authentication.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"password": {
							Description: "Password for NTLM authentication.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"use_ssl_cert": {
							Description: "Whether to use SSL certificate for NTLM authentication.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
		},
	}
	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourceSecretWinRMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secret, err := resourceSecretReadBase(ctx, d, meta, nextgen.SecretTypes.WinRmCredentials)
	if err != nil {
		return err
	}

	if err := readSecretWinRM(d, secret); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretWinRMCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secret := buildSecretWinRM(d)

	newSecret, err := resourceSecretCreateOrUpdateBase(ctx, d, meta, secret)
	if err != nil {
		return err
	}

	if err := readSecretWinRM(d, newSecret); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildSecretWinRM(d *schema.ResourceData) *nextgen.Secret {
	secret := &nextgen.Secret{
		Type_:            nextgen.SecretTypes.WinRmCredentials,
		WinRmCredentials: &nextgen.WinRmCredentialsSpec{},
	}

	if attr, ok := d.GetOk("port"); ok {
		secret.WinRmCredentials.Port = int32(attr.(int))
	}

	if attr, ok := d.GetOk("kerberos"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})

		secret.SSHKey.Auth = &nextgen.SshAuth{
			Type_:          nextgen.SSHAuthenticationTypes.Kerberos,
			KerberosConfig: &nextgen.KerberosConfig{},
		}

		if attr, ok := config["principal"]; ok {
			secret.SSHKey.Auth.KerberosConfig.Principal = attr.(string)
		}

		if attr, ok := config["realm"]; ok {
			secret.SSHKey.Auth.KerberosConfig.Realm = attr.(string)
		}

		if _, ok := d.GetOk("kerberos.0.tgt_key_tab_file_path_spec"); ok {
			if attr, ok := config["tgt_key_tab_file_path_spec"]; ok {
				config = attr.([]interface{})[0].(map[string]interface{})
				secret.SSHKey.Auth.KerberosConfig.KeyTabFilePathSpec = &nextgen.TgtKeyTabFilePathSpecDto{}

				secret.SSHKey.Auth.KerberosConfig.TgtGenerationMethod = nextgen.TgtGenerationMethodTypes.TGTKeyTabFilePathSpecDTO

				if attr, ok := config["key_path"]; ok {
					secret.SSHKey.Auth.KerberosConfig.KeyTabFilePathSpec.KeyPath = attr.(string)
				}
			}
		}

		if _, ok := d.GetOk("kerberos.0.tgt_password_spec"); ok {
			if attr, ok := config["tgt_password_spec"]; ok {
				config = attr.([]interface{})[0].(map[string]interface{})
				secret.SSHKey.Auth.KerberosConfig.PasswordSpec = &nextgen.TgtPasswordSpecDto{}

				secret.SSHKey.Auth.KerberosConfig.TgtGenerationMethod = nextgen.TgtGenerationMethodTypes.TGTPasswordSpecDTO

				if attr, ok := config["password"]; ok {
					secret.SSHKey.Auth.KerberosConfig.PasswordSpec.Password = attr.(string)
				}
			}
		}
	}

	if attr, ok := d.GetOk("ssh"); ok {
		config := attr.([]interface{})[0].(map[string]interface{})

		secret.SSHKey.Auth = &nextgen.SshAuth{
			Type_:     nextgen.SSHAuthenticationTypes.SSH,
			SSHConfig: &nextgen.SshConfig{},
		}

		if _, ok := d.GetOk("ssh.0.sshkey_path_credential"); ok {
			if attr, ok := config["sshkey_path_credential"]; ok {
				config = attr.([]interface{})[0].(map[string]interface{})
				secret.SSHKey.Auth.SSHConfig.Type_ = nextgen.SSHConfigTypes.KeyPath
				secret.SSHKey.Auth.SSHConfig.CredentialType = nextgen.SSHConfigTypes.KeyPath

				secret.SSHKey.Auth.SSHConfig.KeyPathCredential = &nextgen.SshKeyPathCredential{}
				if attr, ok := config["user_name"]; ok {
					secret.SSHKey.Auth.SSHConfig.KeyPathCredential.UserName = attr.(string)
				}

				if attr, ok := config["key_path"]; ok {
					secret.SSHKey.Auth.SSHConfig.KeyPathCredential.KeyPath = attr.(string)
				}

				if attr, ok := config["encrypted_passphrase"]; ok {
					secret.SSHKey.Auth.SSHConfig.KeyPathCredential.EncryptedPassphrase = attr.(string)
				}
			}
		}

		if _, ok := d.GetOk("ssh.0.sshkey_reference_credential"); ok {
			if attr, ok := config["sshkey_reference_credential"]; ok {
				config = attr.([]interface{})[0].(map[string]interface{})
				secret.SSHKey.Auth.SSHConfig.Type_ = nextgen.SSHConfigTypes.KeyReference
				secret.SSHKey.Auth.SSHConfig.CredentialType = nextgen.SSHConfigTypes.KeyReference

				secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential = &nextgen.SshKeyReferenceCredentialDto{}
				if attr, ok := config["user_name"]; ok {
					secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.UserName = attr.(string)
				}

				if attr, ok := config["key"]; ok {
					secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.Key = attr.(string)
				}

				if attr, ok := config["encrypted_passphrase"]; ok {
					secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.EncryptedPassphrase = attr.(string)
				}
			}
		}

		if _, ok := d.GetOk("ssh.0.ssh_password_credential"); ok {
			if attr, ok := config["ssh_password_credential"]; ok {
				config = attr.([]interface{})[0].(map[string]interface{})

				secret.SSHKey.Auth.SSHConfig.Type_ = nextgen.SSHConfigTypes.Password
				secret.SSHKey.Auth.SSHConfig.CredentialType = nextgen.SSHConfigTypes.Password
				secret.SSHKey.Auth.SSHConfig.PasswordCredential = &nextgen.SshPasswordCredentialDto{}
				if attr, ok := config["user_name"]; ok {
					secret.SSHKey.Auth.SSHConfig.PasswordCredential.UserName = attr.(string)
				}

				if attr, ok := config["password"]; ok {
					secret.SSHKey.Auth.SSHConfig.PasswordCredential.Password = attr.(string)
				}
			}
		}
	}

	return secret
}

func readSecretWinRM(d *schema.ResourceData, secret *nextgen.Secret) error {
	switch secret.WinRmCredentials.Auth.Type_ {
	case "SSH":
		d.Set("ssh", []map[string]interface{}{
			{
				"credential_type":             secret.SSHKey.Auth.SSHConfig.CredentialType,
				"sshkey_path_credential":      readWinRMPathCredentialSpec(secret),
				"sshkey_reference_credential": readWinRMReferenceCredentialSpec(secret),
				"ssh_password_credential":     readSshPasswordSpecWinRM(secret),
			},
		})

	case "Kerberos":
		d.Set("kerberos", []map[string]interface{}{
			{
				"principal":                  secret.SSHKey.Auth.KerberosConfig.Principal,
				"realm":                      secret.SSHKey.Auth.KerberosConfig.Realm,
				"tgt_generation_method":      secret.SSHKey.Auth.KerberosConfig.TgtGenerationMethod,
				"tgt_key_tab_file_path_spec": readTgtKeyFilePathSpecWinRM(secret),
				"tgt_password_spec":          readTgtPasswordSpecWinRM(secret),
			},
		})

	}
	return nil
}

func readWinRMReferenceCredentialSpec(secret *nextgen.Secret) []map[string]interface{} {
	var spec []map[string]interface{}
	switch secret.SSHKey.Auth.SSHConfig.CredentialType {
	case nextgen.SSHConfigTypes.KeyPath:
		//noop
	case nextgen.SSHConfigTypes.KeyReference:
		spec = []map[string]interface{}{
			{"user_name": secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.UserName,
				"key":                  secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.Key,
				"encrypted_passphrase": secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.EncryptedPassphrase,
			},
		}
	case nextgen.SSHConfigTypes.Password:
		//noop
	}
	return spec
}

func readWinRMPathCredentialSpec(secret *nextgen.Secret) []map[string]interface{} {
	var spec []map[string]interface{}
	switch secret.SSHKey.Auth.SSHConfig.CredentialType {
	case nextgen.SSHConfigTypes.KeyPath:
		spec = []map[string]interface{}{
			{"user_name": secret.SSHKey.Auth.SSHConfig.KeyPathCredential.UserName,
				"key_path":             secret.SSHKey.Auth.SSHConfig.KeyPathCredential.KeyPath,
				"encrypted_passphrase": secret.SSHKey.Auth.SSHConfig.KeyPathCredential.EncryptedPassphrase,
			},
		}
	case nextgen.SSHConfigTypes.KeyReference:
		//noop
	case nextgen.SSHConfigTypes.Password:
		//noop
	}
	return spec
}

func readTgtPasswordSpecWinRM(secret *nextgen.Secret) []map[string]interface{} {
	var spec []map[string]interface{}
	switch secret.SSHKey.Auth.KerberosConfig.TgtGenerationMethod {
	case nextgen.TgtGenerationMethodTypes.TGTKeyTabFilePathSpecDTO:
		//noop
	case nextgen.TgtGenerationMethodTypes.TGTPasswordSpecDTO:
		spec = []map[string]interface{}{
			{
				"password": secret.SSHKey.Auth.KerberosConfig.PasswordSpec.Password},
		}
	}
	return spec
}

func readSshPasswordSpecWinRM(secret *nextgen.Secret) []map[string]interface{} {
	var spec []map[string]interface{}
	switch secret.SSHKey.Auth.SSHConfig.CredentialType {
	case nextgen.SSHConfigTypes.Password:
		spec = []map[string]interface{}{
			{"user_name": secret.SSHKey.Auth.SSHConfig.PasswordCredential.UserName,
				"password": secret.SSHKey.Auth.SSHConfig.PasswordCredential.Password,
			},
		}
	case nextgen.SSHConfigTypes.KeyReference:
		//noop
	case nextgen.SSHConfigTypes.KeyPath:
		//noop
	}
	return spec
}

func readTgtKeyFilePathSpecWinRM(secret *nextgen.Secret) []map[string]interface{} {
	var spec []map[string]interface{}
	switch secret.SSHKey.Auth.KerberosConfig.TgtGenerationMethod {
	case nextgen.TgtGenerationMethodTypes.TGTKeyTabFilePathSpecDTO:
		spec = []map[string]interface{}{{"key_path": secret.SSHKey.Auth.KerberosConfig.KeyTabFilePathSpec.KeyPath}}
	case nextgen.TgtGenerationMethodTypes.TGTPasswordSpecDTO:
		//noop
	}
	return spec
}
