package flag

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type Variation struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Variables   map[string]interface{} `json:"variables"`
}

func parseVariation(d *schema.ResourceData) []Variation {
	var variations []Variation
	for _, variationMap := range d.Get("variations").([]interface{}) {
		vars := variationMap.(map[string]interface{})["variation"]
		for _, v := range vars.([]interface{}) {
			vMap := v.(map[string]interface{})
			vSchema := Variation{
				Key:         vMap["key"].(string),
				Name:        vMap["name"].(string),
				Description: vMap["description"].(string),
				Variables:   vMap["variables"].(map[string]interface{}),
			}
			variations = append(variations, vSchema)
		}
	}

	return variations
}
