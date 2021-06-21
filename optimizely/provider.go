package optimizely

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"optimizely_feature":  resourceFeature(),
			"optimizely_audience": resourceAudience(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"optimizely_environment": dataSourceEnvironment(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	address := d.Get("host").(string)
	token := d.Get("token").(string)
	projectId := d.Get("project_id").(string)
	projectIdInt64, err := strconv.ParseInt(projectId, 10, 64)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse project_id: %s; %+v", projectId, err),
		})

		return nil, diags
	}

	optimizelyClient := OptimizelyClient{
		Address:   address,
		Token:     token,
		ProjectId: projectIdInt64,
	}

	return optimizelyClient, diags
}
