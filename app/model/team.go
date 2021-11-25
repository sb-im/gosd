package model

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	Model

	Name string `json:"name" form:"name"`
}

type UserTeam struct {
	UserID    uint           `gorm:"primaryKey"`
	TeamID    uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
