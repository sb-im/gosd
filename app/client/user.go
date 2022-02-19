package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) UserIndex() (users []model.User, err error) {
	res, err := http.Get(c.endpoint + "/users")
	if err != nil {
		return users, err
	}

	err = json.NewDecoder(res.Body).Decode(&users)
	return
}

func (c *Client) UserCreate(user interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(user); err != nil {
		return err
	}
	res, err := http.Post(c.endpoint+"/users", "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return json.NewDecoder(res.Body).Decode(user)
	} else {
		return errors.New("Create Failed")
	}
}
