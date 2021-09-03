package client

import (
	"fmt"
	"io"
	"net/http"
)

type OptimizelyClient struct {
	Address string
	Token   string
}

func (c OptimizelyClient) newHttpRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.Address, url), body)
	req.Header.Set("Content-type", "application/json")
	c.configureToken(req)

	return req, err
}

func (c OptimizelyClient) configureToken(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
}

func (c OptimizelyClient) newEmptyRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.Address, url), nil)
	c.configureToken(req)

	return req, err
}

func (c OptimizelyClient) isOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
