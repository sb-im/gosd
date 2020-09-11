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
	rpcs      map[string]*LRpc
	Queue     chan *Task
	State     *state.State
	Runtime   map[string]*lua.LState
	MqttProxy *jsonrpc2mqtt.MqttProxy
	StatusManager *StatusManager
}

func NewWorker(s *state.State) *Worker {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(s.Mqtt)
	return &Worker{
		rpcs:      make(map[string]*LRpc),
		Queue:     make(chan *Task, 1024),
		State:     s,
		Runtime:   make(map[string]*lua.LState),
		MqttProxy: mqttProxy,
		StatusManager: NewStatusManager(&s.Mqtt),
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
	planID := task.PlanID
	w.StatusManager.SetStatus(planID, StatusRunning)
	defer func() {
		if w.StatusManager.GetStatus(planID) == StatusRunning {
			l.Close()
		}

		if r := recover(); r != nil {
			fmt.Printf("Emergency stop planID: %s\n", planID)
			fmt.Printf("Error：%s\n", r)
		}
		w.StatusManager.SetStatus(planID, StatusReady)
	}()

	luajson.Preload(l)

	w.rpcs[task.PlanID] = NewLRpc(w.MqttProxy)
	fmt.Println(w.rpcs)
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

	w.Runtime[task.PlanID] = l
	defer delete(w.Runtime, task.PlanID)

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("run"),
		NRet:    1,
		Protect: true,
	}, lua.LString(task.NodeID)); err != nil {
		return err
	}

	fmt.Println("==> luavm END")
	return nil
}

func (w Worker) Kill(planID string) {
	//l, ok := w.Runtime[planID]
	r, ok := w.rpcs[planID]
	if ok {
		fmt.Println("==> luavm Close")
		//l.Close()
		if err := r.Kill(); err != nil {
			fmt.Println(err)
		}
		w.StatusManager.SetStatus(planID, StatusProtect)
	}
}

func (w Worker) LoadMod(l *lua.LState, task *Task) {
	rpc := w.rpcs[task.PlanID]
	fmt.Println(rpc)

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
