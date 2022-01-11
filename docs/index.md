# Optimizely Provider

Optimizely Terraform Provider allows you to manager Optimizely resources such as Flags and Audiences. 

## Example Usage

```hcl 

terraform {
  required_providers {
    optimizely = {
      source = "pffreitas/optimizely"
      version = "0.0.19"
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

## Argument Reference

* `api_host` - (Required) List arguments this resource takes.
* `api_token` - (Required) List arguments this resource takes.