package luavm

import (
	"fmt"
	"errors"

	redis "github.com/gomodule/redigo/redis"
)

const (
	topic_terminal = "plans/%s/term"
)

func (s *Service) IOGets() (string, error) {
	ch := make(chan []byte)
	topic := fmt.Sprintf(topic_terminal, s.Task.PlanID)

	go func() {
		keyspace := "__keyspace@0__:%s"
		psc := redis.PubSubConn{Conn: s.State.Pool.Get()}
		psc.PSubscribe(fmt.Sprintf(keyspace, topic))
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)

				raw, err := s.State.BytesGet(topic)
				if err != nil {
					// TODO: error handling
				}
				ch <- raw
				return
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
				//return v
			default:
				fmt.Println("default")
			}
		}
	}()

	var raw []byte
	select {
	case <-s.ctx.Done():
		return "", errors.New("Be killed")
	case raw = <-ch:
	}

	return string(raw), nil
}

func (s *Service) IOPuts(str string) error {
	return s.State.Record(fmt.Sprintf(topic_terminal, s.Task.PlanID), []byte(str))
}
