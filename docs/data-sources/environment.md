# Environment Data Source

Provides a Optimizely Environment as datasource

## Example Usage

```hcl
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

## Argument Reference

* `key` - (Required) Environment key

## Attribute Reference

* `id` - Environment id on Optimizely