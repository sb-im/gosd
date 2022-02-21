package model

type Task struct {
	Model

	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	NodeID uint   `json:"node_id" form:"node_id"`
	Job    *Job   `json:"job,omitempty" form:"-"`
	Files  JSON   `json:"files"`
	Extra  JSON   `json:"extra"`
}
