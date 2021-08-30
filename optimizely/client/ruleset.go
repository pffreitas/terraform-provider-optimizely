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
