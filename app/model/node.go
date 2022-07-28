package model

import (
	"gorm.io/gorm"
	"sb.im/gosd/app/helper"
)

type Node struct {
	Model

	UUID   string `json:"uuid" gorm:"uniqueIndex;column:uuid"`
	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	Secret string `json:"-"`
	Points JSON   `json:"points"`
}

func (n *Node) BeforeCreate(tx *gorm.DB) error {
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
