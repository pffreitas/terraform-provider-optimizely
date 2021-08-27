package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	req, err := c.newHttpRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags", feat.ProjectId), bytes.NewBuffer(postBody))
	if err != nil {
		return feat, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("\n\n Create flag - %s -- %+v \n\n", postBody, err)
		return feat, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return feat, err
	}

	fmt.Printf("\n\n Create flag OK- %s -- %+v \n\n", postBody, err)

	var featureResp flag.Flag
	json.Unmarshal(body, &featureResp)

	return featureResp, nil
}

func (c OptimizelyClient) DeleteFlag(projectId int, flagKey string) error {

	req, err := c.newHttpRequest("DELETE", fmt.Sprintf("flags/v1/projects/%d/flags/%s", projectId, flagKey), nil)
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

	return err
}