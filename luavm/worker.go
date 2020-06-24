package luavm

import (
	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type Worker struct {
	Queue chan *Task
	State *state.State
}

func NewWorker(s *state.State) *Worker {
	return &Worker{
		Queue: make(chan *Task, 1024),
		State: s,
	}
}

func (w Worker) Run() {
	for task := range w.Queue {
		w.doRun(task)
	}
}

func (w Worker) doRun(task *Task) error {
	l := lua.NewState()
	defer l.Close()
	luajson.Preload(l)

	w.LoadMod(l, task)

	var err error
	if len(task.Script) == 0 {
		err = l.DoString(LuaMap["default"])
	} else {
		err = l.DoString(string(task.Script))
	}

	if err != nil {
		return err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("run"),
		NRet:    1,
		Protect: false,
	}, lua.LString(task.NodeID)); err != nil {
		return err
	}

	return nil
}

func (w Worker) LoadMod(l *lua.LState, task *Task) {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(w.State.Mqtt)
	rpc := &LRpc{
		MqttProxy: mqttProxy,
	}

	service := &LService{
		State:  w.State,
		NodeID: task.NodeID,
	}

	l.SetGlobal("get_id", l.NewFunction(service.GetID))
	l.SetGlobal("get_msg", l.NewFunction(service.GetMsg))
	l.SetGlobal("get_status", l.NewFunction(service.GetStatus))

	l.SetGlobal("rpc_notify", l.NewFunction(rpc.notify))
	l.SetGlobal("rpc_async", l.NewFunction(rpc.asyncCall))
	l.SetGlobal("rpc_call", l.NewFunction(rpc.call))
}
