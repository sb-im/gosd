package luavm

import (
	"context"
	"time"

	"sb.im/gosd/rpc2mqtt"

	"sb.im/gosd/app/model"
	"sb.im/gosd/app/storage"

	jsonrpc "github.com/sb-im/jsonrpc-lite"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service struct {
	orm *gorm.DB
	rdb *redis.Client
	ofs *storage.Storage

	timeout time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	Rpc    *Rpc
	Task   *model.Task
	Server *rpc2mqtt.Rpc2mqtt
}

func NewService(task *model.Task) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		timeout: time.Hour,

		ctx:    ctx,
		cancel: cancel,
		Task:   task,
		Rpc:    NewRpc(),
	}
}

func (s Service) onStart() error {
	return s.running(s.Task)
}

func (s *Service) Close() error {
	s.cancel()

	for _, ch := range s.Rpc.pendings {
		rpc := jsonrpc.NewErrors("user.killed")
		rpc.Errors.InternalError("Be killed")
		data, err := rpc.ToJSON()
		if err != nil {
			return err
		}
		ch <- data
	}

	// Need Reset Kill status
	// Maybe multiple click Kill button
	// IOGets: need to sleep some time
	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s.onClose()
}

func (s *Service) onClose() error {
	// Clean up the "Dialog" when exiting
	s.CleanDialog()

	return s.running(&struct{}{})
}
