package luavm

import (
	"io/ioutil"
	"os"
	"testing"

	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/storage"
	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"
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

	s := state.NewState()
	s.Mqtt = &jsonrpc2mqtt.MockClient{}
	worker := NewWorker(s, store)

	file, err := os.Open("lua/test_min.lua")
	if err != nil {
		t.Error(err)
	}

	script, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	p := &Task{
		NodeID: "1",
		URL:    "1/12/3/4/4",
		Script: script,
	}

	err = worker.doRun(p)
	if err != nil {
		t.Error(err)
	}
}
