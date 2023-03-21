package luavm

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRunning(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w := newWorker(t)
	task := helpTestNewTask(t, "Unit Test Lua Running", w)
	ch := make(chan error)
	go func() {
		ch <- w.doRun(context.Background(), task, []byte(`
function main(task)
  print("### RUN Running RUN ###")

  sleep("2s")

  print("### END Running END ###")
end
`))
	}()

	time.Sleep(time.Second)
	if val, err := w.rdb.Get(ctx, fmt.Sprintf(topic_running, task.ID)).Result(); err != nil {
		t.Error(val, err)
	} else if val == "{}" {
		t.Error(val)
	}

	time.Sleep(2 * time.Second)
	if val, err := w.rdb.Get(ctx, fmt.Sprintf(topic_running, task.ID)).Result(); err != nil {
		t.Error(val, err)
	} else if val != "{}" {
		t.Error(val)
	}

	select {
	case <-ctx.Done():
		t.Error("Time Out")
	case err := <-ch:
		if err != nil {
			t.Error(err)
		}
	}
}
