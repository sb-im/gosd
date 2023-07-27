package service

import (
	"context"
)

func (s *Service) Reset(ctx context.Context) error {
	prefix := "luavm.*"
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = s.rdb.Scan(ctx, cursor, prefix, 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			s.rdb.Del(ctx, key)
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
