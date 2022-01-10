### Terraform provider for Optimizely




### How to use 


#### Configure Optimizely Provider

```
terraform {
  required_providers {
    optimizely = {
      source = "pffreitas/optimizely"
      version = "0.0.18"
    }
  }
}

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
```

#### Configure Environments and Projects

``` 
data "optimizely_project" "bees_test_cac" { 
	id = 20410805626
}

data "optimizely_environment" "dev" {
  key = "dev"
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
```

#### Create Audiences

```
resource "optimizely_audience" "country_ec" {
    name = "COUNTRY_EC_TERRAFORM"
    conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"ec\"}]"
}
```