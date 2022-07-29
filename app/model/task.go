package model

import "gorm.io/gorm"

type Task struct {
	Model

	Name   string `json:"name" form:"name"`
	NodeID uint   `json:"node_id" form:"node_id"`
	TeamID uint   `json:"-"`
	Index  uint   `json:"index" gorm:"default:0;not null"`
	Files  JSON   `json:"files"`
	Extra  JSON   `json:"extra"`
	Job    *Job   `json:"job,omitempty" form:"-"`
	Jobs   []Job  `json:"jobs,omitempty" form:"-"`
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
