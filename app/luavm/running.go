package luavm

import (
	"context"
	"encoding/json"
	"fmt"

	"sb.im/gosd/app/model"
)

const (
	topic_running = "plans/%d/running"
)

func (w *Worker) SetRunning(planID uint, status interface{}) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return w.rdb.Set(context.Background(), fmt.Sprintf(topic_running, planID), data, 0).Err()
}

func (w *Worker) GetRunning(planID uint) (*model.Task, error) {
	data, err := w.rdb.Get(context.Background(), fmt.Sprintf(topic_running, planID)).Bytes()
	if err != nil {
		return nil, err
	}

	task := &model.Task{}
	if err := json.Unmarshal(data, task); err != nil {
		return task, err
	}
	return task, nil
}
