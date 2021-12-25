package model

type Task struct {
	Model

	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	NodeID int64  `json:"node_id" form:"node_id"`
	Files  JSON   `json:"files" swaggertype:"string"`
	Extra  JSON   `json:"extra" swaggertype:"string"`
}