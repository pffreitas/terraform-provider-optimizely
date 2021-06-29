package optimizely

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Feature struct {
	ID           int64                         `json:"id"`
	ProjectId    int64                         `json:"project_id"`
	Name         string                        `json:"name"`
	Description  string                        `json:"description"`
	Key          string                        `json:"key"`
	Archived     bool                          `json:"archived"`
	Variables    []VariableSchema              `json:"variables"`
	Environments map[string]FeatureEnvironment `json:"environments"`
}

type FeatureEnvironment struct {
	RolloutRules []RolloutRule `json:"rollout_rules"`
}

type RolloutRule struct {
	AudienceConditions string `json:"audience_conditions"`
	Enabled            bool   `json:"enabled"`
	PercentageIncluded int    `json:"percentage_included"`
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

func resourceFeature() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
									"enabled": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"variables": {
										Type:     schema.TypeMap,
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
	client := m.(OptimizelyClient)

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

				audienceConditionsJson, _ := json.Marshal(audConditions)

				rollout := rMap["enabled"].(int)
				rolloutRule := RolloutRule{
					AudienceConditions: string(audienceConditionsJson),
					Enabled:            rollout > 0,
					PercentageIncluded: rollout,
				}

				fmt.Println(rolloutRule)
				if featureEnvironment, ok := envs[env.(string)]; ok {
					// featureEnvironment.RolloutRules = append(featureEnvironment.RolloutRules, rolloutRule)
					envs[env.(string)] = featureEnvironment
				}

				if featureEnvironment, ok := envs[env.(string)]; !ok {
					featureEnvironment = FeatureEnvironment{
						RolloutRules: []RolloutRule{},
					}
					envs[env.(string)] = featureEnvironment
				}
			}

		}
	}

	for key, env := range envs {

		env.RolloutRules = append(env.RolloutRules, RolloutRule{
			AudienceConditions: "everyone",
			Enabled:            false,
			PercentageIncluded: 0,
		})
		envs[key] = env
	}

	feat := Feature{
		ProjectId:    client.ProjectId,
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Key:          d.Get("key").(string),
		Archived:     false,
		Variables:    variablesSchema,
		Environments: envs,
	}

	j, _ := json.Marshal(feat)
	fmt.Printf("\n ---- %s \n", j)

	featResp, err := client.CreateFeature(feat)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create Feature in Optimizely: %+v", err),
		})

		return diags
	}

	d.SetId(strconv.FormatInt(featResp.ID, 10))
	// return resourceFeatureRead(ctx, d, m)
	return diags
}

func resourceFeatureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(OptimizelyClient)
	aud, err := client.GetAudience(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create Audience in Optimizely: %+v", err),
		})

		return diags
	}

	d.SetId(strconv.FormatInt(aud.ID, 10))
	d.Set("name", aud.Name)
	d.Set("description", aud.Description)
	d.Set("conditions", aud.Conditions)

	return diags
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(OptimizelyClient)

	audId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse Audicence ID: %s,  %+v", d.Id(), err),
		})

		return diags
	}

	aud := Audience{
		ProjectId:   client.ProjectId,
		ID:          audId,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Conditions:  d.Get("conditions").(string),
	}

	_, err = client.UpdateAudience(aud)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create Audience in Optimizely: %+v", err),
		})

		return diags
	}

	return resourceAudienceRead(ctx, d, m)
}

func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(OptimizelyClient)

	_, err := client.ArchiveAudience(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create Audience in Optimizely: %+v", err),
		})

		return diags
	}

	return resourceAudienceRead(ctx, d, m)
}
