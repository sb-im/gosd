package luavm

import (
	"testing"
	"time"

	"sb.im/gosd/app/model"
)

func TestLuaFiles(t *testing.T) {
	task := &model.Task{}
	task.ID = 1

	w := newWorker(t)
	ch := make(chan error)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- w.doRun(task, []byte(`
function main(task)
  print("### RUN Files RUN ###")

  xpcall(function()
    print(plan:GetFileContent("test_files"))
  end,
  function()
    plan:SetFileContent("test_files", "test.txt", "233")
  end)

  local filename, content = plan:GetFileContent("test_files")
  if content == "233" then
    plan:SetFileContent("test_files", "test2.txt", "456")
  else
    plan:SetFileContent("test_files", "test.txt", "233")
  end
  print(plan:GetFileContent("test_files"))

  print("### END Files END ###")
end
`))
	}()
}
