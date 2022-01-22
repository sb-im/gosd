package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c Client) TaskIndex() (tasks []model.Task, err error) {
	res, err := http.Get(c.endpoint + "/gosd/api/v3/tasks")
	if err != nil {
		return tasks, err
	}

	err = json.NewDecoder(res.Body).Decode(&tasks)
	return
}

func (c Client) TaskCreate(task *model.Task) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(task); err != nil {
		return err
	}
	res, err := http.Post(c.endpoint+"/gosd/api/v3/tasks", "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return json.NewDecoder(res.Body).Decode(task)
	} else {
		return errors.New("Create Failed")
	}
}
