package luavm

import (
	"fmt"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic_terminal = "plans/%s/term"
)

func (s *Service) IOGets() (string, error) {
	ch := make(chan []byte)
	topic := fmt.Sprintf(topic_terminal, s.Task.PlanID)
	token := s.State.Mqtt.Subscribe(topic, 2, func(client mqtt.Client, msg mqtt.Message) {
		ch <- msg.Payload()
	})

	if err := token.Error(); err != nil {
		return "", err
	}

	var raw []byte
	select {
	case <-s.ctx.Done():
		return "", errors.New("Be killed")
	case raw = <-ch:
	}

	token = s.State.Mqtt.Unsubscribe(topic)
	if err := token.Error(); err != nil {
		return string(raw), err
	}
	return string(raw), nil
}

func (s *Service) IOPuts(str string) error {
	return s.State.Mqtt.Publish(fmt.Sprintf(topic_terminal, s.Task.PlanID), 1, false, []byte(str)).Error()
}
