package luavm

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic_terminal = "plans/%s/term"
	topic_dialog   = "plans/%s/dialog"
)

type Dialog struct {
	Name    string        `json:"name"`
	Message string        `json:"message,omitempty"`
	Level   string        `json:"level,omitempty"`
	Items   []*DialogItem `json:"items,omitempty"`
	Buttons []*DialogItem `json:"buttons,omitempty"`
	Inputs  []*DialogItem `json:"inputs,omitempty"`
}

type DialogItem struct {
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
	Level   string `json:"level,omitempty"`
}

func (s *LService) IOGets() (string, error) {
	ch := make(chan []byte)
	token := s.State.Mqtt.Subscribe(fmt.Sprintf(topic_terminal, s.Task.PlanID), 2, func(client mqtt.Client, msg mqtt.Message) {
		ch <- msg.Payload()
	})

	if err := token.Error(); err != nil {
		return "", err
	}
	return string(<-ch), nil
}

func (s *LService) IOPuts(str string) error {
	return s.State.Mqtt.Publish(fmt.Sprintf(topic_terminal, s.Task.PlanID), 1, false, []byte(str)).Error()
}

func (s *LService) ToggleDialog(dialog *Dialog) error {
	data, err := json.Marshal(dialog)
	if err != nil {
		return err
	}
	return s.State.Mqtt.Publish(fmt.Sprintf(topic_dialog, s.Task.PlanID), 1, true, data).Error()
}
