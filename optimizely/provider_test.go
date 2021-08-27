package optimizely

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pffreitas/optimizely-terraform-provider/optimizely/client"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"optimizely": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccCheckHashicupsOrderConfigBasic() string {
	return `
	provider "optimizely" {
		host  = "https://api.optimizely.com/v2"
		token = "2:myr1dVQxw203jqcj-vr4Sxr1PNAfu2FzPrwauwA_vPcc9HHMB1GY"
		project_id = "19036502365"
	}
	
	data "optimizely_environment" "sit" {}
	
	resource "optimizely_audience" "country_ec" {
		name = "COUNTRY_EC_TERRAFORM"
		conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"ec\"}]"
	}
	`
}

func testAccCheckHashicupsOrderConfigBasic2() string {
	return `
	provider "optimizely" {
		host  = "https://api.optimizely.com/v2"
		token = "2:myr1dVQxw203jqcj-vr4Sxr1PNAfu2FzPrwauwA_vPcc9HHMB1GY"
		project_id = "19036502365"
	}
	
	data "optimizely_environment" "sit" {}
	
	resource "optimizely_audience" "country_ec" {
		name = "COUNTRY_EC_TERRAFORM_2"
		conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"ec\"}]"
	}
	`
}

func testAccCheckHashicupsOrderExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

		fmt.Printf("%+v --- %+v --- %+v \n", rs, ok, s.RootModule().Resources)

		return nil
	}
}

func TestAccHashicupsOrderBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHashicupsOrderConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHashicupsOrderExists("optimizely_audience.country_ec"),
				),
			},
			{
				Config: testAccCheckHashicupsOrderConfigBasic2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHashicupsOrderExists("optimizely_audience.country_ec"),
				),
			},
		},
	})
}
func testFlagConfigUpdate() string {
	return ""
}

func testFlagConfigBasic() (string, error) {
	data := map[string]interface{}{
		"AudienceName": strings.ToUpper(gofakeit.BS()),
		"FlagKey":      gofakeit.BS(),
	}

	tmpl, err := template.New("").Parse(`
	variable "api_host" {
		type = string
		default = "https://api.optimizely.com"
	}

	variable "api_token" {
		type = string
		sensitive = true
	}

	provider "optimizely" {
		host  = var.api_host
		token = var.api_token
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

	data "optimizely_project" "bees_test_cac" { 
		id = 20410805626
	}

	resource "optimizely_audience" "country_us" {
		project	= data.optimizely_project.bees_test_cac.id
		name = "{{.AudienceName}}-US"
		conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"us\"}]"
	}
	
	resource "optimizely_audience" "country_br" {
		project	= data.optimizely_project.bees_test_cac.id
		name = "{{.AudienceName}}-BR"
		conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"br\"}]"
	}

	resource "optimizely_feature" "dynamic_forms_terraform" {
		project	= data.optimizely_project.bees_test_cac.id
		name        = "{{.FlagKey}} - Terraform"
		description = "{{.FlagKey}} - Terraform"
		key         = "{{.FlagKey}}"
	  
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

		variations { 
			variation { 
				key = "blackButtonOnTheRight"
				name = "blackButtonOnTheRight"
				description = "blackButtonOnTheRight"
				variables = {
					buttonPosition = "right"
					buttonColor = "black"
				}
			}
		}
	  
		rules {
		  rule {
			key 		 = "us"
			environments = [data.optimizely_environment.sit.id]
			audience     = [optimizely_audience.country_us.id]
			percentage_included = 50
			deliver = "blackButtonOnTheRight"
		  }
		  
		  rule {
			key 		 = "br"
			environments = [data.optimizely_environment.sit.id]
			audience     = [optimizely_audience.country_br.id]
			percentage_included = 75
			deliver = "on"
		  }

		  rule {
			key 		 = "br-uat"
			environments = [data.optimizely_environment.uat.id]
			audience     = [optimizely_audience.country_br.id]
			percentage_included = 100
			deliver = "on"
		  }
		}
	  }
	`)

	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)

	return buf.String(), err
}

func testAccCheckHashicupsOrderDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(client.OptimizelyClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "optimizely_feature" {
			continue
		}

		flagId := rs.Primary.Attributes["key"]

		err := c.DeleteFlag(20410805626, flagId)
		if err != nil {
			fmt.Printf("\n delete flag error >>>>>>>>> %+v \n", err)
			return err
		}
	}

	return nil
}

func TestFlagBasic(t *testing.T) {
	hcl, _ := testFlagConfigBasic()
	hclUpdate := testFlagConfigUpdate()

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHashicupsOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: hcl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHashicupsOrderExists("optimizely_feature.dynamic_forms_terraform"),
				),
			},
			{
				Config: hclUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHashicupsOrderExists("optimizely_feature.dynamic_forms_terraform"),
				),
			},
		},
	})
}
