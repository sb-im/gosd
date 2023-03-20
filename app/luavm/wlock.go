package luavm

import (
	"errors"
	"strconv"

	"sb.im/gosd/app/model"
)

func (w *Worker) preTaskCheck(task *model.Task) error {
	taskID := strconv.Itoa(int(task.ID))
	nodeID := strconv.Itoa(int(task.NodeID))

	if instance, _ := w.lockTaskGet(taskID); instance != "" {
		return errors.New("This Task already running at: " + instance)
	}

	if instance, _ := w.lockNodeGet(nodeID); instance != "" {
		return errors.New("This Node already running at: " + instance)
	}

	return nil
}

func (w *Worker) lockTaskSet(id string) error {
	return w.srv.LockTaskSet(id)
}

func (w *Worker) lockTaskGet(id string) (string, error) {
	return w.srv.LockTaskGet(id)
}

func (w *Worker) lockTaskDel(id string) error {
	return w.srv.LockTaskDel(id)
}

func (w *Worker) lockNodeSet(id string) error {
	return w.srv.LockNodeSet(id)
}

func (w *Worker) lockNodeGet(id string) (string, error) {
	return w.srv.LockNodeGet(id)
}

func (w *Worker) lockNodeDel(id string) error {
	return w.srv.LockNodeDel(id)
}
