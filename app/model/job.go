package model

import "time"

type Job struct {
	Model

	Index     uint      `json:"index" gorm:"index:idx_index,unique;not null"`
	TaskID    uint      `json:"-" gorm:"index:idx_index,unique;not null"`
	Task      Task      `json:"-"`
	StartedAt time.Time `json:"started_at" form:"started_at"`
	Duration  int       `json:"duration"`
	Files     JSON      `json:"files"`
	Extra     JSON      `json:"extra"`
}
