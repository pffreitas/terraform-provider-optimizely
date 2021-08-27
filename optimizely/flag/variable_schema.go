package flag

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type VariableSchema struct {
	Archived     bool   `json:"archived"`
	DefaultValue string `json:"default_value"`
	Key          string `json:"key"`
	Type         string `json:"type"`
}

func parseVariableSchema(d *schema.ResourceData) []VariableSchema {
	var variablesSchema []VariableSchema
	variableSchema := d.Get("variable_schema").([]interface{})

	for _, variable := range variableSchema {
		vars := variable.(map[string]interface{})["variable"]
		for _, v := range vars.([]interface{}) {
			vMap := v.(map[string]interface{})
			vSchema := VariableSchema{
				Key:          vMap["key"].(string),
				DefaultValue: vMap["default_value"].(string),
				Type:         vMap["type"].(string),
				Archived:     vMap["archived"].(bool),
			}
			variablesSchema = append(variablesSchema, vSchema)
		}
	}

	return variablesSchema
}
