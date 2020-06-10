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

func Test_SetNodeStatus(t *testing.T) {
	id := "233"
	msg := `{"code":0,"msg":"online","timestamp":"1591733101","status":{"link_id":5,"position_ok":true,"lat":"22.6876423001","lng":"114.2248673001","alt":"80.0001"}}`

	state := NewState()
	state.SetNodeStatus(id, []byte(msg))
	if state.Node[id].Status.Code != 0 {
		t.Error(state.Node[id].Status)
	}

	if !state.Node[id].Status.isConnect() {
		t.Error(state.Node[id].Status)
	}
}

func Test_SetNodeStatusOffline(t *testing.T) {
	id := "233"
	msg := `{"code":1,"msg":"offline","timestamp":"1591733101","status":{"link_id":5,"position_ok":true,"lat":"22.6876423001","lng":"114.2248673001","alt":"80.0001"}}`

	state := NewState()
	state.SetNodeStatus(id, []byte(msg))
	if state.Node[id].Status.Code != 1 {
		t.Error(state.Node[id].Status)
	}

	if state.Node[id].Status.isConnect() {
		t.Error(state.Node[id].Status)
	}
}
