package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pffreitas/optimizely-terraform-provider/optimizely/audience"
)

func (c OptimizelyClient) CreateAudience(aud audience.Audience) (audience.Audience, error) {
	postBody, err := json.Marshal(aud)
	if err != nil {
		return aud, err
	}

	respBody, err := c.sendHttpRequest("POST", "v2/audiences", bytes.NewBuffer(postBody))
	if err != nil {
		return aud, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(respBody, &audienceResp)

	return audienceResp, nil
}

func (c OptimizelyClient) GetAudience(audId string) (audience.Audience, error) {

	respBody, err := c.sendHttpRequest("GET", fmt.Sprintf("v2/audiences/%s", audId), nil)
	if err != nil {
		return audience.Audience{}, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(respBody, &audienceResp)

	return audienceResp, nil
}

func (c OptimizelyClient) ArchiveAudience(audId string) (audience.Audience, error) {
	postBody, err := json.Marshal(map[string]interface{}{
		"archived": true,
	})
	if err != nil {
		return audience.Audience{}, err
	}

	respBody, err := c.sendHttpRequest("PATCH", fmt.Sprintf("v2/audiences/%s", audId), bytes.NewBuffer(postBody))
	if err != nil {
		return audience.Audience{}, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(respBody, &audienceResp)

	return audienceResp, nil
}

func (c OptimizelyClient) UpdateAudience(aud audience.Audience) (audience.Audience, error) {
	postBody, err := json.Marshal(aud)
	if err != nil {
		return audience.Audience{}, err
	}

	respBody, err := c.sendHttpRequest("PATCH", fmt.Sprintf("v2/audiences/%d", aud.ID), bytes.NewBuffer(postBody))
	if err != nil {
		return audience.Audience{}, err
	}

	var audienceResp audience.Audience
	json.Unmarshal(respBody, &audienceResp)

	return audienceResp, nil
}
