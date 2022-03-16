package model

import "gorm.io/gorm"

type Node struct {
	Model

	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	Secret string `json:"-"`
	Points JSON   `json:"points"`
}

func (n *Node) BeforeSave(tx *gorm.DB) error {
	if !jsonIsArray(n.Points) {
		n.Points = JSON("[]")
	}
	return nil
}
