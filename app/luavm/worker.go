package luavm

import (
	"context"
	"strconv"

	lualib "sb.im/gosd/app/luavm/lua"
	"sb.im/gosd/app/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

const (
	//defaultLuaFile = "default.lua"
	defaultLuaFile = "test_min.lua"
	defaultLuaTask = "lua"
)

var (
	libs = []string{
		"lib_plan.lua",
		"lib_node.lua",
		"lib_geo.lua",
		"lib_log.lua",
		"lib_main.lua",
	}
)

type Worker struct {
	orm    *gorm.DB
	rdb    *redis.Client
	script []byte

	Queue   chan *model.Task
	Running map[string]*Service
}

func NewWorker(orm *gorm.DB, rdb *redis.Client, script []byte) *Worker {
	// default LuaFile: input > default
	if len(script) == 0 {
		if data, err := lualib.LuaFile.ReadFile(defaultLuaFile); err != nil {
			log.Error(err)
		} else {
			script = data
		}
	}

	// Enable Redis Events
	rdb.ConfigSet(context.Background(), "notify-keyspace-events", "$K")

	return &Worker{
		orm:    orm,
		rdb:    rdb,
		script: script,

		Queue:   make(chan *model.Task, 1024),
		Running: make(map[string]*Service),
	}
}

func (w Worker) Run() {
	for task := range w.Queue {
		// files := make(map[string]string)
		// json.Unmarshal(task.Files, &files)
		// key := files[defaultLuaTask]
		// TODO:
		// script := ???
		script := []byte{}

		//if len(script) == 0 {
		//  script = w.defaultLua
		//}

		if err := w.doRun(task, script); err != nil {
			log.Error(err)
		}
	}
}

func (w Worker) doRun(task *model.Task, script []byte) error {
	w.SetRunning(task.ID, task)
	var err error

	l := lua.NewState()
	defer func() {
		l.Close()

		if r := recover(); r != nil {
			log.Errorf("Emergency stop planID: %d\n", task.ID)
			log.Errorf("Error：%s\n", r)
		}
		w.SetRunning(task.ID, &struct{}{})
	}()

	luajson.Preload(l)

	service := NewService(task)
	service.orm = w.orm
	service.rdb = w.rdb
	w.Running[strconv.Itoa(int(task.ID))] = service
	defer delete(w.Running, strconv.Itoa(int(task.ID)))
	defer log.Warn("==> luavm END")
	l.SetGlobal("SD", luar.New(l, service))

	// Clean up the "Dialog" when exiting
	defer service.CleanDialog()

	// Load Lib
	for _, lib := range libs {
		if f, err := lualib.LuaFile.Open(lib); err != nil {
			log.Error(err)
			continue
		} else {
			if fn, err := l.Load(f, lib); err != nil {
				return err
			} else {
				l.Push(fn)
				err = l.PCall(0, lua.MultRet, nil)
			}
		}
	}

	// Core: Load Script
	err = l.DoString(string(script))
	if err != nil {
		return err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("SD_main"),
		NRet:    1,
		Protect: true,
	}, lua.LString("task.StringNodeID()")); err != nil {
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
