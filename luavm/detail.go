package luavm

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	Plans  map[string]*Detail
	Client mqtt.Client
}

func NewStatusManager(client mqtt.Client) *StatusManager {
	return &StatusManager{
		Plans:  map[string]*Detail{},
		Client: client,
	}
}

func (s *StatusManager) SetRunning(planID string, status interface{}) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return s.Client.Publish(fmt.Sprintf(topic_detail, planID), 1, true, data).Error()
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

	return s.Client.Publish(fmt.Sprintf(topic_detail, planID), 1, true, data).Error()
}

func (s *StatusManager) GetStatus(planID string) string {
	if _, ok := s.Plans[planID]; ok {
		return s.Plans[planID].Status
	} else {
		return StatusReady
	}
}
