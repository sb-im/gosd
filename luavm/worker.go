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
	Queue         chan *Task
	State         *state.State
	Running       map[string]*Service
	MqttProxy     *jsonrpc2mqtt.MqttProxy
	StatusManager *StatusManager
}

func NewWorker(s *state.State) *Worker {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(s.Mqtt)
	return &Worker{
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

	service := NewService(task)
	service.Rpc.MqttProxy = w.MqttProxy
	service.State = w.State
	w.Running[task.PlanID] = service
	defer delete(w.Running, task.PlanID)
	defer fmt.Println("==> luavm END")
	l.SetGlobal("SD", luar.New(l, service))

	// Clean up the "Dialog" when exiting
	defer service.CleanDialog()

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
		Fn:      l.GetGlobal("SD_main"),
		NRet:    1,
		Protect: true,
	}, lua.LString(task.NodeID)); err != nil {
		return err
	}

	fmt.Println("==> luavm no panic")
	return nil
}

func (w Worker) Kill(planID string) {
	if service, ok := w.Running[planID]; ok {
		service.Close()
		fmt.Println("==> luavm Kill")
		w.StatusManager.SetStatus(planID, StatusProtect)
	}
}
