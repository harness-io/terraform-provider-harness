package provider

import (
	"context"

	"github.com/harness-io/harness-go-sdk/harness/api"
	"github.com/harness-io/harness-go-sdk/harness/api/cac"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAWSLambdaService() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource for creating an AWS Lambda service",
		CreateContext: resourceAWSLambdaServiceCreate,
		ReadContext:   resourceAWSLambdaServiceRead,
		UpdateContext: resourceAWSLambdaServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Schema:        commonServiceSchema(),
	}
}

func resourceAWSLambdaServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	svcId := d.Get("id").(string)
	appId := d.Get("app_id").(string)

	svc, err := c.ConfigAsCode().GetServiceById(appId, svcId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", svc.Name)
	d.Set("app_id", svc.ApplicationId)
	d.Set("description", svc.Description)

	if vars := flattenServiceVariables(svc.ConfigVariables); len(vars) > 0 {
		d.Set("variable", vars)
	}

	return nil
}

func resourceAWSLambdaServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	// Setup the object to be created
	svcInput := &cac.Service{
		Name:           d.Get("name").(string),
		ArtifactType:   cac.ArtifactTypes.AWSLambda,
		DeploymentType: cac.DeploymentTypes.AWSLambda,
		ApplicationId:  d.Get("app_id").(string),
		Description:    d.Get("description").(string),
	}

	if vars := d.Get("variable"); vars != nil {
		svcInput.ConfigVariables = expandServiceVariables(vars.(*schema.Set).List())
	}

	// Create Service
	newSvc, err := c.ConfigAsCode().UpsertService(svcInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newSvc.Id)

	return nil
}

func resourceAWSLambdaServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	// Setup the object to create
	svcInput := &cac.Service{
		Name:           d.Get("name").(string),
		ArtifactType:   cac.ArtifactTypes.AWSLambda,
		DeploymentType: cac.DeploymentTypes.AWSLambda,
		ApplicationId:  d.Get("app_id").(string),
		Description:    d.Get("description").(string),
	}

	if vars := d.Get("variable"); vars != nil {
		svcInput.ConfigVariables = expandServiceVariables(vars.(*schema.Set).List())
	}

	// Create Service
	newSvc, err := c.ConfigAsCode().UpsertService(svcInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newSvc.Id)

	return nil
}
