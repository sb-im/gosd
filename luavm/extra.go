package luavm

import (
	"encoding/json"
	"fmt"
)

// GET Plan Extra
func (s *Service) GetExtra(key string) string {
	data, err := s.State.BytesGet(fmt.Sprintf(topic_running, s.Task.PlanID))
	if err != nil {}
	if err := json.Unmarshal(data, s.Task); err != nil {}

	if key, ok := s.Task.Extra[key]; ok {
		return key
	} else {
		return ""
	}
}

// SET Plan Extra
func (s *Service) SetExtra(key, value string) {
	if value == "" {
		delete(s.Task.Extra, key)
	} else {
		s.Task.Extra[key] = value
	}
	data, _ := json.Marshal(s.Task)
	s.State.Record(fmt.Sprintf(topic_running, s.Task.PlanID), data)
}

// GET Job Extra
func (s *Service) GetJobExtra(key string) string {
	data, err := s.State.BytesGet(fmt.Sprintf(topic_running, s.Task.PlanID))
	if err != nil {}
	if err := json.Unmarshal(data, s.Task); err != nil {}

	if key, ok := s.Task.Job.Extra[key]; ok {
		return key
	} else {
		return ""
	}
}

// SET Job Extra
func (s *Service) SetJobExtra(key, value string) {
	if value == "" {
		delete(s.Task.Job.Extra, key)
	} else {
		s.Task.Job.Extra[key] = value
	}
	data, _ := json.Marshal(s.Task)
	s.State.Record(fmt.Sprintf(topic_running, s.Task.PlanID), data)
}
