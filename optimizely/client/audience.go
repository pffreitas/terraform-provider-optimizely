package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pffreitas/optimizely-terraform-provider/optimizely/audience"
)

func (c OptimizelyClient) CreateAudience(aud audience.Audience) (audience.Audience, error) {
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

	fmt.Printf("\n\n >>>>>>>> aud resp:  %s \n\n", body)
	var audienceResp audience.Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

func (c OptimizelyClient) GetAudience(audId string) (audience.Audience, error) {

	req, err := c.newHttpRequest("GET", fmt.Sprintf("v2/audiences/%s", audId), nil)
	if err != nil {
		return audience.Audience{}, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return audience.Audience{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return audience.Audience{}, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

func (c OptimizelyClient) ArchiveAudience(audId string) (audience.Audience, error) {
	postBody, err := json.Marshal(map[string]interface{}{
		"archived": true,
	})
	if err != nil {
		return audience.Audience{}, err
	}

	req, err := c.newHttpRequest("PATCH", fmt.Sprintf("v2/audiences/%s", audId), bytes.NewBuffer(postBody))
	if err != nil {
		return audience.Audience{}, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return audience.Audience{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return audience.Audience{}, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}

func (c OptimizelyClient) UpdateAudience(aud audience.Audience) (audience.Audience, error) {
	postBody, err := json.Marshal(aud)
	if err != nil {
		return audience.Audience{}, err
	}

	req, err := c.newHttpRequest("PATCH", fmt.Sprintf("v2/audiences/%d", aud.ID), bytes.NewBuffer(postBody))
	if err != nil {
		return audience.Audience{}, err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return audience.Audience{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return audience.Audience{}, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(body, &audienceResp)

	return audienceResp, nil
}
