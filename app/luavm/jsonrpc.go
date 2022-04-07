package luavm

import (
	"errors"
	"strconv"
	"sync"
	"time"

	jsonrpc "github.com/sb-im/jsonrpc-lite"
)

type Rpc struct {
	pendings map[string]tRpc
}

type tRpc struct {
	req []byte
	ch  chan []byte
}

func NewRpc() *Rpc {
	return &Rpc{
		pendings: make(map[string]tRpc),
	}
}

var sequence uint64
var sequenceMutex sync.Mutex

func getSequence() string {
	sequenceMutex.Lock()
	id := strconv.FormatUint(sequence, 10)
	sequence++
	sequenceMutex.Unlock()
	return id
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
	s.Rpc.pendings[rpc.ID.String()] = tRpc{
		req: raw,
		ch:  ch,
	}

	// Prevent issuing non-compliant jsonrpc 2.0
	req, err := rpc.ToJSON()
	if err != nil {
		return "", err
	}

	// Server.AsyncRpc is mqtt v5 jsonrpc over mqtt
	err = s.Server.AsyncRpc(nodeId, req, ch)
	if err != nil {
		return "", err
	}
	return rpc.ID.String(), nil
}

func (s *Service) RpcRecv(id string) (string, error) {
	if trpc, ok := s.Rpc.pendings[id]; ok {
		select {
		case <-s.ctx.Done():
			return "", errors.New("Be killed")
		case raw := <-trpc.ch:
			delete(s.Rpc.pendings, jsonrpc.ParseObject(raw).ID.String())
			return string(raw), nil
		}
	}
	return "", errors.New("Not pending")
}
