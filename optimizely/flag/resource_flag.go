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

	variablesSchema := parseVariableSchema(d)
	variations := parseVariation(d)
	envs := parseEnvironment(d)

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
