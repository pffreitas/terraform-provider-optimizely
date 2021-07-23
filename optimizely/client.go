package optimizely

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type OptimizelyClient struct {
	ProjectId int64
	Address   string
	Token     string
}

func (c *OptimizelyClient) newHttpRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.Address, url), body)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-type", "application/json")

	return req, err
}

func (c *OptimizelyClient) CreateAudience(aud Audience) (Audience, error) {
	postBody, err := json.Marshal(aud)
	if err != nil {
		return aud, err
	}

	req, err := c.newHttpRequest("POST", "v2/audiences", bytes.NewBuffer(postBody))
	if err != nil {
		return aud, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return aud, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return aud, err
	}

	var audienceResp Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

func (c *OptimizelyClient) GetAudience(audId string) (Audience, error) {

	req, err := c.newHttpRequest("GET", fmt.Sprintf("v2/audiences/%s", audId), nil)
	if err != nil {
		return Audience{}, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return Audience{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Audience{}, err
	}

	var audienceResp Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

func (c *OptimizelyClient) ArchiveAudience(audId string) (Audience, error) {
	postBody, err := json.Marshal(map[string]interface{}{
		"archived": true,
	})
	if err != nil {
		return Audience{}, err
	}

	req, err := c.newHttpRequest("PATCH", fmt.Sprintf("v2/audiences/%s", audId), bytes.NewBuffer(postBody))
	if err != nil {
		return Audience{}, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return Audience{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Audience{}, err
	}

	var audienceResp Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

func (c *OptimizelyClient) UpdateAudience(aud Audience) (Audience, error) {
	postBody, err := json.Marshal(aud)
	if err != nil {
		return Audience{}, err
	}

	req, err := c.newHttpRequest("PATCH", fmt.Sprintf("v2/audiences/%d", aud.ID), bytes.NewBuffer(postBody))
	if err != nil {
		return Audience{}, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return Audience{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Audience{}, err
	}

	var audienceResp Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

type OptimizelyFlag struct {
	Key                string                                      `json:"key"`
	Name               string                                      `json:"name"`
	Description        string                                      `json:"description"`
	VariableDefinitons map[string]OptimizelyFlagVariableDefinition `json:"variable_definitions"`
}

type OptimizelyFlagVariableDefinition struct {
	Key          string `json:"key"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`
}

func (c *OptimizelyClient) CreateFeature(feat Feature) (Feature, error) {

	var variableDefinitions = make(map[string]OptimizelyFlagVariableDefinition)

	for _, variable := range feat.Variables {
		variableDefinitions[variable.Key] = OptimizelyFlagVariableDefinition{
			Key:          variable.Key,
			Type:         variable.Type,
			DefaultValue: variable.DefaultValue,
		}
	}

	optimizelyFlag := OptimizelyFlag{
		Key:                feat.Key,
		Name:               feat.Name,
		Description:        feat.Description,
		VariableDefinitons: variableDefinitions,
	}

	postBody, err := json.Marshal(optimizelyFlag)
	if err != nil {
		return feat, err
	}

	fmt.Printf(" >>>> flag: %s \n", postBody)

	req, err := c.newHttpRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags", feat.ProjectId), bytes.NewBuffer(postBody))
	if err != nil {
		return feat, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return feat, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return feat, err
	}

	fmt.Printf("\n\n >>>>>>> %s \n\n", body)
	var featureResp Feature
	json.Unmarshal(body, &featureResp)

	return featureResp, nil
}

type RulesetType string

const TargetedDelivery RulesetType = "targeted_delivery"
const ABTesting RulesetType = "a/b"

type RulesetVariation struct {
	Key                string `json:"key"`
	Name               string `json:"name"`
	PercentageIncluded int    `json:"percentage_included"`
}

type AudicenceCondition struct {
}

type OptimizelyRuleset struct {
	Key                 string                      `json:"key"`
	Name                string                      `json:"name"`
	Type                RulesetType                 `json:"type"`
	PercentageIncluded  int                         `json:"percentage_included"`
	Variations          map[string]RulesetVariation `json:"variations"`
	AudicenceConditions []Condition                 `json:"audience_conditions"`
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

func (c *OptimizelyClient) CreateRuleset(feat Feature) error {

	for env, featureEnv := range feat.Environments {
		ops := []OptimizelyOp{}

		for i, rule := range featureEnv.RolloutRules {

			ruleset := OptimizelyRuleset{
				Key:  rule.Key,
				Name: rule.Key,
				Type: TargetedDelivery,
				Variations: map[string]RulesetVariation{
					"on": {
						Key:                "on",
						Name:               "on",
						PercentageIncluded: 10000,
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

		fmt.Printf(" >>>> flag: %s \n", postBody)

		req, err := c.newHttpRequest("PATCH", fmt.Sprintf("flags/v1/projects/%d/flags/%s/environments/%s/ruleset", feat.ProjectId, feat.Key, env), bytes.NewBuffer(postBody))
		if err != nil {
			return err
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("\n\n >>>>>>> %s \n\n", body)

	}
	return nil
}
