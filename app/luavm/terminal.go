package luavm

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	//redis "github.com/gomodule/redigo/redis"
	//"github.com/go-redis/redis/v8"
)

const (
	topic_terminal = "tasks/%d/term"
)

func (s Service) IOGets() (string, error) {
	log.Debug("IOGets")
	ch := make(chan string)
	topic := fmt.Sprintf(topic_terminal, s.Task.ID)

	// TODO: Need test ctx
	// Or: go func() {}
	// return "", errors.New("Be killed")
	//ch := make(chan string, 1)
	//err := s.rdb.Watch(context.Background(), func(t *redis.Tx) error {
	//	log.Debug("On Watch")
	//	ch <- t.Get(context.Background(), topic).String()
	//	return nil
	//}, topic)

	//go func() {
	//	s.rdb.Watch(context.Background(), func(t *redis.Tx) error {
	//		log.Debug("On Watch")
	//		ch <- t.Get(context.Background(), topic).String()
	//		return nil
	//	}, topic)
	//}()
	go func() {
		keyspace := "__keyspace@1__:%s"
		pubsub := s.rdb.Subscribe(context.Background(), fmt.Sprintf(keyspace, topic))
		ch2 := pubsub.Channel()
		m := <-ch2
		log.Debug(m)
		ch <- s.rdb.Get(context.Background(), topic).String()
	}()

	var raw string
	select {
	case <-s.ctx.Done():
		return "", errors.New("Be killed")
	case raw = <-ch:
	}
	return raw, nil
}

func (s Service) IOPuts(str string) error {
	return s.rdb.Set(context.Background(), fmt.Sprintf(topic_terminal, s.Task.ID), str, 0).Err()
}
