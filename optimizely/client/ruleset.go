package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pffreitas/optimizely-terraform-provider/optimizely/flag"
)

type RulesetType string

const TargetedDelivery RulesetType = "targeted_delivery"
const ABTesting RulesetType = "a/b"

type RulesetVariationVariable struct {
	Value string `json:"value"`
}

type RulesetVariation struct {
	Key                string                              `json:"key"`
	Name               string                              `json:"name"`
	PercentageIncluded int                                 `json:"percentage_included"`
	Variables          map[string]RulesetVariationVariable `json:"variables"`
}

type AudicenceCondition struct {
}

type OptimizelyRuleset struct {
	Key                 string                      `json:"key"`
	Name                string                      `json:"name"`
	Type                RulesetType                 `json:"type"`
	PercentageIncluded  int                         `json:"percentage_included"`
	Variations          map[string]RulesetVariation `json:"variations"`
	AudicenceConditions []flag.Condition            `json:"audience_conditions"`
}

type Operation string

type OptimizelyOp struct {
	Op    Operation   `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type OptimizelyRules struct {
	Rules          []OptimizelyRuleset `json:"rules"`
	RulePriorities []string            `json:"rule_priorities"`
}

func (c OptimizelyClient) CreateRuleset(flag flag.Flag) error {

	for env, flagEnv := range flag.Environments {
		ops := []OptimizelyOp{}

		for i, rule := range flagEnv.RolloutRules {

			ruleset := OptimizelyRuleset{
				Key:  rule.Key,
				Name: rule.Key,
				Type: TargetedDelivery,
				Variations: map[string]RulesetVariation{
					rule.Deliver: {
						Key:                rule.Deliver,
						PercentageIncluded: rule.PercentageIncluded,
					},
				},
				AudicenceConditions: rule.AudienceConditions,
				PercentageIncluded:  rule.PercentageIncluded,
			}

			ops = append(ops, OptimizelyOp{
				Op:    "add",
				Path:  fmt.Sprintf("/rules/%s", rule.Key),
				Value: ruleset,
			})

			ops = append(ops, OptimizelyOp{
				Op:    "add",
				Path:  fmt.Sprintf("/rule_priorities/%d", i),
				Value: rule.Key,
			})

		}

		postBody, err := json.Marshal(ops)
		if err != nil {
			return err
		}

		req, err := c.newHttpRequest("PATCH", fmt.Sprintf("flags/v1/projects/%d/flags/%s/environments/%s/ruleset", flag.ProjectId, flag.Key, env), bytes.NewBuffer(postBody))
		if err != nil {
			return err
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

	}
	return nil
}

type getRulesetResponse struct {
	Rules map[string]OptimizelyRuleset `json:"rules"`
}

func (c OptimizelyClient) GetRuleset(flg flag.Flag) (map[string]flag.FeatureEnvironment, error) {
	flagEnvs := make(map[string]flag.FeatureEnvironment)

	for env := range flg.Environments {
		flagEnv := flag.FeatureEnvironment{}

		req, err := c.newEmptyRequest("GET", fmt.Sprintf("flags/v1/projects/%d/flags/%s/environments/%s/ruleset", flg.ProjectId, flg.Key, env))
		if err != nil {
			return flagEnvs, err
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return flagEnvs, err
		}

		defer resp.Body.Close()
		rulesetResponseBodyStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return flagEnvs, err
		}

		var rulesetResponseBody getRulesetResponse
		err = json.Unmarshal(rulesetResponseBodyStr, &rulesetResponseBody)
		if err != nil {
			return flagEnvs, err
		}

		for _, ruleset := range rulesetResponseBody.Rules {

			deliver := ""
			for variationKey := range ruleset.Variations {
				deliver = variationKey
			}

			audienceConditions := []flag.Condition{}
			for _, aud := range ruleset.AudicenceConditions {
				if audienceConditionMap, ok := aud.(map[string]interface{}); ok {
					if audienceId, ok := audienceConditionMap["audience_id"].(float64); ok {
						audienceConditions = append(audienceConditions, flag.AudienceCondition{
							AudienceID: int64(audienceId),
						})
					}
				}
			}

			flagEnv.RolloutRules = append(flagEnv.RolloutRules, flag.RolloutRule{
				Key:                ruleset.Key,
				PercentageIncluded: ruleset.PercentageIncluded / 100,
				AudienceConditions: audienceConditions,
				Deliver:            deliver,
			})
		}

		flagEnvs[env] = flagEnv
	}

	return flagEnvs, nil

}

func (c OptimizelyClient) EnableRuleset(flag flag.Flag) error {

	for env := range flag.Environments {
		req, err := c.newEmptyRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags/%s/environments/%s/ruleset/enabled", flag.ProjectId, flag.Key, env))
		if err != nil {
			return err
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

	}

	return nil
}

func (c OptimizelyClient) DisableRuleset(flag flag.Flag) error {

	for env := range flag.Environments {
		req, err := c.newEmptyRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags/%s/environments/%s/ruleset/disabled", flag.ProjectId, flag.Key, env))
		if err != nil {
			return err
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

	}

	return nil
}
