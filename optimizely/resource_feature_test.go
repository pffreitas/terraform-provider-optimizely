package optimizely

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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

		fmt.Printf("%+v --- %+v --- %+v", rs, ok, s.RootModule().Resources)

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

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
