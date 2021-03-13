package luavm

import (
	"context"
	"testing"
	"time"

	"sb.im/gosd/config"
	"sb.im/gosd/database"
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



}
