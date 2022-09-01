package service

import (
	//"encoding/json"
	"context"
	"errors"

	"sb.im/gosd/app/model"
)

func NewJsonService(s *Service) *JSONService {
	j := &JSONService{s: s}
	j.s.JSON = j
	j.m = map[string]func([]byte) error{
		"cowSay":  j.CowSay,
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
	//if err := json.Unmarshal(raw, params); err != nil {
	//	return err
	//}
	s.s.orm.First(params, string(raw))
	return s.s.TaskRun(context.Background(), params)
}

func (s *JSONService) CowSay(raw []byte) error {
	return s.s.CowSay(string(raw))
}
