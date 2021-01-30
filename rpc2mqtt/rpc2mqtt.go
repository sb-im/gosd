package rpc2mqtt

import (
	"context"
	"errors"
	"sync"

	"sb.im/gosd/mqttd"

	"github.com/sb-im/jsonrpc-lite"
)

type Rpc2mqtt struct {
	pending map[jsonrpc.ID]chan []byte
	mutex   sync.Mutex
	i       chan<- mqttd.MqttRpc
	o       <-chan mqttd.MqttRpc
}

func NewRpc2Mqtt(i chan<- mqttd.MqttRpc, o <-chan mqttd.MqttRpc) *Rpc2mqtt {
	return &Rpc2mqtt{
		i:       i,
		o:       o,
		pending: make(map[jsonrpc.ID]chan []byte),
	}
}

func (t *Rpc2mqtt) AsyncRpc(id string, req []byte, ch_res chan []byte) error {
	select {
	case t.i <- mqttd.MqttRpc{
		ID:      id,
		Payload: req,
	}:
		t.mutex.Lock()
		t.pending[*jsonrpc.ParseObject(req).ID] = ch_res
		t.mutex.Unlock()
	default:
		return errors.New("Buffer full")
	}
	return nil
}

func (t *Rpc2mqtt) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case raw := <-t.o:
			rpc := jsonrpc.ParseObject(raw.Payload)
			if pending := t.pending[*rpc.ID]; pending != nil && (rpc.Type == jsonrpc.TypeSuccess || rpc.Type == jsonrpc.TypeErrors) {
				t.mutex.Lock()
				delete(t.pending, *rpc.ID)
				t.mutex.Unlock()
				pending <- raw.Payload
			}
		}
	}
}
