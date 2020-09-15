package luavm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"sb.im/gosd/jsonrpc2mqtt"

	jsonrpc "github.com/sb-im/jsonrpc-lite"
	jsonrpc2 "github.com/sb-im/jsonrpc2"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

var sequence uint64
var sequenceMutex sync.Mutex

func getSequence() string {
	sequenceMutex.Lock()
	id := strconv.FormatUint(sequence, 10)
	sequence++
	sequenceMutex.Unlock()
	return id
}

type LRpc struct {
	task      *Task
	pendings  map[string]chan []byte
	MqttProxy *jsonrpc2mqtt.MqttProxy
	Status    *StatusManager
}

func NewLRpc(task *Task, proxy *jsonrpc2mqtt.MqttProxy, status *StatusManager) *LRpc {
	return &LRpc{
		task:      task,
		pendings:  make(map[string]chan []byte),
		MqttProxy: proxy,
		Status:    status,
	}
}

func (m *LRpc) Kill() error {
	for id, ch := range m.pendings {
		rpc := jsonrpc.NewErrors(id)
		rpc.Errors.InternalError("Be killed")
		data, err := rpc.ToJSON()
		if err != nil {
			return err
		}
		ch <- data
	}

	return nil
}

func req_jsonrpc(raw []byte) ([]byte, error) {
	bit13_timestamp := string([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))[:13])
	rpc_id := "gosd.0-" + bit13_timestamp + "-" + getSequence()
	fmt.Println(rpc_id)

	jsonrpc_req := jsonrpc2.WireRequest{}
	err := json.Unmarshal(raw, &jsonrpc_req)
	if err != nil {
		fmt.Println(err)
	}
	jsonrpc_req.ID = &jsonrpc2.ID{
		Name: rpc_id,
	}

	return json.Marshal(jsonrpc_req)
}

func res_jsonrpc(raw []byte) ([]byte, error) {
	type rpc struct {
		Result *json.RawMessage `json:"result,omitempty"`
		Error  *jsonrpc2.Error  `json:"error,omitempty"`
	}
	r := rpc{}
	json.Unmarshal(raw, &r)
	return json.Marshal(r)
}

func (m *LRpc) notify(L *lua.LState) int {
	raw, err := luajson.Encode(L.ToTable(2))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	jsonrpc_req := jsonrpc2.WireRequest{}
	err = json.Unmarshal(raw, &jsonrpc_req)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	req, _ := json.Marshal(jsonrpc_req)
	err = m.MqttProxy.Notify(L.ToString(1), req)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LString(""))
	return 1
}

func (m *LRpc) asyncCall(L *lua.LState) int {
	raw, err := luajson.Encode(L.ToTable(2))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	req, err := req_jsonrpc(raw)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	ch_L := L.ToChannel(3)
	ch := make(chan []byte)
	m.pendings[jsonrpc.ParseObject(req).ID.String()] = ch

	err = m.MqttProxy.AsyncRpc(L.ToString(1), req, ch)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	go func() {
		raw := <-ch
		r, _ := res_jsonrpc(raw)
		delete(m.pendings, jsonrpc.ParseObject(raw).ID.String())
		value, _ := luajson.Decode(L, r)
		ch_L <- value
	}()

	L.Push(lua.LString(""))
	return 1
}

func (m *LRpc) call(L *lua.LState) int {
	raw, err := luajson.Encode(L.ToTable(2))
	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}
	req, err := req_jsonrpc(raw)
	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var res []byte
	fmt.Println(m.Status.GetStatus(m.task.PlanID))
	if status := m.Status.GetStatus(m.task.PlanID); status != StatusRunning {
		rpc := jsonrpc.NewErrors(jsonrpc.ParseObject(req).ID)
		rpc.Errors.InternalError("Status is: " + status)
		res, err = rpc.ToJSON()
		if err != nil {
			L.Push(&lua.LTable{})
			L.Push(lua.LString(err.Error()))
		}
	} else {
		ch := make(chan []byte)
		m.pendings[jsonrpc.ParseObject(req).ID.String()] = ch

		err = m.MqttProxy.AsyncRpc(L.ToString(1), req, ch)
		if err != nil {
			L.Push(&lua.LTable{})
			L.Push(lua.LString(err.Error()))
			return 2
		}

		res = <-ch
		delete(m.pendings, jsonrpc.ParseObject(raw).ID.String())
	}

	r, err := res_jsonrpc(res)

	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}

	value, err := luajson.Decode(L, r)
	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(value)
	L.Push(lua.LString(""))
	return 2
}
