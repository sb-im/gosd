package client

import (
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) NodeIndex() (nodes []model.Node, err error) {
	res, err := c.do(http.MethodGet, c.endpoint+"/nodes", nil)
	if err != nil {
		return nodes, err
	}

	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&nodes)
		return
	} else {
		err = &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return
	}
}

func (c *Client) NodeCreate(node interface{}) error {
	res, err := c.do(http.MethodPost, c.endpoint+"/nodes", node)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return json.NewDecoder(res.Body).Decode(node)
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}

func (c *Client) NodeShow(nodeId string) (node model.Node, err error) {
	res, err := c.do(http.MethodGet, c.endpoint+"/nodes/"+nodeId, nil)
	if err != nil {
		return node, err
	}

	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&node)
		return
	} else {
		err = &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return
	}
}