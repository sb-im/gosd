package luavm

import (
	"io/ioutil"
	"os"
	"testing"

	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"
)

func TestNewWorker(t *testing.T) {
	s := state.NewState()
	s.Mqtt = &jsonrpc2mqtt.MockClient{}
	worker := NewWorker(s)

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
