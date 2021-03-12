package luavm

import (
	"encoding/json"
	"fmt"
)

const (
	topic_running = "plans/%d/running"
)

func (w *Worker) SetRunning(planID int64, status interface{}) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return w.State.Record(fmt.Sprintf(topic_running, planID), data)
}

func (w *Worker) GetRunning(planID int64) (*Task, error) {
	data, err := w.State.BytesGet(fmt.Sprintf(topic_running, planID))
	if err != nil {
		return nil, err
	}

	task := &Task{}
	if err := json.Unmarshal(data, task); err != nil {
		return task, err
	}
	return task, nil
}
