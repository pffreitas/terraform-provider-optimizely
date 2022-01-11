# Audience Resource

Manages Optimizely Audiences 

## Example Usage

```hcl
resource "optimizely_audience" "country_us" {
    name = "COUNTRY_US_TERRAFORM"
    conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"us\"}]"
}
```

## Argument Reference

* `name` - (Required) Name.
* `conditions` - (Required) Conditions.

## Attribute Reference

* `id` - Audience Id.