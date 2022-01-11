# Audience Resource

Description of what this resource does, with links to official
app/service documentation.

## Example Usage

```hcl
resource "optimizely_audience" "country_us" {
    name = "COUNTRY_US_TERRAFORM"
    conditions = "[\"and\", {\"type\": \"custom_attribute\", \"name\": \"COUNTRY\", \"value\": \"us\"}]"
}
```

## Argument Reference

* `name` - (Optional/Required) List arguments this resource takes.
* `conditions` - (Optional/Required) List arguments this resource takes.

## Attribute Reference

* `id` - List attributes that this resource exports.