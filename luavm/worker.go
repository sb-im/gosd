package luavm

import (
	"fmt"

	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"
	"sb.im/gosd/task"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type Worker struct {
	Queue chan task.Task
	State *state.State
}

func NewWorker(s *state.State) *Worker {
	return &Worker{
		Queue: make(chan task.Task, 1024),
		State: s,
	}
}

func (w Worker) Run() {
	for plan := range w.Queue {
		w.doRun(plan)
	}
}

func (w Worker) doRun(plan task.Task) {
	l := lua.NewState()
	defer l.Close()
	luajson.Preload(l)

	w.LoadMod(l, plan)

	var err error
	if len(plan.Script()) == 0 {
		err = l.DoString(LuaMap["default"])
	} else {
		err = l.DoString(string(plan.Script()))
	}

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = l.CallByParam(lua.P{
		Fn:      l.GetGlobal("run"),
		NRet:    1,
		Protect: false,
	}, lua.LString(plan.NodeID()))

	// 传递输入参数n=1
	if err != nil {
		panic(err)
	}

}

func (w Worker) LoadMod(l *lua.LState, plan task.Task) {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(w.State.Mqtt)
	rpc := &LRpc{
		MqttProxy: mqttProxy,
	}

	service := &LService{
		State: w.State,
		Plan: &Plan{
			Id:      plan.ID(),
			PlanLog: plan.LogID(),
			NodeId:  plan.NodeID(),
		},
	}

	l.SetGlobal("get_id", l.NewFunction(service.GetID))
	l.SetGlobal("get_msg", l.NewFunction(service.GetMsg))
	l.SetGlobal("get_status", l.NewFunction(service.GetStatus))

	l.SetGlobal("rpc_notify", l.NewFunction(rpc.notify))
	l.SetGlobal("rpc_async", l.NewFunction(rpc.asyncCall))
	l.SetGlobal("rpc_call", l.NewFunction(rpc.call))

	l.SetGlobal("call_service", l.NewFunction(callService))
	l.SetGlobal("filePoolService", lua.LString("FilePoolService"))

	l.SetGlobal("plan_id", lua.LString(plan.ID()))
	l.SetGlobal("plan_log_id", lua.LString(plan.LogID()))
}
