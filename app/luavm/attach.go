package luavm

import (
	"encoding/json"
	"fmt"
)

func (s Service) GetAttach() (string, error) {
	data, err := json.Marshal(s.Task)
	if err != nil {
		return string(data), err
	}
	return string(data), nil
}

func (s Service) SetAttach(raw string) error {
	if err := json.Unmarshal([]byte(raw), s.Task); err != nil {
		return err
	}

	// Task
	if err := s.orm.Model(&s.Task).Select("files", "extra").Updates(&s.Task).Error; err != nil {
		return err
	}

	// Job
	if err := s.orm.Model(&s.Task.Job).Select("files", "extra").Updates(&s.Task.Job).Error; err != nil {
		return err
	}
	return s.rdb.Set(s.ctx, fmt.Sprintf(topic_running, s.Task.ID), raw, s.timeout).Err()
}
