package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pffreitas/optimizely-terraform-provider/optimizely/flag"
)

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

func (c OptimizelyClient) CreateFlag(feat flag.Flag) (flag.Flag, error) {

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

	respBody, err := c.sendHttpRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags", feat.ProjectId), bytes.NewBuffer(postBody))
	if err != nil {
		return feat, err
	}

	var featureResp flag.Flag
	json.Unmarshal(respBody, &featureResp)

	return featureResp, nil
}

func (c OptimizelyClient) GetFlag(projectId int, flagKey string) (flag.Flag, error) {
	respBody, err := c.sendHttpRequest("GET", fmt.Sprintf("flags/v1/projects/%d/flags/%s", projectId, flagKey), nil)
	if err != nil {
		return flag.Flag{}, err
	}

	var flagResp flag.Flag
	json.Unmarshal(respBody, &flagResp)

	return flagResp, nil
}

func (c OptimizelyClient) DeleteFlag(projectId int, flagKey string) error {
	_, err := c.sendHttpRequest("DELETE", fmt.Sprintf("flags/v1/projects/%d/flags/%s", projectId, flagKey), nil)
	return err
}
