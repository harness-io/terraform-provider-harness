package featureflagtargetgroup

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/harness/terraform-provider-harness/internal/service/platform/feature_flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"log"
	"net/http"
	"time"
)

// ResourceFeatureFlagTargetGroup ...
func ResourceFeatureFlagTargetGroup() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness Feature Flag Target Group.",

		ReadContext:   resourceFeatureFlagTargetGroupRead,
		CreateContext: resourceFeatureFlagTargetCreate,
		UpdateContext: resourceFeatureFlagTargetGroupUpdate,
		DeleteContext: resourceFeatureFlagTargetGroupDelete,
		Importer:      helpers.ProjectResourceImporter,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Description: "The unique identifier of the feature flag target group.",
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
			"environment": {
				Description: "Environment Identifier",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"account_id": {
				Description: "Account Identifier",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the feature flag target group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"included": {
				Description: "A list of targets to include in the target group",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				MinItems:    0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded": {
				Description: "A list of targets to exclude from the target group",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				MinItems:    0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rule": {
				Description: "The list of rules used to include targets in the target group.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The rule ID. Gets auto-generated by the system",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"attribute": {
							Description: "The attribute to use in the clause.  This can be any target attribute",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
						},
						"negate": {
							Description: "Is the operation negated?",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    false,
						},
						"op": {
							Description: "The type of operation such as equals, starts_with, contains",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
						},
						"values": {
							Description: "The values that are compared against the operator",
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    false,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}

	return resource
}

type PatchInstruction struct {
	Kind       string                 `json:"kind"`
	Parameters map[string]interface{} `json:"parameters"`
}

// FFTargetGroupQueryParameters ...
type FFTargetGroupQueryParameters struct {
	Identifier  string `json:"identifier,omitempty"`
	OrgID       string `json:"orgId,omitempty"`
	Project     string `json:"project,omitempty"`
	AcountID    string `json:"accountId,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// FFTargetGroupOpts ...
type FFTargetGroupOpts struct {
	Identifier string            `json:"identifier,omitempty"`
	Name       string            `json:"name,omitempty"`
	Included   []*string         `json:"included,omitempty"`
	Excluded   []*string         `json:"excluded,omitempty"`
	Rules      []*nextgen.Clause `json:"rules,omitempty"`
}

// SegmentRequest ...
type SegmentRequest struct {
	Identifier  string            `json:"identifier,omitempty"`
	Project     string            `json:"project,omitempty"`
	Environment string            `json:"environment,omitempty"`
	Name        string            `json:"name,omitempty"`
	Included    []*string         `json:"included,omitempty"`
	Excluded    []*string         `json:"excluded,omitempty"`
	Rules       []*nextgen.Clause `json:"rules,omitempty"`
}

func resourceFeatureFlagTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		d.MarkNewResource()
		return nil
	}

	qp := buildFFTargetGroupQueryParameters(d)

	segment, httpResp, err := c.TargetGroupsApi.GetSegment(ctx, c.AccountId, qp.OrgID, id, qp.Project, qp.Environment)
	if err != nil {
		return feature_flag.HandleCFApiError(err, d, httpResp)
	}
	readFeatureFlagTargetGroup(d, &segment, qp)

	return nil
}

// resourceFeatureFlagTargetGroupCreate ...
func resourceFeatureFlagTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	var err error
	var httpResp *http.Response
	var segment nextgen.Segment
	segmentRequest := buildSegmentRequest(d)
	qp := buildFFTargetGroupQueryParameters(d)
	id := d.Id()
	if id == "" {
		id = d.Get("identifier").(string)
		d.MarkNewResource()
	}

	httpResp, err = c.TargetGroupsApi.CreateSegment(ctx, segmentRequest, c.AccountId, qp.OrgID)

	if err != nil {
		// handle conflict
		if httpResp != nil && httpResp.StatusCode == 409 {
			return diag.Errorf("A target group with identifier [%s] orgIdentifier [%s] project [%s] environment [%s] already exists", segmentRequest.Identifier, qp.OrgID, qp.Project, segmentRequest.Environment)
		}
		return feature_flag.HandleCFApiError(err, d, httpResp)
	}

	if httpResp.StatusCode != http.StatusCreated {
		return diag.Errorf("createstatus: %s", httpResp.Status)
	}

	segment, httpResp, err = c.TargetGroupsApi.GetSegment(ctx, c.AccountId, qp.OrgID, id, qp.Project, qp.Environment)
	if err != nil {
		return feature_flag.HandleCFApiError(err, d, httpResp)
	}

	readFeatureFlagTargetGroup(d, &segment, qp)

	return nil
}

func resourceFeatureFlagTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		return nil
	}

	qp := buildFFTargetGroupQueryParameters(d)
	opts := buildFFTargetGroupPatchOpts(d)

	if opts != nil {

		_, httpResp, err := c.TargetGroupsApi.PatchSegment(ctx, c.AccountId, qp.OrgID, qp.Project, qp.Environment, id, opts)
		if err != nil {
			return feature_flag.HandleCFApiError(err, d, httpResp)
		}
		time.Sleep(1 * time.Second)
	}

	segment, httpResp, err := c.TargetGroupsApi.GetSegment(ctx, c.AccountId, qp.OrgID, id, qp.Project, qp.Environment)
	if err != nil {
		body, _ := io.ReadAll(httpResp.Body)
		return diag.Errorf("readstatus: %s, \nBody:%s", httpResp.Status, body)
	}

	readFeatureFlagTargetGroup(d, &segment, qp)

	return nil
}

func resourceFeatureFlagTargetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()
	if id == "" {
		return nil
	}

	qp := buildFFTargetGroupQueryParameters(d)

	httpResp, err := c.TargetGroupsApi.DeleteSegment(ctx, c.AccountId, qp.OrgID, id, qp.Project, qp.Environment)
	if err != nil {
		return feature_flag.HandleCFApiError(err, d, httpResp)
	}

	return nil
}

// readFeatureFlagTargetGroupRule ...
func readFeatureFlagTargetGroup(d *schema.ResourceData, segment *nextgen.Segment, qp *FFTargetGroupQueryParameters) {
	d.SetId(segment.Identifier)
	d.Set("identifier", segment.Identifier)
	d.Set("org_id", qp.OrgID)
	d.Set("project_id", qp.Project)
	d.Set("account_id", qp.AcountID)
	d.Set("environment", segment.Environment)
	d.Set("name", segment.Name)
	d.Set("included", targetsToStrings(segment.Included))
	d.Set("excluded", targetsToStrings(segment.Excluded))
	rules := expandRules(segment.Rules)
	d.Set("rule", rules)
}

func targetsToStrings(targets []nextgen.Target) []string {
	var targetStrings []string
	for _, target := range targets {
		targetStrings = append(targetStrings, target.Identifier)
	}
	return targetStrings
}

func expandRules(rules []nextgen.Clause) []interface{} {
	var result []interface{}
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"attribute": rule.Attribute,
			"negate":    rule.Negate,
			"op":        rule.Op,
			"values":    rule.Values,
			"id":        rule.Id,
		})
	}
	return result
}

// buildFFTargetGroupQueryParameters ...
func buildFFTargetGroupQueryParameters(d *schema.ResourceData) *FFTargetGroupQueryParameters {
	return &FFTargetGroupQueryParameters{
		Identifier:  d.Get("identifier").(string),
		OrgID:       d.Get("org_id").(string),
		Project:     d.Get("project_id").(string),
		AcountID:    d.Get("account_id").(string),
		Environment: d.Get("environment").(string),
	}
}

// buildSegmentRequest builds a SegmentRequest from a ResourceData
func buildSegmentRequest(d *schema.ResourceData) *SegmentRequest {
	opts := &SegmentRequest{
		Identifier:  d.Get("identifier").(string),
		Project:     d.Get("project_id").(string),
		Environment: d.Get("environment").(string),
		Name:        d.Get("name").(string),
	}

	if included, ok := d.GetOk("included"); ok {
		var targets = make([]*string, 0)
		for _, target := range included.([]interface{}) {
			targets = append(targets, aws.String(target.(string)))
		}
		opts.Included = targets
	}

	if excluded, ok := d.GetOk("excluded"); ok {
		var targets = make([]*string, 0)
		for _, target := range excluded.([]interface{}) {
			targets = append(targets, aws.String(target.(string)))
		}
		opts.Excluded = targets
	}

	if rules, ok := d.GetOk("rule"); ok {
		var rulesList = make([]*nextgen.Clause, 0)
		for _, rule := range rules.([]interface{}) {
			var values []string
			for _, value := range rule.(map[string]interface{})["values"].([]interface{}) {
				values = append(values, value.(string))
			}
			rule := &nextgen.Clause{
				Attribute: rule.(map[string]interface{})["attribute"].(string),
				Negate:    rule.(map[string]interface{})["negate"].(bool),
				Op:        rule.(map[string]interface{})["op"].(string),
				Values:    values,
			}
			rulesList = append(rulesList, rule)
		}
		opts.Rules = rulesList
	}

	return opts
}

// buildFFTargetGroupPatchOpts ...
func buildFFTargetGroupPatchOpts(d *schema.ResourceData) *nextgen.TargetGroupsApiPatchSegmentOpts {
	var instructions []PatchInstruction
	if d.HasChange("included") {
		currentIncludeRules, newIncludeRules := []string{}, []string{}
		oldState, newState := d.GetChange("included")
		// cast to array strings
		if oldState != nil {
			i := oldState.([]interface{})
			for _, v := range i {
				currentIncludeRules = append(currentIncludeRules, v.(string))
			}
		}
		if newState != nil {
			i := newState.([]interface{})
			for _, v := range i {
				newIncludeRules = append(newIncludeRules, v.(string))
			}
		}

		extraIncludeRules, missingIncludeRules := IncludeRuleDiffs(currentIncludeRules, newIncludeRules)
		// remove extra include rules - should be done first so no conflicts of targets being in both lists
		if len(extraIncludeRules) != 0 {
			instructions = append([]PatchInstruction{{
				Kind:       "removeFromIncludeList",
				Parameters: map[string]interface{}{"targets": extraIncludeRules},
			}}, instructions...)
		}
		// add missing include rules
		if len(missingIncludeRules) != 0 {
			instructions = append(instructions, PatchInstruction{
				Kind:       "addToIncludeList",
				Parameters: map[string]interface{}{"targets": missingIncludeRules},
			})
		}
	}

	if d.HasChange("excluded") {
		currentExcludeRules, newExcludeRules := []string{}, []string{}
		oldState, newState := d.GetChange("excluded")
		// cast to array strings
		if oldState != nil {
			i := oldState.([]interface{})
			for _, v := range i {
				currentExcludeRules = append(currentExcludeRules, v.(string))
			}
		}
		if newState != nil {
			i := newState.([]interface{})
			for _, v := range i {
				newExcludeRules = append(newExcludeRules, v.(string))
			}
		}

		extraExcludeRules, missingExcludeRules := IncludeRuleDiffs(currentExcludeRules, newExcludeRules)
		// remove extra exclude rules - should be done first so no conflicts of targets being in both lists
		if len(extraExcludeRules) != 0 {
			instructions = append([]PatchInstruction{{
				Kind:       "removeFromExcludeList",
				Parameters: map[string]interface{}{"targets": extraExcludeRules},
			}}, instructions...)
		}
		// add missing exclude rules
		if len(missingExcludeRules) != 0 {
			instructions = append(instructions, PatchInstruction{
				Kind:       "addToExcludeList",
				Parameters: map[string]interface{}{"targets": missingExcludeRules},
			})
		}
	}

	if d.HasChange("rule") {
		currentRules, newRules := []nextgen.Clause{}, []nextgen.Clause{}
		oldState, newState := d.GetChange("rule")
		if oldState != nil {
			i := oldState.([]interface{})
			for _, v := range i {
				values := []string{}
				for _, value := range v.(map[string]interface{})["values"].([]interface{}) {
					values = append(values, value.(string))
				}
				currentRules = append(currentRules, nextgen.Clause{
					Attribute: v.(map[string]interface{})["attribute"].(string),
					Negate:    v.(map[string]interface{})["negate"].(bool),
					Op:        v.(map[string]interface{})["op"].(string),
					Id:        v.(map[string]interface{})["id"].(string),
					Values:    values,
				})
			}
		}
		if newState != nil {
			i := newState.([]interface{})
			for _, v := range i {
				values := []string{}
				for _, value := range v.(map[string]interface{})["values"].([]interface{}) {
					values = append(values, value.(string))
				}
				newRules = append(newRules, nextgen.Clause{
					Attribute: v.(map[string]interface{})["attribute"].(string),
					Negate:    v.(map[string]interface{})["negate"].(bool),
					Op:        v.(map[string]interface{})["op"].(string),
					Id:        v.(map[string]interface{})["id"].(string),
					Values:    values,
				})
			}
		}

		extraRules, missingRules := RuleDiffs(&currentRules, &newRules)
		// remove extra rules - should be done first so no conflicts of targets being in both lists
		if len(extraRules) != 0 {
			for _, rule := range extraRules {
				if rule.Id == "" {
					fmt.Printf("rule id is empty for %v, skipping", rule)
					continue
				}
				instructions = append(instructions, PatchInstruction{
					Kind:       "removeClause",
					Parameters: map[string]interface{}{"clauseID": rule.Id},
				})
			}
		}
		// add missing rules
		if len(missingRules) != 0 {
			for _, rule := range missingRules {
				instructions = append(instructions, PatchInstruction{
					Kind:       "addClause",
					Parameters: map[string]interface{}{"attribute": rule.Attribute, "op": rule.Op, "values": rule.Values, "negate": rule.Negate},
				})
			}
		}

	}

	if len(instructions) == 0 {
		return nil
	}

	log.Println("segment update request body")
	jsonData, _ := json.Marshal(instructions)
	log.Println(string(jsonData))

	return &nextgen.TargetGroupsApiPatchSegmentOpts{
		Body: optional.NewInterface(map[string]interface{}{"instructions": instructions}),
	}
}

func RuleDiffs(rulesA, rulesB *[]nextgen.Clause) (extraRules, missingRules []nextgen.Clause) {
	rulesMapA, rulesMapB := rulesToMap(rulesA), rulesToMap(rulesB)

	for rule, clause := range rulesMapA {
		if _, ok := rulesMapB[rule]; !ok {
			extraRules = append(extraRules, clause)
		}
	}

	for rule, clause := range rulesMapB {
		if _, ok := rulesMapA[rule]; !ok {
			missingRules = append(missingRules, clause)
		}
	}

	return extraRules, missingRules
}

func rulesToMap(rules *[]nextgen.Clause) map[string]nextgen.Clause {
	rulesMap := map[string]nextgen.Clause{}
	if rules == nil {
		return rulesMap
	}
	for _, clause := range *rules {
		rulesMap[fmt.Sprintf("%s_%s_%s", clause.Attribute, clause.Op, clause.Values)] = clause
	}
	return rulesMap
}

func IncludeRuleDiffs(targetsA, targetsB []string) (extraIncludeRules, missingIncludeRules []string) {
	includeRulesA, includeRulesB := arrayToMap(targetsA), arrayToMap(targetsB)

	for i := range targetsA {
		if _, ok := includeRulesB[targetsA[i]]; !ok {
			extraIncludeRules = append(extraIncludeRules, targetsA[i])
		}
	}

	for i := range targetsB {
		if _, ok := includeRulesA[targetsB[i]]; !ok {
			missingIncludeRules = append(missingIncludeRules, targetsB[i])
		}
	}

	return extraIncludeRules, missingIncludeRules
}

func arrayToMap(identifiers []string) map[string]struct{} {
	m := map[string]struct{}{}
	if identifiers == nil {
		return m
	}
	for _, identifier := range identifiers {
		m[identifier] = struct{}{}
	}
	return m
}
