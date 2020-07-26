package luavm

import mqtt "github.com/eclipse/paho.mqtt.golang"

const (
	StatusReady = "ready"
	StatusError = "error"
	StatusProtect = "protect"
	StatusRunning = "running"
)

type Detail struct {
	Status string `json:"status"`
}

type StatusManager struct {
	Plans map[string]*Detail
	Client *mqtt.Client
}

func NewStatusManager(client *mqtt.Client) *StatusManager {
	return &StatusManager{
		Plans: map[string]*Detail{},
		Client: client,
	}
}

func (s *StatusManager) SetStatus(planID, status string) {
	if detail := s.Plans[planID]; detail == nil {
		detail = &Detail{
			Status: status,
		}
	} else {
		detail.Status = status
	}
}

func (s *StatusManager) GetStatus(planID string) string {
	if detail := s.Plans[planID]; detail == nil {
		return StatusReady
	} else {
		return s.Plans[planID].Status
	}
}

