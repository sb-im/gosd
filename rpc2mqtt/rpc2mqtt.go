package rpc2mqtt

import (
	"context"
	"errors"
	"time"

	"sb.im/gosd/mqttd"

	"github.com/jpillora/backoff"
	"github.com/sb-im/jsonrpc-lite"
	log "github.com/sirupsen/logrus"
)

var (
	// Every 0.1s Detect Send && Resend
	loopWait = 100 * time.Millisecond

	// Max Resend interval
	maxRetryTime   = time.Minute
	commandTimeout = 90 * time.Minute
)

type Pending struct {
	Dst   string
	Msg   []byte
	Time  time.Time
	Reply chan []byte

	nextTime time.Time
	backoff  *backoff.Backoff
}

type Rpc2mqtt struct {
	pending map[jsonrpc.ID]*Pending
	cache   chan Pending
	kill    chan []byte
	i       chan<- mqttd.MqttRpc
	o       <-chan mqttd.MqttRpc
}

func NewRpc2Mqtt(i chan<- mqttd.MqttRpc, o <-chan mqttd.MqttRpc) *Rpc2mqtt {
	return &Rpc2mqtt{
		i:       i,
		o:       o,
		kill:    make(chan []byte, 1024),
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
		case raw := <-t.kill:
			req := jsonrpc.ParseObject(raw)

			rpc := jsonrpc.NewErrors(req.ID)
			rpc.Errors.InternalError("Be killed")
			data, err := rpc.ToJSON()
			if err != nil {
				log.Error(err)
			}

			if rpc.Type == jsonrpc.TypeInvalid {
				continue
			}
			if pending, ok := t.pending[*rpc.ID]; ok && (rpc.Type == jsonrpc.TypeSuccess || rpc.Type == jsonrpc.TypeErrors) {
				delete(t.pending, *rpc.ID)
				pending.Reply <- data
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
				if now := time.Now(); now.After(pending.nextTime) && now.Before(pending.Time.Add(commandTimeout)) {
					select {
					case t.i <- mqttd.MqttRpc{
						ID:      pending.Dst,
						Payload: pending.Msg,
					}:
						pending.nextTime = pending.nextTime.Add(pending.backoff.Duration())
						log.Tracef("SEND: %s", pending.Msg)

					default:
						log.Error("Buffer full")
					}
				} else if now.After(pending.Time.Add(commandTimeout)) {
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
		Reply: ch_res,

		nextTime: time.Now().Add(time.Second),
		backoff: &backoff.Backoff{
			Factor: 2,
			Min:    2 * time.Second,
			Max:    maxRetryTime,
		},
	}:
	default:
		return errors.New("Send cache buffer full")
	}
	return nil
}

func (t *Rpc2mqtt) KillRpc(req []byte) error {
	select {
	case t.kill <- req:
	default:
		return errors.New("Kill cache buffer full")
	}
	return nil
}
