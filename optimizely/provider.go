package optimizely

import (
	"context"

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

	optimizelyClient := client.OptimizelyClient{
		Address: address,
		Token:   token,
	}

	return optimizelyClient, diags
}
