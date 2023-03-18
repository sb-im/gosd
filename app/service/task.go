package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"sb.im/gosd/app/model"
)

func (s *Service) StartTaskWorker(ctx context.Context) {
	jobs := []model.Job{}
	s.orm.Find(&jobs, "started_at > CURRENT_TIMESTAMP AND duration = 0")

	for _, job := range jobs {
		s.setWillJob(ctx, &job)
	}
}

// NOTE: Not Verified. Only Schedule call
func (s *Service) ScheduleCreateJob(ctx context.Context, taskId uint) error {
	var task model.Task
	if err := s.orm.WithContext(ctx).Model(&task).Where("id = ?", taskId).UpdateColumn("index", gorm.Expr("index + ?", 1)).Scan(&task).Error; err != nil {
		return err
	}
	return s.CreateJob(ctx, &model.Job{
		TaskID: task.ID,
		Index: task.Index,
	})
}

func (s *Service) CreateJob(ctx context.Context, job *model.Job) error {
	if job.StartedAt.Before(time.Now()) {
		job.StartedAt = time.Now()
	}

	if err := s.orm.WithContext(ctx).Create(&job).Error; err != nil {
		return err
	}

	return s.setWillJob(ctx, job)
}

func (s *Service) setWillJob(ctx context.Context, job *model.Job) error {
	duration := time.Until(job.StartedAt)
	if duration.Nanoseconds() < time.Second.Nanoseconds() {
		duration = time.Second
	}
	return s.rdb.Set(ctx, fmt.Sprintf("job.%d", job.ID), true, duration).Err()
}

func (s *Service) TaskKill(taskId string) error {
	return s.worker.Kill(taskId)
}
