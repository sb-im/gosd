package luavm

import (
	"context"

	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	jsonrpc "github.com/sb-im/jsonrpc-lite"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
	Rpc    *Rpc
	Task   *Task
	State  *state.State
	Server *rpc2mqtt.Rpc2mqtt
	Store  *storage.Storage
}

func NewService(task *Task) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		ctx:    ctx,
		cancel: cancel,
		Task:   task,
		Rpc:    NewRpc(),
	}
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
	return nil
}
