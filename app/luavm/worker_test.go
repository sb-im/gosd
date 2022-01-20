package luavm

import (
	"testing"
	"time"

	lualib "sb.im/gosd/app/luavm/lua"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/storage"

	"sb.im/gosd/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newWorker(t *testing.T) *Worker {
	return helpTestNewWorker(t, []byte{})
}

func helpTestNewWorker(t *testing.T, script []byte) *Worker {
	parse := config.NewParser()
	opts, err := parse.ParseEnvironmentVariables()
	if err != nil {
		panic(err)
	}

	dsn := opts.DatabaseURL()
	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return NewWorker(orm, redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}), storage.NewStorage("/tmp"), nil, script)
}

func newTestTask(t *testing.T) *model.Task {
	parse := config.NewParser()
	opts, err := parse.ParseEnvironmentVariables()
	if err != nil {
		panic(err)
	}

	if orm, err := gorm.Open(postgres.Open(opts.DatabaseURL()), &gorm.Config{}); err != nil {
		t.Error(err)
		return nil
	} else {
		task := model.Task{
			Name:   "Unit Test",
			TeamID: 1,
			NodeID: 1,
		}
		orm.Create(&task)

		job := model.Job{
			Task: task,
		}
		orm.Create(&job)
		task.Job = &job
		return &task
	}
}

func TestNewWorker(t *testing.T) {
	if newWorker(t) == nil {
		t.Error("No New Worker")
	}
}

func TestLuaScript(t *testing.T) {
	tests := []string{
		"test_min.lua",
		"test_geo.lua",
		//"test_dialog.lua",
		//"test_rpc.lua",
	}

	for _, name := range tests {
		luaScript(t, name)
	}
}

func luaScript(t *testing.T, name string) {
	w := newWorker(t)
	if script, err := lualib.LuaFile.ReadFile(name); err != nil {
		t.Error(err)
	} else {
		w.doRun(&model.Task{}, script)
	}

	w.Close()
}

func TestMultipleTask(t *testing.T) {
	script := []byte(`
function main(task)
  print("### RUN Multiple RUN ###")

  sleep("1s")

  print("### END Multiple END ###")
end
`)

	w := helpTestNewWorker(t, script)

	task := newTestTask(t)
	task2 := newTestTask(t)
	task3 := newTestTask(t)
	task3.NodeID = 3

	if err := w.RunTask(task); err != nil {
		t.Error(err)
	}

	if err := w.RunTask(task); err == nil {
		t.Error("duplicate task")
	}

	if err := w.RunTask(task2); err == nil {
		t.Error("duplicate node, should error")
	}

	if err := w.RunTask(task2); err == nil {
		t.Error("duplicate task2")
	}

	if err := w.RunTask(task3); err != nil {
		t.Error(err)
	}

	if err := w.RunTask(task3); err == nil {
		t.Error("duplicate task")
	}

	time.Sleep(2 * time.Second)
}
