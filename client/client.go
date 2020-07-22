package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"sb.im/gosd/model"
)

const (
	apiLogin = "/oauth/token"
)

// Client holds API procedure calls.
type Client struct {
	Endpoint string
	User     *model.Login
	Token    *model.Token
}

func NewClient(endpoint string, credentials ...string) *Client {
	return &Client{
		Endpoint: endpoint,
		User: &model.Login{
			Username: credentials[0],
			Password: credentials[1],
		},
	}
}

func (c *Client) Login() error {
	toFormData := func(form map[string]string) (b bytes.Buffer, contentType string, err error) {

		w := multipart.NewWriter(&b)
		contentType = w.FormDataContentType()

		var fw io.Writer

		for key, value := range form {
			fw, err = w.CreateFormField(key)
			if err != nil {
				return
			}
			if _, err = io.Copy(fw, bytes.NewReader([]byte(value))); err != nil {
				return
			}
		}
		w.Close()
		return
	}

	f := make(map[string]string)
	f["username"] = c.User.Username
	f["password"] = c.User.Password
	f["grant_type"] = "password"

	b, contentType, err := toFormData(f)

	req, err := http.NewRequest(http.MethodPost, c.Endpoint+apiLogin, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var token model.Token
	if err := json.Unmarshal(d, &token); err != nil {
		return err
	}

	c.Token = &token
	fmt.Println(token)

	return nil
}

func (c *Client) GetNodes() error {
	req, err := http.NewRequest("GET", c.Endpoint+"/api/v1/nodes/", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token.AccessToken)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", d)

	return nil
}
