package luavm

import (
	"context"
	"errors"
	"strconv"
	"time"

	"sb.im/gosd/jsonrpc2mqtt"
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

type Rpc struct {
	pendings  map[string]chan []byte
	MqttProxy *jsonrpc2mqtt.MqttProxy
}

func NewRpc() *Rpc {
	return &Rpc{
		pendings: make(map[string]chan []byte),
	}
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

func (s *Service) GenRpcID() string {
	bit13_timestamp := string([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))[:13])
	return "gosd.0-" + bit13_timestamp + "-" + getSequence()
}

func (s *Service) RpcSend(nodeId string, raw []byte) (string, error) {
	// TODO: Detect online status
	rpc, err := jsonrpc.Parse(raw)
	if err != nil {
		return "", err
	}
	ch := make(chan []byte, 128)
	s.Rpc.pendings[rpc.ID.String()] = ch

	// Prevent issuing non-compliant jsonrpc 2.0
	req, err := rpc.ToJSON()
	if err != nil {
		return "", err
	}

	err = s.Rpc.MqttProxy.AsyncRpc(nodeId, req, ch)
	if err != nil {
		return "", err
	}
	return rpc.ID.String(), nil
}

func (s *Service) RpcRecv(id string) (string, error) {
	if ch, ok := s.Rpc.pendings[id]; ok {

		var raw []byte
		select {
		case <-s.ctx.Done():
			return "", errors.New("Be killed")
		case raw = <-ch:
		}

		delete(s.Rpc.pendings, jsonrpc.ParseObject(raw).ID.String())
		return string(raw), nil
	}
	return "", errors.New("Not pending")
}
