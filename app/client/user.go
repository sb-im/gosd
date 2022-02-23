package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) UserIndex() (users []model.User, err error) {
	req, err := http.NewRequest("GET", c.endpoint+"/users", nil)
	if err != nil {
		return users, err
	}
	res, err := c.do(req)
	if err != nil {
		return users, err
	}

	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&users)
		return
	} else {
		err = &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return
	}
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
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}

func (c *Client) UserUpdate(id string, user interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(user); err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", c.endpoint+"/users/"+id, buf)
	if err != nil {
		return err
	}
	res, err := c.do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return json.NewDecoder(res.Body).Decode(user)
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}

func (c *Client) UserAddTeam(userId, teamId string) error {
	res, err := http.Post(c.endpoint+"/users/"+userId+"/teams/"+teamId, "application/json", nil)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return nil
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}
