package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model

	TeamID uint   `json:"team_id"`
	Team   *Team  `json:"team,omitempty"`
	Teams  []Team `json:"teams,omitempty" gorm:"many2many:user_teams;"`

	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"-"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

func userHashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	password, err := userHashPassword([]byte(u.Password))
	u.Password = string(password)
	return err
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
