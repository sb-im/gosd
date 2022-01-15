package luavm

import (
	"testing"
)

func TestLuaAttach(t *testing.T) {
	task := newTestTask(t)

	w := newWorker(t)
	if err := w.doRun(task, []byte(`
function main(task)
  print("### RUN Attach RUN ###")

  local extra = task:GetExtra()
  extra["xxx"] = "aa"
  extra["ccc"] = "aaa"
  task:SetExtra(extra)

  local job_extra = task:GetJobExtra()
  job_extra["ttt"] = "xxx"
  print("Job Extra:", json.encode(job_extra))
  task:SetJobExtra(job_extra)

  print("### END Attach END ###")
end
`)); err != nil {
		t.Error(err)
	}

	w2 := newWorker(t)
	if err := w2.doRun(task, []byte(`
function main(task)
  print("### RUN Attach 2 RUN ###")

  local extra = task:GetExtra()
  if extra["xxx"] ~= "aa" then
    error("task extra error")
  end
  if extra["ccc"] ~= "aaa" then
    error("task extra error")
  end

  local job_extra = task:GetJobExtra()
  if job_extra["ttt"] ~= "xxx" then
    error("job extra error")
  end

  print("### END Attach 2 END ###")
end
`)); err != nil {
		t.Error(err)
	}

}
