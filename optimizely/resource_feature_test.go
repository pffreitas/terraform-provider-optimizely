package optimizely

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testFeatureConfigBasic() string {
	return `
	provider "optimizely" {
		host  = "https://api.optimizely.com/v2"
		token = "2:myr1dVQxw203jqcj-vr4Sxr1PNAfu2FzPrwauwA_vPcc9HHMB1GY"
		project_id = "19036502365"
	}
	
	data "optimizely_environment" "sit" {
		key = "TOGGLES_SIT"
	}
	data "optimizely_environment" "uat" {
		key = "TOGGLES_UAT"
	}
	data "optimizely_environment" "prod" {
		key = "TOGGLES_PROD"
	}
	
	resource "optimizely_feature" "dynamic_forms_terraform" {
		name        = "Customer Support - Dynamic Forms - Terraform"
		description = "Customer Support - Dynamic Forms - Terraform"
		key         = "dynamic_forms_enabled"
	  
		variable_schema {
		  variable {
			key         = "buttonPosition"
			type         = "string"
			default_value = "left"
		  }

		  variable {
			key         = "buttonColor"
			type         = "string"
			default_value = "black"
		  }
		}
	  
		rules {
		  rule {
			environments = [data.optimizely_environment.sit.id]
			audience     = ["20381209628"]

			enabled = 1000
			variables = {
			  buttonPosition = "left"
			}
		  }

		  rule {
			environments = [data.optimizely_environment.uat.id, data.optimizely_environment.prod.id]
			audience     = ["20381209628"]

			enabled = 100
			variables = {
			  buttonPosition = "right"
			}
		  }
		}
	  }	 
	`
}

func TestFeatureBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFeatureConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHashicupsOrderExists("optimizely_feature.dynamic_forms"),
				),
			},
		},
	})
}
