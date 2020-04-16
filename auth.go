package main

import (
	"crypto/rand"
	"fmt"
)

type User struct {
}

type AccessGrant struct {
	Keys map[string]*User
}

func NewAccessGrant() *AccessGrant {
	return &AccessGrant{
		Keys: make(map[string]*User),
	}
}

func (this *AccessGrant) Grant(user *User) (string, error) {
	n := 32
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	key := fmt.Sprintf("%x", b)
	this.Keys[key] = user
	return key, nil
}
