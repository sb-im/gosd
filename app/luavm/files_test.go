package luavm

import (
	"context"
	"testing"

	"sb.im/gosd/app/model"
)

func TestLuaFiles(t *testing.T) {
	task := &model.Task{}
	task.ID = 1

	job := model.Job{
		TaskID: task.ID,
	}
	task.Job = &job
	task.Job.ID = 1

	w := newWorker(t)
	if err := w.doRun(context.Background(), task, []byte(`
function main(task)
  print("### RUN Files RUN ###")

  xpcall(function()
    print(task:GetFileContent("test_files"))
  end,
  function()
    task:SetFileContent("test_files", "test.txt", "233")
  end)

  local filename, content = task:GetFileContent("test_files")
  if content == "233" then
    task:SetFileContent("test_files", "test2.txt", "456")
  else
    task:SetFileContent("test_files", "test.txt", "233")
  end
  print(task:GetFileContent("test_files"))
  print(task:FileUrl("test_files_url"))

  print("### END Files END ###")
end
`)); err != nil {
		t.Error(err)
	}
}
