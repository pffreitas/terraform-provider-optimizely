package flag

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FeatureEnvironment struct {
	RolloutRules []RolloutRule `json:"rollout_rules"`
}

type RolloutRule struct {
	Key                string      `json:"key"`
	AudienceConditions []Condition `json:"audience_conditions"`
	PercentageIncluded int         `json:"percentage_included"`
	Deliver            string      `json:"deliver"`
}

type Condition interface{}

type AudienceCondition struct {
	AudienceID int64 `json:"audience_id"`
}

func parseEnvironment(d *schema.ResourceData) map[string]FeatureEnvironment {
	var envs = make(map[string]FeatureEnvironment)

	for _, rules := range d.Get("rules").([]interface{}) {
		rule := rules.(map[string]interface{})["rule"]
		for _, r := range rule.([]interface{}) {

			rMap := r.(map[string]interface{})
			environments := rMap["environments"].([]interface{})
			for _, env := range environments {
				audConditions := []Condition{"and"}

				for _, audId := range rMap["audience"].([]interface{}) {
					audIdInt, _ := strconv.ParseInt(audId.(string), 10, 64)
					audConditions = append(audConditions, AudienceCondition{AudienceID: audIdInt})
				}

				rollout := rMap["percentage_included"].(int)
				rolloutRule := RolloutRule{
					Key:                rMap["key"].(string),
					AudienceConditions: audConditions,
					PercentageIncluded: rollout * 100,
					Deliver:            rMap["deliver"].(string),
				}

				if featureEnvironment, ok := envs[env.(string)]; ok {
					featureEnvironment.RolloutRules = append(featureEnvironment.RolloutRules, rolloutRule)
					envs[env.(string)] = featureEnvironment
				}

				if featureEnvironment, ok := envs[env.(string)]; !ok {
					featureEnvironment = FeatureEnvironment{
						RolloutRules: []RolloutRule{rolloutRule},
					}
					envs[env.(string)] = featureEnvironment
				}
			}

		}
	}

	return envs
}
