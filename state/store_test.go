package state

import (
	"bytes"
	"testing"
)

func Test_NodePutGet(t *testing.T) {
	state := NewState()
	id := "233"
	msg := "test"
	str := `{"test":23333333333333}`
	state.NodePut(id, msg, []byte(str))

	payload, err := state.NodeGet(id, msg)
	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(payload, []byte(str)) {
		t.Errorf("%s\n", payload)
	}
}

func Test_NodeGetIdNil(t *testing.T) {
	state := NewState()
	payload, err := state.NodeGet("2", "test")
	if err == nil {
		t.Error("Should no payload")
		t.Errorf("%s\n", payload)
	}
}

func Test_NodeGetMsgNil(t *testing.T) {
	state := NewState()
	id := "233"
	msg := "test"
	str := `{"test":23333333333333}`
	state.NodePut(id, msg, []byte(str))

	payload, err := state.NodeGet(id, "test2")
	if err == nil {
		t.Error("Should no Message")
	}

	if bytes.Equal(payload, []byte(str)) {
		t.Errorf("%s\n", payload)
	}
}
