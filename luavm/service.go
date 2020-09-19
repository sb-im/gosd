package luavm

import (
	"context"

	"sb.im/gosd/state"

	jsonrpc "github.com/sb-im/jsonrpc-lite"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
	Rpc    *Rpc
	Task   *Task
	State  *state.State
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

	for id, ch := range s.Rpc.pendings {
		rpc := jsonrpc.NewErrors(id)
		rpc.Errors.InternalError("Be killed")
		data, err := rpc.ToJSON()
		if err != nil {
			return err
		}
		ch <- data
	}
	return nil
}
