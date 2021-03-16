package luavm

import (
	"context"
	"testing"
	"time"

	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/model"
	"sb.im/gosd/mqttd"
	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"
)

func TestNewWorker(t *testing.T) {

	parse := config.NewParser()
	opts, err := parse.ParseEnvironmentVariables()
	if err != nil {
		panic(err)
	}

	db, err := database.NewConnectionPool(
		opts.DatabaseURL(),
		opts.DatabaseMinConns(),
		opts.DatabaseMaxConns(),
	)

	store := storage.NewStorage(db)

	// This unit test need:
	// planID == 1
	// jobID == 1
	if plan, _ := store.PlanByID(1); plan == nil {
		if err := store.CreatePlan(model.NewPlan()); err != nil {
			panic(err)
		}
	}
	if job, _ := store.PlanLogByID(1); job == nil {
		if err := store.CreatePlanLog(model.NewPlanLog()); err != nil {
			panic(err)
		}
	}

	s := state.NewState(opts.RedisURL())

	chI := make(chan mqttd.MqttRpc)
	chO := make(chan mqttd.MqttRpc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqttd := mqttd.NewMqttd(opts.MqttURL(), s, chI, chO)
	go mqttd.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chO, chI)
	go rpcServer.Run(ctx)

	worker := NewWorker(s, store, rpcServer)


	// test_min
	task := NewTask(0, 0, 0)
	task.Script = []byte(LuaMap["test_min"])

	err = worker.doRun(task)
	if err != nil {
		t.Error(err)
	}

	// test_dialog
	taskDialog := NewTask(0, 0, 0)
	taskDialog.Script = []byte(LuaMap["test_dialog"])

	go func() {
		for {
			worker.State.Record("plans/0/term", []byte("fine"))
			time.Sleep(1*time.Second)
			worker.State.Record("plans/0/term", []byte("confirm"))
			time.Sleep(1*time.Second)
		}
	}()
	err = worker.doRun(taskDialog)
	if err != nil {
		t.Error(err)
	}

	// test_extra
	taskExtra := NewTask(1, 1, 1)
	taskExtra.Script = []byte(LuaMap["test_extra"])

	err = worker.doRun(taskExtra)
	if err != nil {
		t.Error(err)
	}

	// test_files
	taskFiles := NewTask(1, 1, 1)
	taskFiles.Script = []byte(LuaMap["test_files"])

	err = worker.doRun(taskFiles)
	if err != nil {
		t.Error(err)
	}

	// test_blobs
	taskBlobs := NewTask(1, 1, 1)
	taskBlobs.Script = []byte(LuaMap["test_blobs"])

	err = worker.doRun(taskBlobs)
	if err != nil {
		t.Error(err)
	}

}
