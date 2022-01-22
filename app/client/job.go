package client

import (
	"errors"
	"fmt"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c Client) JobCreate(task *model.Task) error {
	res, err := http.Post(c.endpoint+fmt.Sprintf("/gosd/api/v3/tasks/%d/jobs", task.ID), "application/json", nil)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return nil
	} else {
		return errors.New("Create Failed")
	}
}
