package client

import (
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) TaskIndex() (tasks []model.Task, err error) {
	res, err := c.do(http.MethodGet, c.endpoint+"/tasks", nil)
	if err != nil {
		return tasks, err
	}

	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&tasks)
		return
	} else {
		err = &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return
	}
}

func (c *Client) TaskCreate(task *model.Task) error {
	res, err := c.do(http.MethodPost, c.endpoint+"/tasks", task)
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
