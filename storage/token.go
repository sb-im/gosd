package storage

import "sb.im/gosd/model"

type userToken struct {
	token map[string]*model.User
}

func (s *Storage) CreateToken(token string, user *model.User) (err error) {
	s.cache.token[token] = user

	return nil
}

func (s *Storage) GetCurrentUser(token string) *model.User {
	return s.cache.token[token]
}
