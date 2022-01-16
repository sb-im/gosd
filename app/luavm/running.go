package luavm

import (
	"encoding/json"
	"fmt"
)

const (
	topic_running = "tasks/%d/running"
)

func (s Service) running(status interface{}) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return s.rdb.Set(s.ctx, fmt.Sprintf(topic_running, s.Task.ID), data, s.timeout).Err()
}
