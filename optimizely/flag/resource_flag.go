package flag

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Flag struct {
	ID           int64                         `json:"id"`
	ProjectId    int                           `json:"project_id"`
	Name         string                        `json:"name"`
	Description  string                        `json:"description"`
	Key          string                        `json:"key"`
	Archived     bool                          `json:"archived"`
	Variables    map[string]VariableSchema     `json:"variable_definitions"`
	Variations   []Variation                   `json:"variations"`
	Environments map[string]FeatureEnvironment `json:"environments"`
}

func ResourceFeature() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Feature key, also acts as it's unique ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Human readable name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A description of this feature",
			},
			"variable_schema": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"default_value": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"variations": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variation": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"variables": {
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"environments": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"audience": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"percentage_included": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"deliver": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceFeatureCreate,
		ReadContext:   resourceFeatureRead,
		DeleteContext: resourceFeatureDelete,
		UpdateContext: resourceFeatureUpdate,
	}
}

func parseFlag(d *schema.ResourceData) Flag {
	variablesSchema := parseVariableSchema(d)
	variations := parseVariation(d)
	envs := parseEnvironment(d)

	return Flag{
		ProjectId:    d.Get("project").(int),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Key:          d.Get("key").(string),
		Archived:     false,
		Variables:    variablesSchema,
		Variations:   variations,
		Environments: envs,
	}
}

func resourceFeatureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(FlagClient)

	flag := parseFlag(d)

	featResp, err := client.CreateFlag(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create flag in Optimizely: %+v", err),
		})

		return diags
	}

	for _, variation := range flag.Variations {
		err := client.CreateVariation(flag, variation)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to create flag variations in Optimizely: %+v", err),
			})

			return diags
		}
	}

	err = client.CreateRuleset(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create ruleset in Optimizely: %+v", err),
		})

		return diags
	}

	err = client.EnableRuleset(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to enable ruleset in Optimizely: %+v", err),
		})

		return diags
	}

	d.SetId(strconv.FormatInt(featResp.ID, 10))
	// return resourceFeatureRead(ctx, d, m)
	return diags
}

func resourceFeatureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(FlagClient)

	flag := parseFlag(d)

	flagResp, err := client.GetFlag(flag.ProjectId, flag.Key)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed fetch flag from Optimizely: %+v", err),
		})

		return diags
	}

	flagResp.Variations, err = client.GetVariation(flag.ProjectId, flag.Key)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed fetch flag from Optimizely; failed to fetch variations: %+v", err),
		})

		return diags
	}

	flagResp.Environments, err = client.GetRuleset(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed fetch flag from Optimizely; failed to fetch environment rules: %+v", err),
		})

		return diags
	}

	flagRespJson, _ := json.Marshal(flagResp)
	fmt.Printf("Read Flag: %s", flagRespJson)

	return diags
}

func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(FlagClient)

	flag := parseFlag(d)

	fmt.Printf("DELETE %s \n", flag.Key)

	err := client.DisableRuleset(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to disable ruleset while deleting flag in Optimizely: %+v", err),
		})

		return diags
	}

	err = client.DeleteFlag(flag.ProjectId, flag.Key)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to delete flag in Optimizely: %+v", err),
		})

		return diags
	}

	return diags
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(FlagClient)

	flag := parseFlag(d)

	err := client.UpdateRuleset(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create ruleset in Optimizely: %+v", err),
		})

		return diags
	}

	err = client.EnableRuleset(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to enable ruleset in Optimizely: %+v", err),
		})

		return diags
	}

	// return resourceFeatureRead(ctx, d, m)
	return diags
}
