package luavm

import (
	"encoding/json"
	"fmt"
)

func (s *Service) GetAttach() (string, error) {
	// Plan
	plan, err := s.Store.PlanByID(s.Task.PlanID)
	if err != nil {
		return "", err
	}
	s.Task.Files = plan.Files
	s.Task.Extra = plan.Extra

	// Job
	job, err := s.Store.PlanLogByID(s.Task.Job.JobID)
	if err != nil {
		return "", err
	}
	s.Task.Job.Files = job.Files
	s.Task.Job.Extra = job.Extra

	data, err := json.Marshal(s.Task)
	if err != nil {
		return string(data), err
	}
	return string(data), nil
}

func (s *Service) SetAttach(raw string) error {
	if err := json.Unmarshal([]byte(raw), s.Task); err != nil {
		// NOTE: is null, json.Unmarshal return error
		//return err
	}

	// Plan
	plan, err := s.Store.PlanByID(s.Task.PlanID)
	if err != nil {
		return err
	}
	plan.Files = s.Task.Files
	plan.Extra = s.Task.Extra
	if err := s.Store.UpdatePlan(plan); err != nil {
		return err
	}

	// Job
	job, err := s.Store.PlanLogByID(s.Task.Job.JobID)
	if err != nil {
		return err
	}
	job.Files = s.Task.Job.Files
	job.Extra = s.Task.Job.Extra
	if err := s.Store.UpdatePlanLog(job); err != nil {
		return err
	}

	return s.State.Record(fmt.Sprintf(topic_running, s.Task.PlanID), []byte(raw))
}
