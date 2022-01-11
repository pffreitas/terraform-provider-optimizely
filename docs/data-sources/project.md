# Project Data Source

Provides a Optimizely Project as datasource

## Example Usage

```hcl
data "optimizely_project" "my_project" { 
	id = 20410805626
}
```

## Argument Reference

* `id` - (Required) Project Id

## Attribute Reference

* `id` - Project id on Optimizely