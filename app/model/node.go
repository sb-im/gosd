package model

type Node struct {
	Model

	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	Points JSON   `json:"points"`
}
