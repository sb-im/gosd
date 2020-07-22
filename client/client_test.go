package client

import (
	"testing"
)

func TestLogin(t *testing.T) {
	client := NewClient("http://localhost:8000/gosd", "ttt", "123456")
	if err := client.Login(); err != nil {
		t.Error(err)
	}

	if err := client.GetNodes(); err != nil {
		t.Error(err)
	}
}
