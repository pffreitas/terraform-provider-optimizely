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

resource "optimizely_audience" "country_br" {
    name = "COUNTRY_BR_TERRAFORM"
    conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"br\"}]"
}

resource "optimizely_audience" "country_us" {
    name = "COUNTRY_US_TERRAFORM"
    conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"us\"}]"
}
```


#### Create Feature Flags 

```
resource "optimizely_feature" "dynamic_forms_terraform" {
  project     = data.optimizely_project.bees_test_cac.id
  name        = "Dynamic Forms"
  description = "Dynamic Forms"
  key         = "dynamic_forms"

  variable_schema {
    variable {
      key           = "buttonPosition"
      type          = "string"
      default_value = "left"
    }

    variable {
      key           = "buttonColor"
      type          = "string"
      default_value = "black"
    }
  }

  variations {
    variation {
      key         = "blackButtonOnTheRight"
      name        = "blackButtonOnTheRight"
      description = "blackButtonOnTheRight"
      variables = {
        buttonPosition = "right"
        buttonColor    = "black"
      }
    }
  }

  rules {
    rule {
      key                 = "us-50pct-blackButtonOnTheRight"
      environments        = [data.optimizely_environment.sit.id]
      audience            = [optimizely_audience.country_us.id]
      percentage_included = 50
      deliver             = "blackButtonOnTheRight"
    }

    rule {
      key                 = "br-75pct-on"
      environments        = [data.optimizely_environment.sit.id]
      audience            = [optimizely_audience.country_br.id]
      percentage_included = 75
      deliver             = "on"
    }

    rule {
      key                 = "br-uat"
      environments        = [data.optimizely_environment.uat.id]
      audience            = [optimizely_audience.country_br.id]
      percentage_included = 100
      deliver             = "on"
    }

    rule {
      key                 = "br-prod"
      environments        = [data.optimizely_environment.prod.id]
      audience            = [optimizely_audience.country_br.id]
      percentage_included = 100
      deliver             = "on"
    }

    rule {
      key                 = "br-dev"
      environments        = [data.optimizely_environment.dev.id]
      audience            = [optimizely_audience.country_br.id]
      percentage_included = 100
      deliver             = "on"
    }
  }
}

```