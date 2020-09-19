package luavm

import (
	"fmt"
	"strings"

	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

type Worker struct {
	rpcs          map[string]*LRpc
	Queue         chan *Task
	State         *state.State
	Running       map[string]*Service
	MqttProxy     *jsonrpc2mqtt.MqttProxy
	StatusManager *StatusManager
}

func NewWorker(s *state.State) *Worker {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(s.Mqtt)
	return &Worker{
		rpcs:          make(map[string]*LRpc),
		Queue:         make(chan *Task, 1024),
		State:         s,
		Running:       make(map[string]*Service),
		MqttProxy:     mqttProxy,
		StatusManager: NewStatusManager(s.Mqtt),
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

	w.rpcs[task.PlanID] = NewLRpc(task, w.MqttProxy, w.StatusManager)
	fmt.Println(w.rpcs)

	service := NewService(task)
	service.Rpc.MqttProxy = w.MqttProxy
	service.State = w.State
	w.Running[task.PlanID] = service
	defer delete(w.Running, task.PlanID)

	// Clean up the "Dialog" when exiting
	defer service.CleanDialog()

	w.LoadMod(l, task)

	var err error
	if fn, err := l.Load(strings.NewReader(LuaMap["lib"]), "lib.lua"); err != nil {
		return err
	} else {
		l.Push(fn)
		err = l.PCall(0, lua.MultRet, nil)
	}

	if len(task.Script) == 0 {
		err = l.DoString(LuaMap["default"])
	} else {
		err = l.DoString(string(task.Script))
	}

	if err != nil {
		return err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("main"),
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

	if service, ok := w.Running[planID]; ok {
		service.Close()
	}

}

func (w Worker) LoadMod(l *lua.LState, task *Task) {
	rpc := w.rpcs[task.PlanID]
	service := w.Running[task.PlanID]

	l.SetGlobal("SD", luar.New(l, service))

	l.SetGlobal("rpc_notify", l.NewFunction(rpc.notify))
	l.SetGlobal("rpc_async", l.NewFunction(rpc.asyncCall))
	l.SetGlobal("rpc_call", l.NewFunction(rpc.call))
}
