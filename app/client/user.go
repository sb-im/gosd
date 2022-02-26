package client

import (
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) UserIndex() (users []model.User, err error) {
	res, err := c.do(http.MethodGet, c.endpoint+"/users", nil)
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
	res, err := c.do(http.MethodPost, c.endpoint+"/users", user)
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
	res, err := c.do(http.MethodPatch, c.endpoint+"/users/"+id, user)
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
	res, err := c.do(http.MethodPost, c.endpoint+"/users/"+userId+"/teams/"+teamId, nil)
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
