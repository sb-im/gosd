package rpc2mqtt

import (
	"context"
	"errors"
	"math"
	"time"

	"sb.im/gosd/mqttd"

	"github.com/sb-im/jsonrpc-lite"
	log "github.com/sirupsen/logrus"
)

var (
	// Every 0.1s Detect Send && Resend
	loopWait = 100 * time.Millisecond

	// Max Resend interval 1h
	timeout = float64(3600)
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
	cache   chan Pending
	i       chan<- mqttd.MqttRpc
	o       <-chan mqttd.MqttRpc
}

func NewRpc2Mqtt(i chan<- mqttd.MqttRpc, o <-chan mqttd.MqttRpc) *Rpc2mqtt {
	return &Rpc2mqtt{
		i:       i,
		o:       o,
		cache:   make(chan Pending, 1024),
		pending: make(map[jsonrpc.ID]*Pending),
	}
}

func (t *Rpc2mqtt) Run(ctx context.Context) {
	ticker := time.NewTicker(loopWait)
	defer ticker.Stop()
	for {
		select {
		case p := <-t.cache:
			// Send
			select {
			case t.i <- mqttd.MqttRpc{
				ID:      p.Dst,
				Payload: p.Msg,
			}:
				t.pending[*jsonrpc.ParseObject(p.Msg).ID] = &p
			default:
				log.Error("Buffer full")
			}
		case raw := <-t.o:
			// Recv
			log.Tracef("RECV: %s", raw.Payload)
			rpc := jsonrpc.ParseObject(raw.Payload)
			if rpc.Type == jsonrpc.TypeInvalid {
				continue
			}
			if pending, ok := t.pending[*rpc.ID]; ok && (rpc.Type == jsonrpc.TypeSuccess || rpc.Type == jsonrpc.TypeErrors) {
				log.Debugf("res: %s", raw.Payload)
				delete(t.pending, *rpc.ID)
				pending.Reply <- raw.Payload
			}
		case <-ticker.C:
			// Resend
			for id, pending := range t.pending {
				if now := time.Since(pending.Time).Seconds(); now > math.Ldexp(1, pending.Count) && now <= timeout {
					select {
					case t.i <- mqttd.MqttRpc{
						ID:      pending.Dst,
						Payload: pending.Msg,
					}:
						pending.Count++
						log.Tracef("SEND: %s", pending.Msg)

					default:
						log.Error("Buffer full")
					}
				} else if now > timeout {
					rpc := jsonrpc.NewError(id, 1, "timeout", nil)
					data, _ := rpc.ToJSON()
					delete(t.pending, *rpc.ID)
					pending.Reply <- data
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (t *Rpc2mqtt) AsyncRpc(id string, req []byte, ch_res chan []byte) error {
	select {
	case t.cache <- Pending{
		Time:  time.Now(),
		Dst:   id,
		Msg:   req,
		Count: 1,
		Reply: ch_res,
	}:
	default:
		return errors.New("Send cache buffer full")
	}
	return nil
}
