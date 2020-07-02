package model

import (
	"errors"
	"time"

	"miniflux.app/timezone"
)

// User represents a user in the system.
type User struct {
	ID          int64             `json:"id"`
	Username    string            `json:"username"`
	Password    string            `json:"password,omitempty"`
	Language    string            `json:"language"`
	Timezone    string            `json:"timezone"`
	LastLoginAt *time.Time        `json:"last_login_at,omitempty"`
	Group       *Group            `json:"group"`
	Extra       map[string]string `json:"extra"`
}

// NewUser returns a new User.
func NewUser() *User {
	return &User{
		Group: NewGroup(),
		Extra: make(map[string]string),
	}
}

// ValidateUserCreation validates new user.
func (u User) ValidateUserCreation() error {
	if err := u.ValidateUserLogin(); err != nil {
		return err
	}
	return u.ValidatePassword()
}

// ValidateUserModification validates user modification payload.
func (u User) ValidateUserModification() error {
	if u.Password != "" {
		return u.ValidatePassword()
	}
	return nil
}

// ValidateUserLogin validates user credential requirements.
func (u User) ValidateUserLogin() error {
	if u.Username == "" {
		return errors.New("The username is mandatory")
	}

	if u.Password == "" {
		return errors.New("The password is mandatory")
	}
	return nil
}

// ValidatePassword validates user password requirements.
func (u User) ValidatePassword() error {
	if u.Password != "" && len(u.Password) < 6 {
		return errors.New("The password must have at least 6 characters")
	}
	return nil
}

// UseTimezone converts last login date to the given timezone.
func (u *User) UseTimezone(tz string) {
	if u.LastLoginAt != nil {
		*u.LastLoginAt = timezone.Convert(tz, *u.LastLoginAt)
	}
}

// Users represents a list of users.
type Users []*User

// UseTimezone converts last login timestamp of all users to the given timezone.
func (u Users) UseTimezone(tz string) {
	for _, user := range u {
		user.UseTimezone(tz)
	}
}
