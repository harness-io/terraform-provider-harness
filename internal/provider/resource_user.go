package provider

import (
	"context"
	"strings"

	"github.com/harness-io/harness-go-sdk/harness/api"
	"github.com/harness-io/harness-go-sdk/harness/api/graphql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Resource for creating a Harness user",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Unique identifier of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email of the user.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"is_email_verified": {
				Description: "Flag indicating whether or not the users email has been verified.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_imported_from_identity_provider": {
				Description: "Flag indicating whether or not the user was imported from an identity provider.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_password_expired": {
				Description: "Flag indicating whether or not the users password has expired.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_two_factor_auth_enabled": {
				Description: "Flag indicating whether or not two-factor authentication is enabled for the user.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_user_locked": {
				Description: "Flag indicating whether or not the user is locked out.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	input := &graphql.CreateUserInput{
		Name:  d.Get("name").(string),
		Email: d.Get("email").(string),
	}

	user, err := c.Users().CreateUser(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.Id)

	return nil
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	email := d.Get("email").(string)

	user, err := c.Users().GetUserByEmail(email)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.Id)
	d.Set("name", user.Name)
	d.Set("email", user.Email)
	d.Set("is_email_verified", user.IsEmailVerified)
	d.Set("is_imported_from_identity_provider", user.IsImportedFromIdentityProvider)
	d.Set("is_password_expired", user.IsPasswordExpired)
	d.Set("is_two_factor_auth_enabled", user.IsTwoFactorAuthenticationEnabled)
	d.Set("is_user_locked", user.IsUserLocked)

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	input := &graphql.UpdateUserInput{
		Name: d.Get("name").(string),
		Id:   d.Id(),
	}

	_, err := c.Users().UpdateUser(input)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	if err := c.Users().DeleteUser(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
