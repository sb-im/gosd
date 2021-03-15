package luavm

import (
	"encoding/json"
	"fmt"
)

// GET Plan Extra
func (s *Service) GetExtra(key string) (string, error) {
	data, err := s.State.BytesGet(fmt.Sprintf(topic_running, s.Task.PlanID))
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, s.Task); err != nil {
		return "", err
	}

	if key, ok := s.Task.Extra[key]; ok {
		return key, nil
	} else {
		return "", nil
	}
}

// SET Plan Extra
func (s *Service) SetExtra(key, value string) error {
	if value == "" {
		delete(s.Task.Extra, key)
	} else {
		if s.Task.Extra == nil {
			s.Task.Extra = make(map[string]string)
		}
		s.Task.Extra[key] = value
	}
	data, _ := json.Marshal(s.Task)
	fmt.Println(string(data))
	return s.State.Record(fmt.Sprintf(topic_running, s.Task.PlanID), data)
}

// GET Job Extra
func (s *Service) GetJobExtra(key string) (string, error) {
	data, err := s.State.BytesGet(fmt.Sprintf(topic_running, s.Task.PlanID))
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, s.Task); err != nil {
		return "", err
	}

	if key, ok := s.Task.Job.Extra[key]; ok {
		return key, nil
	} else {
		return "", nil
	}
}

// SET Job Extra
func (s *Service) SetJobExtra(key, value string) error {
	if value == "" {
		delete(s.Task.Job.Extra, key)
	} else {
		if s.Task.Extra == nil {
			s.Task.Extra = make(map[string]string)
		}
		s.Task.Job.Extra[key] = value
	}
	data, _ := json.Marshal(s.Task)
	return s.State.Record(fmt.Sprintf(topic_running, s.Task.PlanID), data)
}
