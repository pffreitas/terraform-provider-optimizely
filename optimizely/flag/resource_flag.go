package flag

import (
	"context"
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
	Variables    []VariableSchema              `json:"variables"`
	Variations   []Variation                   `json:"variations"`
	Environments map[string]FeatureEnvironment `json:"environments"`
}

type FeatureEnvironment struct {
	RolloutRules []RolloutRule `json:"rollout_rules"`
}

type RolloutRule struct {
	Key                string      `json:"key"`
	AudienceConditions []Condition `json:"audience_conditions"`
	Enabled            bool        `json:"enabled"`
	PercentageIncluded int         `json:"percentage_included"`
	Deliver            string      `json:"deliver"`
}

type Condition interface{}

type AudienceCondition struct {
	AudienceID int64 `json:"audience_id"`
}

type VariableSchema struct {
	Archived     bool   `json:"archived"`
	DefaultValue string `json:"default_value"`
	Key          string `json:"key"`
	Type         string `json:"type"`
}

type Variation struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Variables   map[string]interface{} `json:"variables"`
}

func ResourceFeature() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Project ID",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Feature key, also acts as it's unique ID",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human readable name",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of this feature",
			},
			"variable_schema": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"archived": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"default_value": {
										Type:     schema.TypeString,
										Required: true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variation": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"variables": {
										Type:     schema.TypeMap,
										Optional: true,
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
		UpdateContext: resourceFeatureUpdate,
		DeleteContext: resourceFeatureDelete,
	}
}

func resourceFeatureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(FlagClient)

	var variablesSchema []VariableSchema
	variableSchema := d.Get("variable_schema").([]interface{})
	for _, variable := range variableSchema {
		vars := variable.(map[string]interface{})["variable"]
		for _, v := range vars.([]interface{}) {
			vMap := v.(map[string]interface{})
			vSchema := VariableSchema{
				Key:          vMap["key"].(string),
				DefaultValue: vMap["default_value"].(string),
				Type:         vMap["type"].(string),
				Archived:     vMap["archived"].(bool),
			}
			variablesSchema = append(variablesSchema, vSchema)
		}
	}

	var variations []Variation
	for _, variationMap := range d.Get("variations").([]interface{}) {
		vars := variationMap.(map[string]interface{})["variation"]
		for _, v := range vars.([]interface{}) {
			vMap := v.(map[string]interface{})
			vSchema := Variation{
				Key:         vMap["key"].(string),
				Name:        vMap["name"].(string),
				Description: vMap["description"].(string),
				Variables:   vMap["variables"].(map[string]interface{}),
			}
			variations = append(variations, vSchema)
		}
	}

	var envs = make(map[string]FeatureEnvironment)
	for _, rules := range d.Get("rules").([]interface{}) {
		rule := rules.(map[string]interface{})["rule"]
		for _, r := range rule.([]interface{}) {

			rMap := r.(map[string]interface{})
			environments := rMap["environments"].([]interface{})
			for _, env := range environments {
				audConditions := []Condition{"and"}

				for _, audId := range rMap["audience"].([]interface{}) {
					audIdInt, _ := strconv.ParseInt(audId.(string), 10, 64)
					audConditions = append(audConditions, AudienceCondition{AudienceID: audIdInt})
				}

				rollout := rMap["percentage_included"].(int)
				rolloutRule := RolloutRule{
					Key:                rMap["key"].(string),
					AudienceConditions: audConditions,
					Enabled:            rollout > 0,
					PercentageIncluded: rollout * 100,
					Deliver:            rMap["deliver"].(string),
				}

				if featureEnvironment, ok := envs[env.(string)]; ok {
					featureEnvironment.RolloutRules = append(featureEnvironment.RolloutRules, rolloutRule)
					envs[env.(string)] = featureEnvironment
				}

				if featureEnvironment, ok := envs[env.(string)]; !ok {
					featureEnvironment = FeatureEnvironment{
						RolloutRules: []RolloutRule{rolloutRule},
					}
					envs[env.(string)] = featureEnvironment
				}
			}

		}
	}

	flag := Flag{
		ProjectId:    d.Get("project").(int),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Key:          d.Get("key").(string),
		Archived:     false,
		Variables:    variablesSchema,
		Variations:   variations,
		Environments: envs,
	}

	featResp, err := client.CreateFlag(flag)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create flag in Optimizely: %+v", err),
		})

		return diags
	}

	for _, variation := range variations {
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
	return nil
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
