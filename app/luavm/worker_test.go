package luavm

import (
	"testing"
	"time"

	"sb.im/gosd/app/config"
	lualib "sb.im/gosd/app/luavm/lua"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/storage"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newWorker(t *testing.T) *Worker {
	return helpTestNewWorker(t, []byte{})
}

func helpTestNewWorker(t *testing.T, script []byte) *Worker {
	cfg := config.Parse()

	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}

	return NewWorker(orm, redis.NewClient(redisOpt), storage.NewStorage(cfg.StorageURL), nil, script)
}

func newTestTask(t *testing.T) *model.Task {
	cfg := config.Parse()

	if orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{}); err != nil {
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
