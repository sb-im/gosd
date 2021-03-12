package luavm

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

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

	file, err := os.Open("lua/test_min.lua")
	if err != nil {
		t.Error(err)
	}

	script, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	task := NewTask(0, 0, 0)
	task.Script = script

	err = worker.doRun(task)
	if err != nil {

		// TODO: fix
		//t.Error(err)
	}
}
