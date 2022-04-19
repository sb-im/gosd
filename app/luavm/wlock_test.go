package luavm

import (
	"strconv"
	"testing"
)

func TestPreTaskCheck(t *testing.T) {
	task := newTestTask(t)
	taskID := strconv.Itoa(int(task.ID))
	nodeID := task.NodeID

	w := newWorker(t)

	// check pass
	if err := w.preTaskCheck(task); err != nil {
		t.Error(err)
	}

	if err := w.lockTaskSet(taskID); err != nil {
		t.Error(err)
	}

	// check not pass
	if err := w.preTaskCheck(task); err == nil {
		t.Error("check should not pass")
	}

	if err := w.lockTaskDel(taskID); err != nil {
		t.Error(err)
	}

	// check pass
	if err := w.preTaskCheck(task); err != nil {
		t.Error(err)
	}

	if err := w.lockNodeSet(nodeID); err != nil {
		t.Error(err)
	}

	// check not pass
	if err := w.preTaskCheck(task); err == nil {
		t.Error("check should not pass")
	}

	if err := w.lockTaskSet(taskID); err != nil {
		t.Error(err)
	}

	// check not pass
	if err := w.preTaskCheck(task); err == nil {
		t.Error("check should not pass")
	}

	if err := w.lockNodeDel(nodeID); err != nil {
		t.Error(err)
	}

	// check not pass
	if err := w.preTaskCheck(task); err == nil {
		t.Error("check should not pass")
	}

	if err := w.lockTaskDel(taskID); err != nil {
		t.Error(err)
	}

	// check pass
	if err := w.preTaskCheck(task); err != nil {
		t.Error(err)
	}
}
