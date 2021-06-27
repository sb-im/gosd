package rpc2mqtt

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"sb.im/gosd/mqttd"

	log "github.com/sirupsen/logrus"
	"github.com/sb-im/jsonrpc-lite"
)

type Pending struct {
	Dst   string
	Msg   []byte
	Time  time.Time
	Reply chan []byte
	Count int
}

type Rpc2mqtt struct {
	pending map[jsonrpc.ID]*Pending
	mutex   sync.Mutex
	i       chan<- mqttd.MqttRpc
	o       <-chan mqttd.MqttRpc
}

func NewRpc2Mqtt(i chan<- mqttd.MqttRpc, o <-chan mqttd.MqttRpc) *Rpc2mqtt {
	return &Rpc2mqtt{
		i:       i,
		o:       o,
		pending: make(map[jsonrpc.ID]*Pending),
	}
}

func (t *Rpc2mqtt) AsyncRpc(id string, req []byte, ch_res chan []byte) error {
	select {
	case t.i <- mqttd.MqttRpc{
		ID:      id,
		Payload: req,
	}:
		t.mutex.Lock()
		t.pending[*jsonrpc.ParseObject(req).ID] = &Pending{
			Time:  time.Now(),
			Dst:   id,
			Msg:   req,
			Count: 1,
			Reply: ch_res,
		}
		t.mutex.Unlock()
	default:
		return errors.New("Buffer full")
	}
	return nil
}

func (t *Rpc2mqtt) Run(ctx context.Context) {
	go t.Resend(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case raw := <-t.o:
			log.Tracef("RECV: %s", raw.Payload)
			rpc := jsonrpc.ParseObject(raw.Payload)
			if rpc.Type == jsonrpc.TypeInvalid {
				continue
			}
			if pending, ok := t.pending[*rpc.ID]; ok && (rpc.Type == jsonrpc.TypeSuccess || rpc.Type == jsonrpc.TypeErrors) {
				log.Debugf("res: %s", raw.Payload)
				t.mutex.Lock()
				delete(t.pending, *rpc.ID)
				t.mutex.Unlock()
				pending.Reply <- raw.Payload
			}
		}
	}
}

func (t *Rpc2mqtt) Resend(ctx context.Context) {
	// Max Resend interval 1h
	timeout := float64(3600)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for id, pending := range t.pending {
				if now := time.Since(pending.Time).Seconds(); now > math.Ldexp(1, pending.Count) && now <= timeout {
					select {
					case t.i <- mqttd.MqttRpc{
						ID:      pending.Dst,
						Payload: pending.Msg,
					}:
						pending.Count++

					default:
						log.Error("Buffer full")
					}
				} else if now > timeout {
					rpc := jsonrpc.NewError(id, 1, "timeout", nil)
					data, _ := rpc.ToJSON()
					t.mutex.Lock()
					delete(t.pending, *rpc.ID)
					t.mutex.Unlock()
					pending.Reply <- data
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
