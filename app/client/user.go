package client

import (
	"encoding/json"
	"net/http"

	"sb.im/gosd/app/model"
)

func (c Client) UserIndex() (users []model.User, err error) {
	res, err := http.Get(c.endpoint + "/users")
	if err != nil {
		return users, err
	}

	err = json.NewDecoder(res.Body).Decode(&users)
	return
}
