package luavm

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s Service) GetAttach() (string, error) {
	// TODO: Job
	//job, err := s.Store.PlanLogByID(s.Task.Job.JobID)
	//if err != nil {
	//	return "", err
	//}
	//s.Task.Job.Files = job.Files
	//s.Task.Job.Extra = job.Extra

	data, err := json.Marshal(s.Task)
	if err != nil {
		return string(data), err
	}
	return string(data), nil
}

func (s Service) SetAttach(raw string) error {
	if err := json.Unmarshal([]byte(raw), s.Task); err != nil {
		// NOTE: is null, json.Unmarshal return error
		//return err
	}

	err := s.orm.Updates(&s.Task).Select("files", "extra").Updates(&s.Task).Error
	if err != nil {
		return err
	}

	//// Job
	//job, err := s.Store.PlanLogByID(s.Task.Job.JobID)
	//if err != nil {
	//	return err
	//}
	//job.Files = s.Task.Job.Files
	//job.Extra = s.Task.Job.Extra
	//if err := s.Store.UpdatePlanLog(job); err != nil {
	//	return err
	//}

	return s.rdb.Set(context.Background(), fmt.Sprintf(topic_running, s.Task.ID), raw, 0).Err()
}
