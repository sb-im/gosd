package luavm

import (
	"errors"
	"strconv"

	"sb.im/gosd/app/model"
)

const (
	lockTaskPrefix = "luavm.lock.task."
	lockNodePrefix = "luavm.lock.node."
)

func (w Worker) preTaskCheck(task *model.Task) error {
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

func (w Worker) lockTaskSet(id string) error {
	return w.rdb.Set(w.ctx, lockTaskPrefix+id, w.cfg.Instance, w.timeout).Err()
}

func (w Worker) lockTaskGet(id string) (string, error) {
	return w.rdb.Get(w.ctx, lockTaskPrefix+id).Result()
}

func (w Worker) lockTaskDel(id string) error {
	return w.rdb.Del(w.ctx, lockTaskPrefix+id).Err()
}

func (w Worker) lockNodeSet(id string) error {
	return w.rdb.Set(w.ctx, lockNodePrefix+id, w.cfg.Instance, w.timeout).Err()
}

func (w Worker) lockNodeGet(id string) (string, error) {
	return w.rdb.Get(w.ctx, lockNodePrefix+id).Result()
}

func (w Worker) lockNodeDel(id string) error {
	return w.rdb.Del(w.ctx, lockNodePrefix+id).Err()
}
