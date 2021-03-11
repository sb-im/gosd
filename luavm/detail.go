package luavm

import (
	"encoding/json"
	"fmt"

	"sb.im/gosd/state"
)

const (
	topic_detail = "plans/%s/running"
)

const (
	StatusReady   = "ready"
	StatusError   = "error"
	StatusProtect = "protect"
	StatusRunning = "running"
)

type Detail struct {
	Status string `json:"status"`
}

type StatusManager struct {
	Plans map[string]*Detail
	State *state.State
}

func NewStatusManager(store *state.State) *StatusManager {
	return &StatusManager{
		Plans: map[string]*Detail{},
		State: store,
	}
}

func (s *StatusManager) SetRunning(planID string, status interface{}) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return s.State.Record(fmt.Sprintf(topic_detail, planID), data)
}

func (s *StatusManager) SetStatus(planID, status string) error {
	if _, ok := s.Plans[planID]; ok {
		s.Plans[planID].Status = status
	} else {
		s.Plans[planID] = &Detail{
			Status: status,
		}
	}
	data, err := json.Marshal(s.Plans[planID])
	if err != nil {
		return err
	}

	return s.State.Record(fmt.Sprintf(topic_detail, planID), data)
}

func (s *StatusManager) GetStatus(planID string) string {
	if _, ok := s.Plans[planID]; ok {
		return s.Plans[planID].Status
	} else {
		return StatusReady
	}
}
