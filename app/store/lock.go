package store

import (
	"context"
	"time"
)

const (
	lockTaskPrefix = "luavm.lock.task."
	lockNodePrefix = "luavm.lock.node."

	lockTimeout = 2 * time.Hour
)

func (s *Store) LockTaskSet(id string) error {
	return s.rdb.Set(context.Background(), lockTaskPrefix+id, s.cfg.Instance, lockTimeout).Err()
}

func (s *Store) LockTaskGet(id string) (string, error) {
	return s.rdb.Get(context.Background(), lockTaskPrefix+id).Result()
}

func (s *Store) LockTaskDel(id string) error {
	return s.rdb.Del(context.Background(), lockTaskPrefix+id).Err()
}

func (s *Store) LockNodeSet(id string) error {
	return s.rdb.Set(context.Background(), lockNodePrefix+id, s.cfg.Instance, lockTimeout).Err()
}

func (s *Store) LockNodeGet(id string) (string, error) {
	return s.rdb.Get(context.Background(), lockNodePrefix+id).Result()
}

func (s *Store) LockNodeDel(id string) error {
	return s.rdb.Del(context.Background(), lockNodePrefix+id).Err()
}
