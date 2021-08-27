package optimizely

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pffreitas/optimizely-terraform-provider/optimizely/audience"
	"github.com/pffreitas/optimizely-terraform-provider/optimizely/client"
	"github.com/pffreitas/optimizely-terraform-provider/optimizely/environment"
	"github.com/pffreitas/optimizely-terraform-provider/optimizely/flag"
	"github.com/pffreitas/optimizely-terraform-provider/optimizely/project"
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
			"optimizely_feature":  flag.ResourceFeature(),
			"optimizely_audience": audience.ResourceAudience(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"optimizely_environment": environment.DataSourceEnvironment(),
			"optimizely_project":     project.DataSourceProject(),
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

	optimizelyClient := client.OptimizelyClient{
		Address:   address,
		Token:     token,
		ProjectId: projectIdInt64,
	}

	return optimizelyClient, diags
}
