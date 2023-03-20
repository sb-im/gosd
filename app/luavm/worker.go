package luavm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"sb.im/gosd/rpc2mqtt"

	"sb.im/gosd/app/logger"
	"sb.im/gosd/app/luavm/lib"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	luar "layeh.com/gopher-luar"
)

type Worker struct {
	cfg Config
	ctx context.Context
	srv *service.Service
	orm *gorm.DB
	rdb *redis.Client
	ofs *storage.Storage

	script []byte
	mutex  *sync.Mutex

	rpc     *rpc2mqtt.Rpc2mqtt
	Running map[string]*Service
}

func NewWorker(cfg Config, srv *service.Service, rpc *rpc2mqtt.Rpc2mqtt, script []byte) *Worker {
	ctx := context.TODO()
	// default LuaFile: input > default
	if len(script) == 0 {
		if data, err := lib.File.ReadFile(defaultFileLua); err != nil {
			logger.WithContext(ctx).Error(err)
		} else {
			script = data
		}
	}

	return &Worker{
		cfg: cfg,
		ctx: ctx,
		srv: srv,
		orm: srv.Orm(),
		rdb: srv.Rdb(),
		ofs: srv.Ofs(),

		script: script,
		mutex:  &sync.Mutex{},

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
		if key, ok := files[defaultFileKey]; ok {
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

func (w *Worker) Run(ctx context.Context) {
	chRun := w.rdb.Subscribe(ctx, fmt.Sprintf("__keyevent@%d__:expired", w.rdb.Options().DB)).Channel()
	chEnd := w.rdb.Subscribe(ctx, "luavm.kill").Channel()
	for {
		select {
		case m := <-chRun:
			// <prefix>.<jobId>
			// job.1
			if data := strings.Split(m.Payload, "."); len(data) > 1 {
				if data[0] == "job" {
					jobId := data[1]
					var job model.Job
					if err := w.orm.Preload("Task").Find(&job, jobId).Error; err != nil {
						logger.WithContext(ctx).Errorln(err)
					}

					logger.WithContext(ctx).Infof("%+v", job)

					task := job.Task
					task.Job = &job

					files := make(map[string]string)
					extra := make(map[string]string)

					json.Unmarshal(task.Files, &files)
					json.Unmarshal(task.Extra, &extra)

					logger.WithContext(ctx).Infof("%+v\t%v\t%v", task, files, extra)

					w.AddTask(ctx, &task)
				}
			}
		case m := <-chEnd:
			w.Kill(m.Payload)
		case <-ctx.Done():
		}
	}
}

func (w *Worker) RunTask(task *model.Task, script []byte) error {
	return w.doRun(context.Background(), task, script)
}

func (w *Worker) doRun(ctx context.Context, task *model.Task, script []byte) error {
	w.setDuration(task.Job, 1)
	l := lua.NewState()
	defer func() {
		l.Close()

		if r := recover(); r != nil {
			logger.WithContext(ctx).Errorf("Emergency stop taskID: %d", task.ID)
			logger.WithContext(ctx).Error(r)
		}
	}()

	luajson.Preload(l)

	service := NewService(ctx, task)
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

	// Patch
	l.SetGlobal("print", l.NewFunction(patchBasePrint(ctx)))

	// Main
	l.SetGlobal("SD", luar.New(l, service))

	// Load Lib
	for _, name := range libs {
		if f, err := lib.File.Open(name); err != nil {
			logger.WithContext(ctx).Error(err)
			continue
		} else {
			if fn, err := l.Load(f, name); err != nil {
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
		w.setDuration(task.Job, -1)
		return err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("SD_main"),
		NRet:    1,
		Protect: true,
	}, lua.LString(node.UUID)); err != nil {
		// Min duration 3s
		duration := int(time.Since(task.Job.StartedAt).Seconds())
		if duration < 3 {
			duration = 3
		}
		w.setDuration(task.Job, -duration)
		return err
	}

	logger.WithContext(ctx).Warn("==> luavm no panic")

	// Min duration 3s
	duration := int(time.Since(task.Job.StartedAt).Seconds())
	if duration < 3 {
		duration = 3
	}
	return w.setDuration(task.Job, duration)
}

func (w *Worker) setDuration(job *model.Job, duration int) error {
	return w.orm.Model(job).Where("id = ?", job.ID).UpdateColumn("duration", duration).Error
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
