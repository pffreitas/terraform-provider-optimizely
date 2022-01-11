# Flag Resource

Manages Optimizely Flags

## Example Usage

```hcl
resource "optimizely_feature" "out-of-stock" {
  project     = data.optimizely_project.bees_test_cac.id
  name        = "Out of stock"
  description = "Out of stock"
  key         = "oos"

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
    variation {
      key         = "blackButtonOnTheLeft"
      name        = "blackButtonOnTheLeft"
      description = "blackButtonOnTheLeft"
      variables = {
        buttonPosition = "left"
        buttonColor    = "black"
      }
    }
  }

  rules {
    rule {
      key                 = "us-right"
      environments        = [data.optimizely_environment.sit.id]
      audience            = [optimizely_audience.country_us.id]
      percentage_included = 50
      deliver             = "blackButtonOnTheRight"
    }

    rule {
      key                 = "us-left"
      environments        = [data.optimizely_environment.sit.id]
      audience            = [optimizely_audience.country_us.id]
      percentage_included = 50
      deliver             = "blackButtonOnTheLeft"
    }

    rule {
      key                 = "br-on"
      environments        = [data.optimizely_environment.sit.id]
      audience            = [optimizely_audience.country_br.id]
      percentage_included = 75
      deliver             = "on"
    }

    rule {
      key                 = "br-off"
      environments        = [data.optimizely_environment.sit.id]
      audience            = [optimizely_audience.country_br.id]
      percentage_included = 25
      deliver             = "off"
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

    # rule {
    #   key                 = "br-dev"
    #   environments        = [data.optimizely_environment.dev.id]
    #   audience            = [optimizely_audience.country_br.id]
    #   percentage_included = 100
    #   deliver             = "on"
    # }
  }
}

```

## Argument Reference


## Attribute Reference

* `id` - Flag Id.