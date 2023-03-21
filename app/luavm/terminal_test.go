package luavm

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestLuaTerminal(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w := newWorker(t)
	task := helpTestNewTask(t, "Unit Test Lua Terminal", w)
	ch := make(chan error)
	go func() {
		ch <- w.doRun(context.Background(), task, []byte(`
function main(task)
  print("### RUN Terminal RUN ###")
  err = task:Puts("puts_1")
  if err ~= nil then
    error(err)
  end

  sleep("1s")

  msg, err = task:Gets()
  if err ~= nil then
    error(err)
  end

  if msg ~= 'get_1' then
    error(msg)
  end
  print("### END Terminal END ###")
end
`))
	}()

	time.Sleep(time.Second)
	if val, err := w.rdb.Get(context.Background(), fmt.Sprintf(topic_terminal, task.ID)).Result(); err != nil {
		t.Error(val, err)
	} else if val != "puts_1" {
		t.Error(val)
	}

	time.Sleep(2 * time.Second)
	if err := w.rdb.Set(context.Background(), fmt.Sprintf(topic_terminal, task.ID), "get_1", time.Second).Err(); err != nil {
		t.Error(err)
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
