package luavm

import (
	"fmt"
	"strings"

	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/model"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

type Worker struct {
	RpcServer     *rpc2mqtt.Rpc2mqtt
	Queue         chan *Task
	Log           chan *model.PlanLog
	State         *state.State
	Store         *storage.Storage
	Running       map[string]*Service
}

func NewWorker(s *state.State, store *storage.Storage, rpcServer *rpc2mqtt.Rpc2mqtt) *Worker {
	return &Worker{
		Queue:         make(chan *Task, 1024),
		Log:           make(chan *model.PlanLog, 1024),
		State:         s,
		Store:         store,
		Running:       make(map[string]*Service),
		RpcServer:     rpcServer,
	}
}

func (w Worker) Run() {
	go func() {
		for l := range w.Log {
			w.SetRunning(l.PlanID, l)
		}
	}()

	for task := range w.Queue {
		if err := w.doRun(task); err != nil {
			fmt.Println(err)
		}
	}
}

func (w Worker) doRun(task *Task) error {
	var err error

	l := lua.NewState()
	planID := task.planID
	defer func() {
		l.Close()

		if r := recover(); r != nil {
			fmt.Printf("Emergency stop planID: %d\n", planID)
			fmt.Printf("Error：%s\n", r)
		}
		w.SetRunning(planID, &struct{}{})
	}()

	luajson.Preload(l)

	service := NewService(task)
	service.State = w.State
	service.Store = w.Store
	service.Server = w.RpcServer
	w.Running[task.PlanID()] = service
	defer delete(w.Running, task.PlanID())
	defer fmt.Println("==> luavm END")
	l.SetGlobal("SD", luar.New(l, service))

	// Clean up the "Dialog" when exiting
	defer service.CleanDialog()

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
	}, lua.LString(task.NodeID())); err != nil {
		return err
	}

	fmt.Println("==> luavm no panic")
	return nil
}

func (w Worker) Kill(planID string) {
	if service, ok := w.Running[planID]; ok {
		service.Close()
		fmt.Println("==> luavm Kill")
	}
}
