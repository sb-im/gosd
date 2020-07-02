package model

type Login struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
