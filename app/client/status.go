package client

import (
	"errors"
	"net/http"
)

func (c *Client) ServerStatus() error {
	res, err := c.do(http.MethodGet, c.endpoint+"/status", nil)
	if res.StatusCode != http.StatusOK {
		return errors.New("status error: " + res.Status)
	}
	return err
}
