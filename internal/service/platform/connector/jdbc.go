package connector

import (
	"context"
	"fmt"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceConnectorJDBC() *schema.Resource {
	resource := &schema.Resource{
		Description:   "Resource for creating a JDBC connector.",
		ReadContext:   resourceConnectorJDBCRead,
		CreateContext: resourceConnectorJDBCCreateOrUpdate,
		UpdateContext: resourceConnectorJDBCCreateOrUpdate,
		DeleteContext: resourceConnectorDelete,
		Importer:      helpers.MultiLevelResourceImporter,

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL of the database server.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"delegate_selectors": {
				Description: "Tags to filter delegates for connection.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"credentials": {
				Description: "The credentials to use for the database server.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Description:  "Authentication types for JDBC connector",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      nextgen.JDBCAuthTypes.UsernamePassword.String(),
							ValidateFunc: validation.StringInSlice([]string{nextgen.JDBCAuthTypes.UsernamePassword.String(), nextgen.JDBCAuthTypes.ServiceAccount.String()}, false),
						},
						"username": {
							Description:   "The username to use for the database server.",
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.username_ref", "credentials.0.username_password", "credentials.0.service_account"},
							AtLeastOneOf: []string{
								"credentials.0.username",
								"credentials.0.username_ref",
								"credentials.0.username_password",
								"credentials.0.service_account",
							},
						},
						"username_ref": {
							Description:   "The reference to the Harness secret containing the username to use for the database server." + secret_ref_text,
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.username", "credentials.0.username_password", "credentials.0.service_account"},
							AtLeastOneOf: []string{
								"credentials.0.username",
								"credentials.0.username_ref",
								"credentials.0.username_password",
								"credentials.0.service_account",
							},
						},
						"password_ref": {
							Description:   "The reference to the Harness secret containing the password to use for the database server." + secret_ref_text,
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.username_password", "credentials.0.service_account"},
							AtLeastOneOf: []string{
								"credentials.0.password_ref",
								"credentials.0.username_password",
								"credentials.0.service_account",
							},
						},
						"username_password": {
							Description:   "Authenticate using username password.",
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.username", "credentials.0.username_ref", "credentials.0.password_ref", "credentials.0.service_account"},
							AtLeastOneOf: []string{
								"credentials.0.username_password",
								"credentials.0.service_account",
								"credentials.0.password_ref",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Description:   "Username to use for authentication.",
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"credentials.0.username_password.0.username_ref"},
										AtLeastOneOf: []string{
											"credentials.0.username_password.0.username",
											"credentials.0.username_password.0.username_ref",
										},
									},
									"username_ref": {
										Description:   "Reference to a secret containing the username to use for authentication." + secret_ref_text,
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"credentials.0.username_password.0.username"},
										AtLeastOneOf: []string{
											"credentials.0.username_password.0.username",
											"credentials.0.username_password.0.username_ref",
										},
									},
									"password_ref": {
										Description: "Reference to a secret containing the password to use for authentication." + secret_ref_text,
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
						"service_account": {
							Description:   "Authenticate using service account.",
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: []string{"credentials.0.username", "credentials.0.username_ref", "credentials.0.password_ref", "credentials.0.username_password"},
							AtLeastOneOf: []string{
								"credentials.0.username_password",
								"credentials.0.service_account",
								"credentials.0.password_ref",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token_ref": {
										Description: "Reference to a secret containing the token to use for authentication." + secret_ref_text,
										Type:        schema.TypeString,
										Required:    true,
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

func resourceConnectorJDBCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := resourceConnectorReadBase(ctx, d, meta, nextgen.ConnectorTypes.JDBC)
	if err != nil {
		return err
	}

	if conn == nil {
		return nil
	}

	if err := readConnectorJDBC(d, conn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorJDBCCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := buildConnectorJDBC(d)

	newConn, err := resourceConnectorCreateOrUpdateBase(ctx, d, meta, conn)
	if err != nil {
		return err
	}

	if err := readConnectorJDBC(d, newConn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildConnectorJDBC(d *schema.ResourceData) *nextgen.ConnectorInfo {
	config := d.Get("credentials").([]interface{})[0].(map[string]interface{})
	authType := config["auth_type"].(string)

	connector1 := &nextgen.ConnectorInfo{
		Type_: nextgen.ConnectorTypes.JDBC,
		JDBC: &nextgen.JdbcConnector{
			Url: d.Get("url").(string),
			// As currently we support through delegate only
			ExecuteOnDelegate:    true,
			IgnoreTestConnection: false,
		},
	}

	switch authType {
	case nextgen.JDBCAuthTypes.UsernamePassword.String():
		{
			connector1.JDBC.Auth = &nextgen.JdbcAuthenticationDto{
				Type_:            nextgen.JDBCAuthTypes.UsernamePassword,
				UsernamePassword: &nextgen.JdbcUserNamePasswordDto{},
			}
			if usernamePasswordConfig, ok := config["username_password"]; ok && len(usernamePasswordConfig.([]interface{})) > 0 {
				config = usernamePasswordConfig.([]interface{})[0].(map[string]interface{})
				if attr, ok := config["username"]; ok {
					connector1.JDBC.Auth.UsernamePassword.Username = attr.(string)
				}
				if attr, ok := config["username_ref"]; ok {
					connector1.JDBC.Auth.UsernamePassword.UsernameRef = attr.(string)
				}
				if attr, ok := config["password_ref"]; ok {
					connector1.JDBC.Auth.UsernamePassword.PasswordRef = attr.(string)
				}
			} else {
				if attr, ok := config["username"]; ok {
					connector1.JDBC.Auth.UsernamePassword.Username = attr.(string)
				}

				if attr, ok := config["username_ref"]; ok {
					connector1.JDBC.Auth.UsernamePassword.UsernameRef = attr.(string)
				}

				if attr, ok := config["password_ref"]; ok {
					connector1.JDBC.Auth.UsernamePassword.PasswordRef = attr.(string)
				}
			}
		}
	case nextgen.JDBCAuthTypes.ServiceAccount.String():
		{
			connector1.JDBC.Auth = &nextgen.JdbcAuthenticationDto{
				Type_:          nextgen.JDBCAuthTypes.ServiceAccount,
				ServiceAccount: &nextgen.JdbcServiceAccountDto{},
			}
			if serviceAccountConfig, ok := config["service_account"]; ok && len(serviceAccountConfig.([]interface{})) > 0 {
				config = serviceAccountConfig.([]interface{})[0].(map[string]interface{})
				if attr, ok := config["token_ref"]; ok {
					connector1.JDBC.Auth.ServiceAccount.ServiceAccountTokenRef = attr.(string)
				}
			}
		}
	default:
		panic(fmt.Sprintf("unknown jdbc auth method type %s", authType))
	}

	if attr, ok := d.GetOk("delegate_selectors"); ok {
		connector1.JDBC.DelegateSelectors = utils.InterfaceSliceToStringSlice(attr.(*schema.Set).List())
	}

	return connector1
}

func readConnectorJDBC(d *schema.ResourceData, connector *nextgen.ConnectorInfo) error {
	d.Set("url", connector.JDBC.Url)
	d.Set("delegate_selectors", connector.JDBC.DelegateSelectors)
	switch connector.JDBC.Auth.Type_ {
	case nextgen.JDBCAuthTypes.UsernamePassword:
		useFirstClassUserAndPasswordVars := false
		_, ok := d.GetOk("credentials")
		if ok {
			cred := d.Get("credentials").([]interface{})[0].(map[string]interface{})
			passwordRefValue, passwordRefKeyPresent := cred["password_ref"]
			if passwordRefKeyPresent && passwordRefValue != "" {
				useFirstClassUserAndPasswordVars = true
			}
		}
		if useFirstClassUserAndPasswordVars {
			d.Set("credentials", []map[string]interface{}{
				{
					"auth_type":    nextgen.JDBCAuthTypes.UsernamePassword.String(),
					"username":     connector.JDBC.Auth.UsernamePassword.Username,
					"username_ref": connector.JDBC.Auth.UsernamePassword.UsernameRef,
					"password_ref": connector.JDBC.Auth.UsernamePassword.PasswordRef,
				},
			})
		} else {
			d.Set("credentials", []map[string]interface{}{
				{
					"auth_type": nextgen.JDBCAuthTypes.UsernamePassword.String(),
					"username_password": []map[string]interface{}{
						{
							"username":     connector.JDBC.Auth.UsernamePassword.Username,
							"username_ref": connector.JDBC.Auth.UsernamePassword.UsernameRef,
							"password_ref": connector.JDBC.Auth.UsernamePassword.PasswordRef,
						},
					},
				},
			})
		}
	case nextgen.JDBCAuthTypes.ServiceAccount:
		d.Set("credentials", []map[string]interface{}{
			{
				"auth_type": nextgen.JDBCAuthTypes.ServiceAccount.String(),
				"service_account": []map[string]interface{}{
					{
						"token_ref": connector.JDBC.Auth.ServiceAccount.ServiceAccountTokenRef,
					},
				},
			},
		})
	default:
		return fmt.Errorf("unknown jdbc auth method type %s", connector.JDBC.Auth.Type_)
	}

	return nil
}
