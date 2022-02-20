package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c *Client) TeamIndex() (teams []model.Team, err error) {
	res, err := http.Get(c.endpoint + "/teams")
	if err != nil {
		return teams, err
	}

	err = json.NewDecoder(res.Body).Decode(&teams)
	return
}

func (c *Client) TeamCreate(team interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(team); err != nil {
		return err
	}
	res, err := http.Post(c.endpoint+"/teams", "application/json", buf)
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
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(team); err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", c.endpoint+"/teams/"+id, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := (&http.Client{}).Do(req)
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
