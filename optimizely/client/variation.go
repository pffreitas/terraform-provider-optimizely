package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	req, err := c.newHttpRequest("POST", fmt.Sprintf("flags/v1/projects/%d/flags/%s/variations", flag.ProjectId, flag.Key), bytes.NewBuffer(postBody))
	if err != nil {
		return err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("\n\n Create variation - %s -- %+v \n\n", postBody, err)
		return err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("\n\n Create variation OK - reqBody: %s -- respBody: %s --  err: %+v \n\n", postBody, respBody, err)

	return nil
}
