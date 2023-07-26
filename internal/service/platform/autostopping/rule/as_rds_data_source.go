package as_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRDSRule() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a Harness Variable.",
		ReadContext: resourceASRuleRead,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Description: "Unique identifier of the resource",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"name": {
				Description: "Name of the rule",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloud_connector_id": {
				Description: "Id of the cloud connector",
				Type:        schema.TypeString,
				Required:    true,
			},
			"idle_time_mins": {
				Description: "Idle time in minutes. This is the time that the AutoStopping rule waits before stopping the idle instances.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     15,
			},
			"database": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the database",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region to which database belong to",
						},
					},
				},
			},
			"tcp": {
				Description: "TCP routing configuration",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_id": {
							Description: "Id of the Proxy",
							Type:        schema.TypeString,
							Required:    true,
						},
						"forward_rule": {
							Description: "Additional tcp forwarding rules",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connect_on": {
										Description: "Port to listen on the proxy",
										Type:        schema.TypeInt,
										Optional:    true,
									},
									"port": {
										Description: "Port to listen on the vm",
										Type:        schema.TypeInt,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"depends": {
				Description: "Dependent rules",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Description: "Rule id of the dependent rule",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"delay_in_sec": {
							Description: "Number of seconds the rule should wait after warming up the dependent rule",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5,
						},
					},
				},
			},
		},
	}

	return resource
}
