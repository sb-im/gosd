package state

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func Test_NodePutGet(t *testing.T) {
	state := NewState(os.Getenv("REDIS_URL"))
	id := "233"
	msg := "test"
	str := `{"test":23333333333333}`
	state.Record(fmt.Sprintf("nodes/%s/msg/%s", id, msg), []byte(str))

	payload, err := state.GetNodeMsg(id, msg)
	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(payload, []byte(str)) {
		t.Errorf("%s\n", payload)
	}
}

func Test_NodeGetIdNil(t *testing.T) {
	state := NewState(os.Getenv("REDIS_URL"))
	payload, err := state.GetNodeMsg("2", "test")
	if err == nil {
		t.Error("Should no payload")
		t.Errorf("%s\n", payload)
	}
}

func Test_NodeGetMsgNil(t *testing.T) {
	state := NewState(os.Getenv("REDIS_URL"))
	id := "233"
	msg := "test"
	str := `{"test":23333333333333}`
	state.Record(fmt.Sprintf("nodes/%s/msg/%s", id, msg), []byte(str))

	payload, err := state.GetNodeMsg(id, "test2")
	if err == nil {
		t.Error("Should no Message")
	}

	if bytes.Equal(payload, []byte(str)) {
		t.Errorf("%s\n", payload)
	}
}
