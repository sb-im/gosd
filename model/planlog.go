package model

import (
	"time"
)

type PlanLog struct {
	ID          int64             `json:"id"`
	LogID       int64             `json:"log_id"`
	PlanID      int64             `json:"plan_id"`
	CreateAt    *time.Time        `json:"create_at,omitempty"`
	UpdateAt    *time.Time        `json:"update_at,omitempty"`
	Attachments map[string]string `json:"attachments"`
	Extra       map[string]string `json:"extra"`
}

func NewPlanLog() *PlanLog {
	return &PlanLog{
		Attachments: make(map[string]string),
		Extra:       make(map[string]string),
	}
}

type PlanLogs []*PlanLog
