package model

import (
	"time"

	"gorm.io/gorm"
	"sb.im/gosd/app/helper"
)

type Node struct {
	ID_ID     int            `json:"id_id" gorm:"primaryKey"`
	ID        string         `json:"id" gorm:"unique"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	Secret string `json:"-"`
	Points JSON   `json:"points"`
}

func (n *Node) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = helper.GenNumberSecret(6)
	}

	if n.Secret == "" {
		n.Secret = helper.GenSecret(16)
	}
	return nil
}

func (n *Node) BeforeSave(tx *gorm.DB) error {
	if !jsonIsArray(n.Points) {
		n.Points = JSON("[]")
	}
	return nil
}
