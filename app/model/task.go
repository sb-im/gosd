package model

type Task struct {
	Model

	Name   string `json:"name"`
	NodeID int64  `json:"node_id"`
	Files  JSON   `json:"files"`
	Extra  JSON   `json:"extra"`
}
