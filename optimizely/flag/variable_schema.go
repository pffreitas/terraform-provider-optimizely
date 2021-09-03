package flag

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type VariableSchema struct {
	DefaultValue string `json:"default_value"`
	Key          string `json:"key"`
	Type         string `json:"type"`
}

func parseVariableSchema(d *schema.ResourceData) map[string]VariableSchema {
	variableSchemaByKey := make(map[string]VariableSchema)
	variableSchemaList := d.Get("variable_schema").([]interface{})

	for _, variable := range variableSchemaList {
		vars := variable.(map[string]interface{})["variable"]
		for _, v := range vars.([]interface{}) {
			vMap := v.(map[string]interface{})

			key := vMap["key"].(string)
			variableSchema := VariableSchema{
				Key:          key,
				DefaultValue: vMap["default_value"].(string),
				Type:         vMap["type"].(string),
			}
			variableSchemaByKey[key] = variableSchema
		}
	}

	return variableSchemaByKey
}
