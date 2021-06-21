package optimizely

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Audience struct {
	ID          int64  `json:"id"`
	ProjectId   int64  `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Conditions  string `json:"conditions"`
	Archived    bool   `json:"archived"`
}

func resourceAudience() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
	client := m.(OptimizelyClient)

	aud := Audience{
		ProjectId:   client.ProjectId,
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

func resourceAudienceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceAudienceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
