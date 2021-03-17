package luavm

import (
	"fmt"
	"encoding/json"
)

const (
	topic_notification = "plans/%d/notification"
)

type Notification struct {
	Time	int64  `json:"time"`
	Level uint8  `json:"level"`
	Msg		string `json:"msg"`
}

func (s *Service) Notification(notification *Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return s.State.Record(fmt.Sprintf(topic_notification, s.Task.PlanID), data)
}
