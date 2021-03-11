package state

import (
	"fmt"

	redis "github.com/gomodule/redigo/redis"
)

type State struct {
	Plan map[int]PlanState
	Pool *redis.Pool
}

type PlanState struct {
	id int64
}

func NewState(rawURL string) *State {
	c, err := redis.DialURL(rawURL)
	if err != nil {
		// TODO: handle connection error
		panic(err)
	}
	defer c.Close()
	c.Do("CONFIG", "SET", "notify-keyspace-events", "K$")

	return &State{
		Pool: redis.NewPool(func() (redis.Conn, error) { return redis.DialURL(rawURL) }, 5),
	}
}

func (s *State) do(commandName string, args ...interface{}) (interface{}, error) {
	return s.Pool.Get().Do(commandName, args...)
}

func (s *State) StringGet(key string) (string, error) {
	return redis.String(s.do("GET", key))
}

func (s *State) BytesGet(key string) ([]byte, error) {
	return redis.Bytes(s.do("GET", key))
}

func (s *State) Record(key string, value []byte) error {
	if _, err := s.do("SET", key, value); err != nil {
		return err
	}
	return nil
}

func (s *State) GetNode(id, msg string) (string, error) {
	return redis.String(s.do("GET", fmt.Sprintf("nodes/%s/%s", id, msg)))
}

func (s *State) GetNodeMsg(id, msg string) ([]byte, error) {
	return redis.Bytes(s.do("GET", fmt.Sprintf("nodes/%s/msg/%s", id, msg)))
}
