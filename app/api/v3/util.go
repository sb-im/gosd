package v3

import (
	"encoding/base64"
	"errors"
	"strings"
)

func basicAuthDecode(raw string) (string, string, error) {
	// example: "Basic ZGVtbzpwQDU1dzByZA=="
	v := strings.Split(raw, " ")
	if len(v) < 2 || v[0] != "Basic" {
		return "", "", errors.New("Not basic auth")
	}

	data, err := base64.StdEncoding.DecodeString(v[1])
	if err != nil {
		return "", "", err
	}
	r := strings.Split(string(data), ":")
	return r[0], r[1], nil
}

func basicAuthEncode(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
