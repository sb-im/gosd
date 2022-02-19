package client

import (
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
