package luavm

import (
	"fmt"

	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

type Worker struct {
	Queue     chan *Task
	State     *state.State
	MqttProxy *jsonrpc2mqtt.MqttProxy
}

func NewWorker(s *state.State) *Worker {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(s.Mqtt)
	return &Worker{
		Queue:     make(chan *Task, 1024),
		State:     s,
		MqttProxy: mqttProxy,
	}
}

func (w Worker) Run() {
	for task := range w.Queue {
		if err := w.doRun(task); err != nil {
			fmt.Println(err)
		}
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
		Protect: true,
	}, lua.LString(task.NodeID)); err != nil {
		return err
	}

	return nil
}

func (w Worker) LoadMod(l *lua.LState, task *Task) {
	rpc := &LRpc{
		MqttProxy: w.MqttProxy,
	}

	service := &LService{
		Task:   task,
		State:  w.State,
		NodeID: task.NodeID,
	}

	l.SetGlobal("SD", luar.New(l, service))

	l.SetGlobal("get_id", l.NewFunction(service.GetID))
	l.SetGlobal("get_msg", l.NewFunction(service.GetMsg))
	l.SetGlobal("get_status", l.NewFunction(service.GetStatus))

	l.SetGlobal("rpc_notify", l.NewFunction(rpc.notify))
	l.SetGlobal("rpc_async", l.NewFunction(rpc.asyncCall))
	l.SetGlobal("rpc_call", l.NewFunction(rpc.call))
}
