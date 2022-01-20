package luavm

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	topic_terminal = "tasks/%d/term"
)

func (s Service) IOGets() (string, error) {
	log.Debug("IOGets")
	topic := fmt.Sprintf(topic_terminal, s.Task.ID)

	keyspace := "__keyspace@*__:%s"
	pubsub := s.rdb.PSubscribe(s.ctx, fmt.Sprintf(keyspace, topic))
	ch2 := pubsub.Channel()
	m := <-ch2
	pubsub.Close()
	log.Debug(m)
	return s.rdb.Get(s.ctx, topic).Result()
}

func (s Service) IOPuts(str string) error {
	return s.rdb.Set(context.Background(), fmt.Sprintf(topic_terminal, s.Task.ID), str, s.timeout).Err()
}
