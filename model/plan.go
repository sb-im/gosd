package model

import (
	"time"
)

type Plan struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	NodeID      int64             `json:"node_id"`
	CreateAt    *time.Time        `json:"create_at,omitempty"`
	UpdateAt    *time.Time        `json:"update_at,omitempty"`
	Attachments map[string]string `json:"attachments"`
}

func NewPlan() *Plan {
	return &Plan{Attachments: make(map[string]string)}
}

type Plans []*Plan
