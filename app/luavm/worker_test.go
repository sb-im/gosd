package luavm

import (
	"context"
	"testing"
	"time"

	"sb.im/gosd/app/config"
	"sb.im/gosd/app/luavm/lib"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testTeamId = 1
	testNodeId = 1
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

	rdb := redis.NewClient(redisOpt)
	rdb.ConfigSet(context.Background(), "notify-keyspace-events", "$KEx")

	// use tempDir
	return NewWorker(Config{
		Instance: cfg.Instance,
		BaseURL:  cfg.BaseURL,
	}, service.NewService(cfg, orm, rdb, storage.NewStorage(t.TempDir())), nil, script)
}

func helpTestNewTask(t *testing.T, name string, w *Worker) *model.Task {
	task := &model.Task{
		Name:   name,
		TeamID: testTeamId,
		NodeID: testNodeId,
	}
	if err := w.srv.Orm().FirstOrCreate(task, model.Task{Name: name}).Error; err != nil {
		t.Error(err)
	}

	job, err := w.srv.ScheduleCreateJob(context.Background(), task.ID)
	if err != nil {
		t.Error(err)
	}
	task.Job = job
	return task
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
	task := helpTestNewTask(t, "Unit Test Lua Scriptl", w)
	if script, err := lib.File.ReadFile(name); err != nil {
		t.Error(err)
	} else {
		if err := w.doRun(context.Background(), task, script); err != nil {
			t.Error(err)
		}
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

	task := helpTestNewTask(t, "Unit Test MultipleTask 1", w)
	task2 := helpTestNewTask(t, "Unit Test MultipleTask 2", w)
	task3 := helpTestNewTask(t, "Unit Test MultipleTask 3", w)
	task3.NodeID = 3

	ctx := context.Background()

	if err := w.AddTask(ctx, task); err != nil {
		t.Error(err)
	}

	if err := w.AddTask(ctx, task); err == nil {
		t.Error("duplicate task")
	}

	if err := w.AddTask(ctx, task2); err == nil {
		t.Error("duplicate node, should error")
	}

	if err := w.AddTask(ctx, task2); err == nil {
		t.Error("duplicate task2")
	}

	if err := w.AddTask(ctx, task3); err != nil {
		t.Error(err)
	}

	if err := w.AddTask(ctx, task3); err == nil {
		t.Error("duplicate task")
	}

	time.Sleep(2 * time.Second)
}
