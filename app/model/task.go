package model

import "gorm.io/gorm"

type Task struct {
	Model

	Name   string `json:"name" form:"name"`
	TeamID uint   `json:"-"`
	NodeID uint   `json:"node_id" form:"node_id"`
	Job    *Job   `json:"job,omitempty" form:"-"`
	Files  JSON   `json:"files"`
	Extra  JSON   `json:"extra"`
}

func (t *Task) BeforeSave(tx *gorm.DB) error {
	if !jsonIsObject(t.Files) {
		t.Files = JSON("{}")
	}
	if !jsonIsObject(t.Extra) {
		t.Extra = JSON("{}")
	}
	return nil
}
