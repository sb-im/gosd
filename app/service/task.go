package service

import (
	"context"
	"fmt"
	"time"

	"sb.im/gosd/app/model"
)

func (s *Service) StartTaskWorker(ctx context.Context) {
	jobs := []model.Job{}
	s.orm.Find(&jobs, "started_at > CURRENT_TIMESTAMP AND duration = 0")

	for _, job := range jobs {
		s.setWillJob(ctx, &job)
	}
}

func (s *Service) setWillJob(ctx context.Context, job *model.Job) error {
	duration := time.Until(job.StartedAt)
	if duration.Nanoseconds() < time.Second.Nanoseconds() {
		duration = time.Second
	}
	return s.rdb.Set(ctx, fmt.Sprintf("job.%d", job.ID), true, duration).Err()
}

func (s *Service) TaskRun(ctx context.Context, task *model.Task) error {
	return s.setWillJob(ctx, task.Job)
}

func (s *Service) TaskKill(taskId string) error {
	return s.worker.Kill(taskId)
}
