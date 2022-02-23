package client

import "net/http"

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

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headerApiKey, c.apiKey)
	return (&http.Client{}).Do(req)
}
