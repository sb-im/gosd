package luavm

import (
	"context"
	"encoding/json"
	"errors"

	"sb.im/gosd/jsonrpc2mqtt"

	jsonrpc "github.com/sb-im/jsonrpc-lite"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
	Rpc    *Rpc
}

type Rpc struct {
	MqttProxy *jsonrpc2mqtt.MqttProxy
}

func NewService() *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		ctx:    ctx,
		cancel: cancel,
		Rpc:    &Rpc{},
	}
}

func (s *Service) Close() {
	s.cancel()
}

func (s *Service) SyncCall(id, method string, params []byte) (string, error) {
	rpc := jsonrpc.NewRequest(id, method, nil)
	p := json.RawMessage(params)
	rpc.Params = &p
	req, err := rpc.ToJSON()
	if err != nil {
		return string(req), err
	}

	var res []byte

	// TODO: Detect online status
	ch := make(chan []byte)

	err = s.Rpc.MqttProxy.AsyncRpc(id, req, ch)
	if err != nil {
		return "", err
	}

	select {
	case <-s.ctx.Done():
		return "", errors.New("Be killed")
	case res = <-ch:
		return string(res), nil
	}
}
