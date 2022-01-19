package luavm

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"sb.im/gosd/rpc2mqtt"

	lualib "sb.im/gosd/app/luavm/lua"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/storage"

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
	ctx    context.Context
	orm    *gorm.DB
	rdb    *redis.Client
	ofs    *storage.Storage
	script []byte
	mutex  *sync.Mutex

	instance string
	timeout  time.Duration

	rpc     *rpc2mqtt.Rpc2mqtt
	Queue   chan *model.Task
	Running map[string]*Service
}

func NewWorker(orm *gorm.DB, rdb *redis.Client, ofs *storage.Storage, rpc *rpc2mqtt.Rpc2mqtt, script []byte) *Worker {
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
		ctx:    context.TODO(),
		orm:    orm,
		rdb:    rdb,
		ofs:    ofs,
		script: script,
		mutex:  &sync.Mutex{},

		instance: "gosd.0",
		timeout:  time.Hour,

		rpc:     rpc,
		Queue:   make(chan *model.Task, 1024),
		Running: make(map[string]*Service),
	}
}

func (w Worker) RunTask(task *model.Task) error {
	if err := w.preTaskCheck(task); err != nil {
		return err
	}

	taskID := strconv.Itoa(int(task.ID))
	nodeID := strconv.Itoa(int(task.NodeID))

	// Lock
	w.lockTaskSet(taskID)
	w.lockNodeSet(nodeID)

	go func() {
		if err := w.doRun(task, w.getScript(task)); err != nil {
			log.Error(err)
		}

		// Unlock
		w.lockTaskDel(taskID)
		w.lockNodeDel(nodeID)
	}()

	return nil
}

func (w Worker) getScript(task *model.Task) (script []byte) {
	files := make(map[string]string)
	if err := json.Unmarshal(task.Files, &files); err == nil {
		if key, ok := files[defaultLuaTask]; ok {
			if data, err := w.ofs.Get(key); err == nil {
				script = data
			}
		}
	}

	// if nil, use system default
	if len(script) == 0 {
		script = w.script
	}

	return
}

func (w Worker) Run() {
	for task := range w.Queue {
		// Task Lua Script
		if err := w.doRun(task, w.getScript(task)); err != nil {
			log.Error(err)
		}
	}
}

func (w Worker) doRun(task *model.Task, script []byte) error {
	l := lua.NewState()
	defer func() {
		l.Close()

		if r := recover(); r != nil {
			log.Errorf("Emergency stop taskID: %d", task.ID)
			log.Error(r)
		}
	}()

	luajson.Preload(l)

	service := NewService(task)
	service.orm = w.orm
	service.rdb = w.rdb
	service.ofs = w.ofs
	service.Server = w.rpc

	w.mutex.Lock()
	w.Running[strconv.Itoa(int(task.ID))] = service
	w.mutex.Unlock()

	defer func() {
		w.mutex.Lock()
		delete(w.Running, strconv.Itoa(int(task.ID)))
		w.mutex.Unlock()
		log.Warn("==> luavm END")
	}()

	service.onStart()
	defer service.onClose()
	l.SetGlobal("SD", luar.New(l, service))

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
	if err := l.DoString(string(script)); err != nil {
		return err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("SD_main"),
		NRet:    1,
		Protect: true,
	}, lua.LString(strconv.Itoa(int(task.NodeID)))); err != nil {
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
