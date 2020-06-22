package luavm

import (
	"io/ioutil"
	"os"
	"testing"

	"sb.im/gosd/state"
)

func TestNewWorker(t *testing.T) {
	s := state.NewState()
	worker := NewWorker(s)

	file, err := os.Open("test_min.lua")
	if err != nil {
		t.Error(err)
	}

	script, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	p := &plan{
		Id:     "1",
		LogID:  "2",
		NodeID: "3",
		Url:    "1/12/3/4/4",
		Script: script,
	}

	worker.doRun(p)

}
