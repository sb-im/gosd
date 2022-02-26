package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	headerApiKey = "X-API-Key"
)

type Client struct {
	apiKey   string
	endpoint string
}

func NewClient(endpoint string, apiKey string) *Client {
	return &Client{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

func (c *Client) do(method string, url string, body interface{}) (*http.Response, error) {
	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	return c.request(req)
}

func (c *Client) request(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headerApiKey, c.apiKey)
	return (&http.Client{}).Do(req)
}
