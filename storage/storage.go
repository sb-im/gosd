package storage

import (
	"database/sql"

	"sb.im/gosd/model"
)

// Storage handles all operations related to the database.
type Storage struct {
	db    *sql.DB
	cache userToken
}

// NewStorage returns a new Storage.
func NewStorage(db *sql.DB) *Storage {
	return &Storage{db, userToken{token: make(map[string]*model.User)}}
}

func (s *Storage) Database() *sql.DB {
	return s.db
}
