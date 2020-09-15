package luavm

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic_detail = "plans/%s/status"
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
	Plans  map[string]*Detail
	Client mqtt.Client
}

func NewStatusManager(client mqtt.Client) *StatusManager {
	return &StatusManager{
		Plans:  map[string]*Detail{},
		Client: client,
	}
}

func (s *StatusManager) SetStatus(planID, status string) error {
	detail := s.Plans[planID]
	if detail == nil {
		detail = &Detail{
			Status: status,
		}
	} else {
		detail.Status = status
	}
	data, err := json.Marshal(detail)
	if err != nil {
		return err
	}

	return s.Client.Publish(fmt.Sprintf(topic_detail, planID), 1, true, data).Error()
}

func (s *StatusManager) GetStatus(planID string) string {
	if detail := s.Plans[planID]; detail == nil {
		return StatusReady
	} else {
		return s.Plans[planID].Status
	}
}
