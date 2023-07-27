package service

import (
	"context"
	"testing"

	"sb.im/gosd/app/config"
	"sb.im/gosd/app/storage"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestReset(t *testing.T) {
	id := "test_id"
	srv := testNewService(t)

	if err := srv.LockTaskSet(id); err != nil {
		t.Error(err)
	}
	if _, err := srv.LockTaskGet(id); err != nil {
		t.Error(err)
	}

	if err := srv.LockNodeSet(id); err != nil {
		t.Error(err)
	}
	if _, err := srv.LockNodeGet(id); err != nil {
		t.Error(err)
	}

	srv.LockTaskGet(id)
	if err := srv.Reset(context.Background()); err != nil {
		t.Error(err)
	}

	if key, err := srv.LockTaskGet(id); err == nil {
		t.Error(key)
	}
	if key, err := srv.LockNodeGet(id); err == nil {
		t.Error(key)
	}
}

func testNewService(t *testing.T) *Service {
	cfg := config.Parse()
	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), nil)
	if err != nil {
		t.Error(err)
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		t.Error(err)
	}
	rdb := redis.NewClient(redisOpt)

	// Enable Redis Events
	// K: store
	// Ex: luavm
	rdb.ConfigSet(context.Background(), "notify-keyspace-events", "$KEx")

	ofs := storage.NewStorage(cfg.StorageURL)

	return NewService(cfg, orm, rdb, ofs)
}
