package audience

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Audience struct {
	ID          int64  `json:"id"`
	ProjectId   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Conditions  string `json:"conditions"`
	Archived    bool   `json:"archived"`
}

func ResourceAudience() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Project ID",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Audience",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Audience",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A short description of the Audience",
			},
			"conditions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A string defining the targeting rules for an Audience",
			},
		},
		CreateContext: resourceAudienceCreate,
		ReadContext:   resourceAudienceRead,
		UpdateContext: resourceAudienceUpdate,
		DeleteContext: resourceAudienceDelete,
	}
}

func resourceAudienceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(AudienceClient)

	aud := Audience{
		ProjectId:   d.Get("project").(int),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Conditions:  d.Get("conditions").(string),
	}

	audResp, err := client.CreateAudience(aud)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create Audience in Optimizely: %+v", err),
		})

		return diags
	}

	d.SetId(strconv.FormatInt(audResp.ID, 10))
	return resourceAudienceRead(ctx, d, m)
}

func resourceAudienceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(AudienceClient)
	aud, err := client.GetAudience(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create Audience in Optimizely: %+v", err),
		})

		return diags
	}

	compactConditions := new(bytes.Buffer)
	err = json.Compact(compactConditions, []byte(aud.Conditions))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to compact conditions: %s --- %+v", aud.Conditions, err),
		})

		return diags
	}

	d.SetId(strconv.FormatInt(aud.ID, 10))
	d.Set("name", aud.Name)
	d.Set("description", aud.Description)
	d.Set("conditions", compactConditions.String())

	return diags
}

func resourceAudienceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(AudienceClient)

	audId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse Audicence ID: %s,  %+v", d.Id(), err),
		})

		return diags
	}

	aud := Audience{
		ProjectId:   d.Get("project").(int),
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

func resourceAudienceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(AudienceClient)

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
