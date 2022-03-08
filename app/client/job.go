package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) JobCreate(task *model.Task) error {
	res, err := c.do(http.MethodPost, c.endpoint+fmt.Sprintf("/tasks/%d/jobs", task.ID), task)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return json.NewDecoder(res.Body).Decode(task)
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}
