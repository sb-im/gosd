package model

import (
	"encoding/json"
	"strconv"
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
	Extra       map[string]string `json:"extra"`
}

func NewPlan() *Plan {
	return &Plan{
		Attachments: make(map[string]string),
		Extra:       make(map[string]string),
	}
}

func (plan *Plan) MarshalJSON() ([]byte, error) {
	cycle_types_id, err := strconv.Atoi(plan.Extra["cycle_types_id"])
	if err != nil {
		cycle_types_id = 0
	}

	return json.Marshal(&struct {
		ID             int64             `json:"id"`
		Name           string            `json:"name"`
		Description    string            `json:"description"`
		File           string            `json:"map_path"`
		Node_id        int64             `json:"node_id"`
		Cycle_types_id int               `json:"cycle_types_id"`
		Attachments    map[string]string `json:"attachments"`
		Extra          map[string]string `json:"extra"`
	}{
		ID:             plan.ID,
		Name:           plan.Name,
		Description:    plan.Description,
		File:           plan.Attachments["file"],
		Node_id:        plan.NodeID,
		Attachments:    plan.Attachments,
		Extra:          plan.Extra,
		Cycle_types_id: cycle_types_id,
	})
}

type Plans []*Plan
