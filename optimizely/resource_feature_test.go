package optimizely

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testFeatureConfigBasic() string {
	return `
	provider "optimizely" {
		host  = "https://api.optimizely.com"
		token = "2:myr1dVQxw203jqcj-vr4Sxr1PNAfu2FzPrwauwA_vPcc9HHMB1GY"
		project_id = "20410805626"
	}
	
	data "optimizely_environment" "sit" {
		key = "sit"
	}
	data "optimizely_environment" "uat" {
		key = "uat"
	}
	data "optimizely_environment" "prod" {
		key = "prod"
	}

	resource "optimizely_audience" "country_us" {
		name = "COUNTRY_US_TERRAFORM"
		conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"us\"}]"
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
			key 		 = "us"
			environments = [data.optimizely_environment.sit.id]
			audience     = [optimizely_audience.country_us.id]

			enabled = 1000
			variables = {
			  buttonPosition = "left"
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
