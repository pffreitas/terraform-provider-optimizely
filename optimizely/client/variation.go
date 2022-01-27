package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pffreitas/optimizely-terraform-provider/optimizely/flag"
)

type OptimizelyVariation struct {
	flag.Variation
	Variables map[string]OptimizelyVariationVariable `json:"variables"`
}

type OptimizelyVariationVariable struct {
	Value interface{} `json:"value"`
}

func (c OptimizelyClient) CreateVariation(flag flag.Flag, variation flag.Variation) error {

	optVariationVariables := make(map[string]OptimizelyVariationVariable)
	for key, value := range variation.Variables {
		optVariationVariables[key] = OptimizelyVariationVariable{Value: value}
	}

	optVariation := OptimizelyVariation{
		variation,
		optVariationVariables,
	}

	postBody, err := json.Marshal(optVariation)
	if err != nil {
		return err
	}

	_, err = c.sendHttpRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags/%s/variations", flag.ProjectId, flag.Key), bytes.NewBuffer(postBody))
	return err
}

type getVariationResponse struct {
	Items []flag.Variation `json:"items"`
}

func (c OptimizelyClient) GetVariation(projectId int, flagKey string) ([]flag.Variation, error) {
	var variations []flag.Variation
	respBody, err := c.sendHttpRequest("GET", fmt.Sprintf("flags/v1/projects/%d/flags/%s/variations", projectId, flagKey), nil)
	if err != nil {
		return variations, err
	}

	var getVariationResponse getVariationResponse
	err = json.Unmarshal(respBody, &getVariationResponse)
	if err != nil {
		return variations, err
	}

	return getVariationResponse.Items, nil
}
