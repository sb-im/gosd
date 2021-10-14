package service

import (
	"encoding/json"
	"errors"

	"sb.im/gosd/app/model"
)

func NewJsonService(s *Service) *JSONService {
	j := &JSONService{s: s}
	j.s.JSON = j
	j.m = map[string]func([]byte) error{
		"taskRun": j.TaskRun,
	}
	return j
}

type JSONService struct {
	s *Service
	m map[string]func([]byte) error
}

func (s *JSONService) Call(method string, params []byte) error {
	if fn, ok := s.m[method]; ok {
		return fn(params)
	}
	return errors.New("Not Method Found")
}

func (s *JSONService) TaskRun(raw []byte) error {
	params := &model.Task{}
	if err := json.Unmarshal(raw, params); err != nil {
		return err
	}
	return s.s.TaskRun(params)
}
