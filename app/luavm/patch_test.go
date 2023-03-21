package luavm

import (
	"context"
	"testing"
	"time"
)

func TestLuaPatchPrint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), "traceid", "CI-"+t.Name()), 3*time.Second)
	defer cancel()

	w := newWorker(t)
	task := helpTestNewTask(t, "Unit Test Lua PatchPrint", w)
	ch := make(chan error)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- w.doRun(ctx, task, []byte(`
function main(task)
  print("### RUN Patch RUN ###")

  print(0)
  print(0, 1, 2)

  print("aa")
  print("aa", "bb", "cc")

  print(true)
  print(true, false, true, false)


  print()
  print(0, "aa", true)
  print({"aa", "bb", "cc"})
  print({
    aa = 0,
    bb = "bb",
    cc = true,
  })

  print("### END Patch END ###")
end
`))
	}()

	select {
	case <-ctx.Done():
		t.Error("Time Out")
	case err := <-ch:
		if err != nil {
			t.Error(err)
		}
	}
}
