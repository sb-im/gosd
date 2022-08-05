package client

import (
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) TeamIndex() (teams []model.Team, err error) {
	res, err := c.do(http.MethodGet, c.endpoint+"/teams", nil)
	if err != nil {
		return teams, err
	}

	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&teams)
		return
	} else {
		err = &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return
	}
}

func (c *Client) TeamCreate(team interface{}) error {
	res, err := c.do(http.MethodPost, c.endpoint+"/teams", team)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		return json.NewDecoder(res.Body).Decode(team)
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}

func (c *Client) TeamUpdate(id string, team interface{}) error {
	res, err := c.do(http.MethodPatch, c.endpoint+"/teams/"+id, team)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return json.NewDecoder(res.Body).Decode(team)
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}

func (c *Client) TeamDestroy(id string) error {
	res, err := c.do(http.MethodDelete, c.endpoint+"/teams/"+id, nil)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	} else {
		err := &errMsg{
			status: res.Status,
		}
		json.NewDecoder(res.Body).Decode(err)
		return err
	}
}
