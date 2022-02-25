package v3

import "testing"

func TestBasicAuth(t *testing.T) {
	username, password := "demo", "p@55w0rd"
	value := basicAuthEncode(username, password)
	if value != "Basic ZGVtbzpwQDU1dzByZA==" {
		t.Error(value)
	}
	user, pass, err := basicAuthDecode(value)
	if err != nil {
		t.Error(err)
	}
	if user != username || pass != password {
		t.Error(user, pass)
	}
}
