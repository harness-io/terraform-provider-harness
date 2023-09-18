package feature_flag

import (
	"context"
	"net/http"
	"strings"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFeatureFlag() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for managing Feature Flags.",

		ReadContext:   resourceFeatureFlagRead,
		DeleteContext: resourceFeatureFlagDelete,
		CreateContext: resourceFeatureFlagCreate,
		UpdateContext: resourceFeatureFlagUpdate,
		Importer:      helpers.ProjectResourceImporter,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Description: "Identifier of the Feature Flag",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Name of the Feature Flag",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"org_id": {
				Description: "Organization Identifier",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Description: "Project Identifier",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"archived": {
				Description: "Whether or not the flag is archived",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"default_off_variation": {
				Description: "Which of the variations to use when the flag is toggled to off state",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"default_on_variation": {
				Description: "Which of the variations to use when the flag is toggled to on state",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"git_details": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commit_msg": {
							Description: "The commit message to use as part of a gitsync operation",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"kind": {
				Description: "The type of data the flag represents. Valid values are `boolean`, `int`, `string`, `json`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"owner": {
				Description: "The owner of the flag",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"permanent": {
				Description: "Whether or not the flag is permanent. If it is, it will never be flagged as stale",
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
			},
			"environment": {
				Description: "Environment Identifier",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"variation": {
				Description: "The options available for your flag",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Description: "The identifier of the variation",
							Type:        schema.TypeString,
							Required:    true,
						},
						"description": {
							Description: "The description of the variation",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "The user friendly name of the variation",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "The value of the variation",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"add_target_rule": {
				Description: "The targeting rules for the flag",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variation": {
							Description: "The identifier of the variation. Valid values are `enabled`, `disabled`",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"targets": {
							Description: "The targets of the rule",
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    0,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"add_target_group_rule": {
				Description: "The targeting rules for the flag",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Description: "The name of the target group",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"variation": {
							Description: "The identifier of the variation. Valid values are `enabled`, `disabled`",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"distribution": {
							Description: "The distribution of the rule",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"variations": {
										Description: "The variations of the rule",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"variation": {
													Description: "The identifier of the variation",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"weight": {
													Description: "The weight of the variation",
													Type:        schema.TypeInt,
													Optional:    true,
												},
											},
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

	return resource
}

type FFQueryParameters struct {
	Identifier     string
	OrganizationId string
	ProjectId      string
	Environment    string
}

// KindMap is a map of the kind to the actual kind
var KindMap = map[string]string{
	"removeTargets": "removeTargetsToVariationTargetMap",
	"removeRule":    "removeRule",
	"addRule":       "addRule",
	"addTargets":    "addTargetsToVariationTargetMap",
}

// TargetRules is the target rules for the feature flag
type TargetRules struct {
	Kind      string   `json:"kind,omitempty"`
	Variation string   `json:"variation,omitempty"`
	Targets   []string `json:"targets,omitempty"`
}

// Variation is the variation for the feature flag
type Variation struct {
	Variation string `json:"variation,omitempty"`
	Weight    int    `json:"weight,omitempty"`
}

// Distribution is the distribution for the feature flag
type Distribution struct {
	BuckedBy   string      `json:"buckedBy,omitempty"`
	Variations []Variation `json:"variations,omitempty"`
}

// TargetGroupRules is the target group rules for the feature flag
type TargetGroupRules struct {
	Kind      string `json:"kind,omitempty"`
	GroupName string `json:"groupName,omitempty"`
	Variation string `json:"variation,omitempty"`
}

// Serve ...
type Serve struct {
	Variation    string       `json:"variation,omitempty"`
	Distribution Distribution `json:"distribution,omitempty"`
}

// Parameter ...
type Parameter struct {
	Variation string           `json:"variation,omitempty"`
	Targets   []string         `json:"targets,omitempty"`
	Priority  string           `json:"priority,omitempty"`
	Clauses   []nextgen.Clause `json:"clauses,omitempty"`
	Serve     Serve            `json:"serve,omitempty"`
}

// Instruction defines the instruction for the feature flag
type Instruction struct {
	Kind       string    `json:"kind,omitempty"`
	Parameters Parameter `json:"parameters,omitempty"`
}

type FFOpts struct {
	Identifier          string              `json:"identifier"`
	Name                string              `json:"name"`
	Description         string              `json:"description,omitempty"`
	Environment         string              `json:"environment,omitempty"`
	Archived            bool                `json:"archived,omitempty"`
	DefaultOffVariation string              `json:"defaultOffVariation"`
	DefaultOnVariation  string              `json:"defaultOnVariation"`
	GitDetails          nextgen.GitDetails  `json:"gitDetails,omitempty"`
	Kind                string              `json:"kind"`
	Owner               string              `json:"owner,omitempty"`
	Permanent           bool                `json:"permanent"`
	Project             string              `json:"project"`
	Variations          []nextgen.Variation `json:"variations"`
	Instructions        []Instruction       `json:"instructions,omitempty"`
}

// FFPatchOpts is the options for patching a feature flag
type FFPatchOpts struct {
	Identifier          string              `json:"identifier"`
	Name                string              `json:"name"`
	Environment         string              `json:"environment,omitempty"`
	Description         string              `json:"description,omitempty"`
	Archived            bool                `json:"archived,omitempty"`
	DefaultOffVariation string              `json:"defaultOffVariation"`
	DefaultOnVariation  string              `json:"defaultOnVariation"`
	GitDetails          nextgen.GitDetails  `json:"gitDetails,omitempty"`
	Kind                string              `json:"kind"`
	Owner               string              `json:"owner,omitempty"`
	Permanent           bool                `json:"permanent"`
	Project             string              `json:"project"`
	Variations          []nextgen.Variation `json:"variations"`
	Instructions        []Instruction       `json:"instructions,omitempty"`
}

func resourceFeatureFlagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		return nil
	}

	qp := buildFFQueryParameters(d)
	opts := buildFFPatchOpts(d)

	feature, httpResp, err := c.FeatureFlagsApi.PatchFeature(ctx, c.AccountId, qp.OrganizationId, qp.ProjectId, id, opts)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readFeatureFlag(d, &feature, qp)

	return nil
}

func resourceFeatureFlagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		d.MarkNewResource()
		return nil
	}

	qp := buildFFQueryParameters(d)
	opts := buildFFReadOpts(d)

	resp, httpResp, err := c.FeatureFlagsApi.GetFeatureFlag(ctx, id, c.AccountId, qp.OrganizationId, qp.ProjectId, opts)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readFeatureFlag(d, &resp, qp)

	return nil
}

func resourceFeatureFlagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		id = d.Get("identifier").(string)
		d.MarkNewResource()
	}

	qp := buildFFQueryParameters(d)
	opts := buildFFCreateOpts(d)
	readOpts := buildFFReadOpts(d)

	var err error
	var resp nextgen.Feature
	var httpResp *http.Response

	httpResp, err = c.FeatureFlagsApi.CreateFeatureFlag(ctx, c.AccountId, qp.OrganizationId, opts)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	resp, httpResp, err = c.FeatureFlagsApi.GetFeatureFlag(ctx, id, c.AccountId, qp.OrganizationId, qp.ProjectId, readOpts)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readFeatureFlag(d, &resp, qp)

	// patchOpts := buildFFPatchOpts(d)

	// update the feature flag with the git details
	// feature, httpResp, err := c.FeatureFlagsApi.PatchFeature(ctx, c.AccountId, qp.OrganizationId, qp.ProjectId, id, patchOpts)

	// if err != nil {
	//	return helpers.HandleApiError(err, d, httpResp)
	// }

	return nil
}

func resourceFeatureFlagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		return nil
	}
	qp := buildFFQueryParameters(d)

	httpResp, err := c.FeatureFlagsApi.DeleteFeatureFlag(ctx, d.Id(), c.AccountId, qp.OrganizationId, qp.ProjectId, &nextgen.FeatureFlagsApiDeleteFeatureFlagOpts{CommitMsg: optional.EmptyString()})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func readFeatureFlag(d *schema.ResourceData, flag *nextgen.Feature, qp *FFQueryParameters) {
	d.SetId(flag.Identifier)
	d.Set("identifier", flag.Identifier)
	d.Set("name", flag.Name)
	d.Set("project_id", flag.Project)
	d.Set("default_on_variation", flag.DefaultOnVariation)
	d.Set("default_off_variation", flag.DefaultOffVariation)
	d.Set("description", flag.Description)
	d.Set("kind", flag.Kind)
	d.Set("permanent", flag.Permanent)
	d.Set("owner", strings.Join(flag.Owner, ","))
	d.Set("org_id", qp.OrganizationId)
	d.Set("variation", expandVariations(flag.Variations))
	d.Set("environment", qp.Environment)
}

func expandVariations(variations []nextgen.Variation) []interface{} {
	var result []interface{}
	for _, variation := range variations {
		result = append(result, map[string]interface{}{
			"identifier":  variation.Identifier,
			"name":        variation.Name,
			"description": variation.Description,
			"value":       variation.Value,
		})
	}

	return result
}

func buildFFQueryParameters(d *schema.ResourceData) *FFQueryParameters {
	return &FFQueryParameters{
		Identifier:     d.Get("identifier").(string),
		OrganizationId: d.Get("org_id").(string),
		ProjectId:      d.Get("project_id").(string),
		Environment:    d.Get("environment").(string),
	}
}

func buildFFCreateOpts(d *schema.ResourceData) *nextgen.FeatureFlagsApiCreateFeatureFlagOpts {
	opts := &FFOpts{
		Identifier:          d.Get("identifier").(string),
		Name:                d.Get("name").(string),
		DefaultOffVariation: d.Get("default_off_variation").(string),
		DefaultOnVariation:  d.Get("default_on_variation").(string),
		Project:             d.Get("project_id").(string),
		Kind:                d.Get("kind").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		opts.Description = desc.(string)
	}

	if owner, ok := d.GetOk("owner"); ok {
		opts.Owner = owner.(string)
	}

	if archived, ok := d.GetOk("archived"); ok {
		opts.Archived = archived.(bool)
	}

	var variations []nextgen.Variation
	variationsData := d.Get("variation").([]interface{})
	for _, variationData := range variationsData {
		vMap := variationData.(map[string]interface{})
		variation := nextgen.Variation{
			Identifier:  vMap["identifier"].(string),
			Value:       vMap["value"].(string),
			Name:        vMap["name"].(string),
			Description: vMap["description"].(string),
		}
		variations = append(variations, variation)
	}
	opts.Variations = variations

	return &nextgen.FeatureFlagsApiCreateFeatureFlagOpts{
		Body: optional.NewInterface(opts),
	}
}

func buildFFPatchOpts(d *schema.ResourceData) *nextgen.FeatureFlagsApiPatchFeatureOpts {
	opts := &FFPatchOpts{
		Identifier:          d.Get("identifier").(string),
		Name:                d.Get("name").(string),
		DefaultOffVariation: d.Get("default_off_variation").(string),
		DefaultOnVariation:  d.Get("default_on_variation").(string),
		Project:             d.Get("project_id").(string),
		Kind:                d.Get("kind").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		opts.Description = desc.(string)
	}

	if owner, ok := d.GetOk("owner"); ok {
		opts.Owner = owner.(string)
	}

	if archived, ok := d.GetOk("archived"); ok {
		opts.Archived = archived.(bool)
	}

	var variations []nextgen.Variation
	variationsData := d.Get("variation").([]interface{})
	for _, variationData := range variationsData {
		vMap := variationData.(map[string]interface{})
		variation := nextgen.Variation{
			Identifier:  vMap["identifier"].(string),
			Value:       vMap["value"].(string),
			Name:        vMap["name"].(string),
			Description: vMap["description"].(string),
		}
		variations = append(variations, variation)
	}
	opts.Variations = variations

	var targetRules []TargetRules
	if targetRulesData, ok := d.GetOk("add_target_rule"); ok {
		for _, targetRuleData := range targetRulesData.([]interface{}) {
			vMap := targetRuleData.(map[string]interface{})
			var targets []string = make([]string, 0)
			for _, target := range vMap["targets"].([]interface{}) {
				targets = append(targets, target.(string))
			}
			targetRule := TargetRules{
				Kind:      "addTargetsToVariationTargetMap",
				Variation: vMap["variation"].(string),
				Targets:   targets,
			}
			targetRules = append(targetRules, targetRule)
		}
	}

	var targetGroupRules []TargetGroupRules
	var distribution Distribution
	if targetGroupRulesData, ok := d.GetOk("add_target_group_rule"); ok {
		for _, targetGroupRuleData := range targetGroupRulesData.([]interface{}) {
			vMap := targetGroupRuleData.(map[string]interface{})
			targetGroupRule := TargetGroupRules{
				Kind:      "addRule",
				GroupName: vMap["group_name"].(string),
				Variation: vMap["variation"].(string),
			}

			for _, distributionData := range vMap["distribution"].([]interface{}) {
				vMap := distributionData.(map[string]interface{})
				distribution = Distribution{
					BuckedBy: "identifier",
				}
				var variations []Variation
				for _, variationData := range vMap["variations"].([]interface{}) {
					vMap := variationData.(map[string]interface{})
					variation := Variation{
						Variation: vMap["variation"].(string),
						Weight:    vMap["weight"].(int),
					}
					variations = append(variations, variation)
				}
				distribution.Variations = variations
			}
			targetGroupRules = append(targetGroupRules, targetGroupRule)
		}
	}

	var instructions []Instruction
	for _, target := range targetRules {
		instruction := Instruction{
			Kind: target.Kind,
			Parameters: Parameter{
				Variation: target.Variation,
				Targets:   target.Targets,
			},
		}
		instructions = append(instructions, instruction)
	}

	for _, targetGroup := range targetGroupRules {
		instruction := Instruction{
			Kind: targetGroup.Kind,
			Parameters: Parameter{
				Serve: Serve{
					Variation:    targetGroup.Variation,
					Distribution: distribution,
				},
				Clauses: []nextgen.Clause{
					{
						Op:     "segmentMatch",
						Values: []string{targetGroup.GroupName},
					},
				},
			},
		}
		instructions = append(instructions, instruction)
	}

	opts.Instructions = instructions

	return &nextgen.FeatureFlagsApiPatchFeatureOpts{
		Body:                  optional.NewInterface(opts),
		EnvironmentIdentifier: optional.NewString(d.Get("environment").(string)),
	}
}

func buildFFReadOpts(d *schema.ResourceData) *nextgen.FeatureFlagsApiGetFeatureFlagOpts {

	return &nextgen.FeatureFlagsApiGetFeatureFlagOpts{
		EnvironmentIdentifier: optional.EmptyString(),
	}

}
