package luavm

import (
	"strings"

	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	log "github.com/sirupsen/logrus"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

type Worker struct {
	defaultLua    []byte
	RpcServer     *rpc2mqtt.Rpc2mqtt
	Queue         chan *Task
	State         *state.State
	Store         *storage.Storage
	Running       map[string]*Service
}

func NewWorker(s *state.State, store *storage.Storage, rpcServer *rpc2mqtt.Rpc2mqtt, defaultLua []byte) *Worker {
	if len(defaultLua) == 0 {
		defaultLua = []byte(LuaMap["default"])
	}
	return &Worker{
		Queue:         make(chan *Task, 1024),
		State:         s,
		Store:         store,
		Running:       make(map[string]*Service),
		RpcServer:     rpcServer,
		defaultLua:    defaultLua,
	}
}

func (w Worker) Run() {
	for task := range w.Queue {
		if err := w.doRun(task); err != nil {
			log.Error(err)
		}
	}
}

func (w Worker) doRun(task *Task) error {
	w.SetRunning(task.PlanID, task)
	var err error

	l := lua.NewState()
	defer func() {
		l.Close()

		if r := recover(); r != nil {
			log.Errorf("Emergency stop planID: %d\n", task.PlanID)
			log.Errorf("Error：%s\n", r)
		}
		w.SetRunning(task.PlanID, &struct{}{})
	}()

	luajson.Preload(l)

	service := NewService(task)
	service.State = w.State
	service.Store = w.Store
	service.Server = w.RpcServer
	w.Running[task.StringPlanID()] = service
	defer delete(w.Running, task.StringPlanID())
	defer log.Warn("==> luavm END")
	l.SetGlobal("SD", luar.New(l, service))

	// Clean up the "Dialog" when exiting
	defer service.CleanDialog()

	// Load Lib
	for _, lib := range []string{"lib_plan", "lib_node", "lib_geo", "lib_log", "lib_main"} {
		if fn, err := l.Load(strings.NewReader(LuaMap[lib]), lib + ".lua"); err != nil {
			return err
		} else {
			l.Push(fn)
			err = l.PCall(0, lua.MultRet, nil)
		}
	}

	if len(task.Script) == 0 {
		err = l.DoString(string(w.defaultLua))
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
	}, lua.LString(task.StringNodeID())); err != nil {
		return err
	}

	log.Warn("==> luavm no panic")
	return nil
}

func (w Worker) Kill(planID string) {
	if service, ok := w.Running[planID]; ok {
		service.Close()
		log.Warn("==> luavm Kill")
	}
}

func (w Worker) Close() {
	// Stop Create running
	close(w.Queue)

	for _, service := range w.Running {
		service.Close()
	}

	log.Warn("==> luaVM Worker Close")
}
