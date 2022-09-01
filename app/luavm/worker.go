package luavm

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"sb.im/gosd/rpc2mqtt"

	"sb.im/gosd/app/logger"
	lualib "sb.im/gosd/app/luavm/lua"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/storage"
	"sb.im/gosd/app/store"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

var (
	libs = []string{
		"lib_task.lua",
		"lib_node.lua",
		"lib_geo.lua",
		"lib_log.lua",
		"lib_main.lua",
	}
)

type Worker struct {
	cfg    Config
	ctx    context.Context
	orm    *gorm.DB
	rdb    *redis.Client
	ofs    *storage.Storage
	store  *store.Store
	script []byte
	mutex  *sync.Mutex

	timeout time.Duration

	rpc     *rpc2mqtt.Rpc2mqtt
	Running map[string]*Service
}

func NewWorker(cfg Config, s *store.Store, rpc *rpc2mqtt.Rpc2mqtt, script []byte) *Worker {
	ctx := context.TODO()
	// default LuaFile: input > default
	if len(script) == 0 {
		if data, err := lualib.LuaFile.ReadFile(cfg.LuaFile); err != nil {
			logger.WithContext(ctx).Error(err)
		} else {
			script = data
		}
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		logger.WithContext(ctx).Error(err)
		timeout = 2 * time.Hour
	}

	return &Worker{
		cfg:    cfg,
		ctx:    ctx,
		orm:    s.Orm(),
		rdb:    s.Rdb(),
		ofs:    s.Ofs(),
		store:  s,
		script: script,
		mutex:  &sync.Mutex{},

		timeout: timeout,

		rpc:     rpc,
		Running: make(map[string]*Service),
	}
}

func (w *Worker) AddTask(ctx context.Context, task *model.Task) error {
	if err := w.preTaskCheck(task); err != nil {
		return err
	}

	taskID := strconv.Itoa(int(task.ID))
	nodeID := strconv.Itoa(int(task.NodeID))

	// Lock
	w.lockTaskSet(taskID)
	w.lockNodeSet(nodeID)

	go func() {
		if err := w.doRun(ctx, task, w.getScript(task)); err != nil {
			logger.WithContext(ctx).Error(err)
		}

		// Unlock
		w.lockTaskDel(taskID)
		w.lockNodeDel(nodeID)
	}()

	return nil
}

func (w *Worker) getScript(task *model.Task) (script []byte) {
	files := make(map[string]string)
	if err := json.Unmarshal(task.Files, &files); err == nil {
		if key, ok := files[w.cfg.LuaTask]; ok {
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

func (w Worker) Run(ctx context.Context) {
	<-ctx.Done()
}

func (w Worker) RunTask(task *model.Task, script []byte) error {
	return w.doRun(context.Background(), task, script)
}

func (w *Worker) doRun(ctx context.Context, task *model.Task, script []byte) error {
	l := lua.NewState()
	defer func() {
		l.Close()

		if r := recover(); r != nil {
			logger.WithContext(ctx).Errorf("Emergency stop taskID: %d", task.ID)
			logger.WithContext(ctx).Error(r)
		}
	}()

	luajson.Preload(l)

	service := NewService(task)
	var nodes []model.Node
	w.orm.Find(&nodes, "team_id = ?", task.TeamID)
	service.cfg = w.cfg
	service.orm = w.orm
	service.rdb = w.rdb
	service.ofs = w.ofs
	service.nodes = nodes
	service.Server = w.rpc

	var node model.Node
	for _, n := range nodes {
		if task.NodeID == n.ID {
			node = n
		}
	}
	logger.WithContext(ctx).Debugf("Task: %+v", task)
	logger.WithContext(ctx).Debugf("Node: %+v", node)
	if node.ID == 0 {
		// TODO: unit test need change
		//return errors.New("Not Found This Node: " + task.NodeID)
	}

	w.mutex.Lock()
	w.Running[strconv.Itoa(int(task.ID))] = service
	w.mutex.Unlock()

	defer func() {
		w.mutex.Lock()
		delete(w.Running, strconv.Itoa(int(task.ID)))
		w.mutex.Unlock()
		logger.WithContext(ctx).Warn("==> luavm END")
	}()

	service.onStart()
	defer service.onClose()
	l.SetGlobal("SD", luar.New(l, service))

	// Load Lib
	for _, lib := range libs {
		if f, err := lualib.LuaFile.Open(lib); err != nil {
			logger.WithContext(ctx).Error(err)
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

	logger.WithContext(ctx).Debug(string(script))

	// Core: Load Script
	if err := l.DoString(string(script)); err != nil {
		return err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("SD_main"),
		NRet:    1,
		Protect: true,
	}, lua.LString(node.UUID)); err != nil {
		return err
	}

	logger.WithContext(ctx).Warn("==> luavm no panic")
	return nil
}

func (w *Worker) Kill(taskID string) error {
	w.mutex.Lock()
	service, ok := w.Running[taskID]
	w.mutex.Unlock()
	if ok {
		service.Kill()
		logger.WithContext(w.ctx).Warn("==> luavm Kill")
		return nil
	}
	return errors.New("Not Found This Task")
}

func (w *Worker) Close() {
	w.mutex.Lock()
	for _, service := range w.Running {
		service.Close()
	}
	w.mutex.Unlock()

	logger.WithContext(w.ctx).Warn("==> luaVM Worker Close")
}
