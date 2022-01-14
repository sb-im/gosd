package luavm

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	topic_notification = "tasks/%d/notification"
)

type Notification struct {
	Time  int64  `json:"time"`
	Level uint8  `json:"level"`
	Msg   string `json:"msg"`
}

func (s Service) Notification(notification *Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return s.rdb.Set(context.Background(), fmt.Sprintf(topic_notification, s.Task.ID), data, 0).Err()
}
