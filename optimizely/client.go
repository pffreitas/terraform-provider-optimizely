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

	return req, err
}

func (c *OptimizelyClient) CreateAudience(aud Audience) (Audience, error) {
	postBody, err := json.Marshal(aud)
	if err != nil {
		return aud, err
	}

	req, err := c.newHttpRequest("POST", "audiences", bytes.NewBuffer(postBody))
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

	req, err := c.newHttpRequest("GET", fmt.Sprintf("audiences/%s", audId), nil)
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

	req, err := c.newHttpRequest("PATCH", fmt.Sprintf("audiences/%s", audId), bytes.NewBuffer(postBody))
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

	req, err := c.newHttpRequest("PATCH", fmt.Sprintf("audiences/%d", aud.ID), bytes.NewBuffer(postBody))
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

func (c *OptimizelyClient) CreateFeature(feat Feature) (Feature, error) {
	postBody, err := json.Marshal(feat)
	if err != nil {
		return feat, err
	}

	req, err := c.newHttpRequest("POST", "features", bytes.NewBuffer(postBody))
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
