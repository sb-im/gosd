package store

import (
	//"context"

	"sb.im/gosd/app/config"
	"sb.im/gosd/app/storage"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Store struct {
	orm *gorm.DB
	rdb *redis.Client
	cfg *config.Config
	ofs *storage.Storage
}

func NewStore(cfg *config.Config, orm *gorm.DB, rdb *redis.Client, ofs *storage.Storage) *Store {
	// Enable Redis Events
	//rdb.ConfigSet(context.Background(), "notify-keyspace-events", "$K")

	return &Store{
		cfg: cfg,
		orm: orm,
		rdb: rdb,
		ofs: ofs,
	}
}

func (s *Store) Cfg() *config.Config {
	return s.cfg
}

func (s *Store) Orm() *gorm.DB {
	return s.orm
}

func (s *Store) Rdb() *redis.Client {
	return s.rdb
}

func (s *Store) Ofs() *storage.Storage {
	return s.ofs
}
