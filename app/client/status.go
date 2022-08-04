package client

import (
	"net/http"
)

func (c *Client) ServerStatus() error {
	_, err := c.do(http.MethodGet, c.endpoint+"/status", nil)
	return err
}
