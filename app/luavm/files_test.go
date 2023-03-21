package luavm

import (
	"context"
	"testing"
)

func TestLuaFiles(t *testing.T) {
	w := newWorker(t)
	task := helpTestNewTask(t, "Unit Test Lua Files", w)
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
