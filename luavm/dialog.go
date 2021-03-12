package luavm

import (
	"encoding/json"
	"fmt"
)

const (
	topic_dialog = "plans/%d/dialog"
)

type Dialog struct {
	Name    string        `json:"name,omitempty"`
	Message string        `json:"message,omitempty"`
	Level   string        `json:"level,omitempty"`
	Items   []*DialogItem `json:"items,omitempty"`
	Buttons []*DialogItem `json:"buttons,omitempty"`
}

type DialogItem struct {
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
	Level   string `json:"level,omitempty"`
}

func (s *Service) CleanDialog() error {
	return s.ToggleDialog(&Dialog{})
}

func (s *Service) ToggleDialog(dialog *Dialog) error {
	data, err := json.Marshal(dialog)
	if err != nil {
		return err
	}
	return s.State.Record(fmt.Sprintf(topic_dialog, s.Task.planID), data)
}
