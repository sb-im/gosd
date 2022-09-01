package luavm

import (
	"context"
	"testing"
)

func TestLuaTaskName(t *testing.T) {
	task := newTestTask(t)

	w := newWorker(t)

	if err := w.doRun(context.Background(), task, []byte(`
function main(task)
  print("### RUN Task Name RUN ###")

  if task.name ~= "`+task.Name+`" then
    error("task name is: " .. task.name)
  end

  print("### END Task Name END ###")
end
`)); err != nil {
		t.Error(err)
	}
}
