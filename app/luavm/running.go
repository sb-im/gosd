package luavm

import (
	"context"
	"encoding/json"
	"fmt"

	"sb.im/gosd/app/model"
)

const (
	topic_running = "tasks/%d/running"
)

func (w Worker) SetRunning(id uint, status interface{}) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return w.rdb.Set(context.Background(), fmt.Sprintf(topic_running, id), data, 0).Err()
}

func (w Worker) GetRunning(id uint) (*model.Task, error) {
	data, err := w.rdb.Get(context.Background(), fmt.Sprintf(topic_running, id)).Bytes()
	if err != nil {
		return nil, err
	}

	task := &model.Task{}
	if err := json.Unmarshal(data, task); err != nil {
		return task, err
	}
	return task, nil
}
