package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"sb.im/gosd/app/logger"
	"sb.im/gosd/app/model"
)

const (
	schedulePrefix = "schedule.%d"
)

func (s *Service) RunSchedule(ctx context.Context) {
	schedules := []model.Schedule{}
	if err := s.orm.WithContext(ctx).Find(&schedules, "enable = TRUE"); err != nil {
		logger.WithContext(ctx).Error(err)
	}

	for _, i := range schedules {
		s.ScheduleAdd(i)
	}

	s.cron.Start()
	<-ctx.Done()
	s.cron.Stop()
}

func (s *Service) scheduleEntrySet(id uint, entryID int) error {
	return s.rdb.Set(context.Background(), fmt.Sprintf(schedulePrefix, id), entryID, 0).Err()
}

func (s *Service) scheduleEntryGet(id uint) (int, error) {
	return s.rdb.Get(context.Background(), fmt.Sprintf(schedulePrefix, id)).Int()
}

func (s *Service) scheduleEntryDel(id uint) error {
	return s.rdb.Del(context.Background(), fmt.Sprintf(schedulePrefix, id)).Err()
}

func (s *Service) ScheduleAdd(schedule model.Schedule) {
	if schedule.Enable {
		if entryID, err := s.cron.AddFunc(schedule.Cron, func() {

			ctx := context.WithValue(context.Background(), "traceid", uuid.New().String())
			logger.WithContext(ctx).WithField("src", "cron")

			s.ScheduleCreateJob(ctx, schedule.TaskID)
		}); err != nil {
			fmt.Println(err)
		} else {
			s.scheduleEntrySet(schedule.ID, int(entryID))
		}
	}
}

func (s *Service) ScheduleDel(schedule model.Schedule) {
	if entryID, err := s.scheduleEntryGet(schedule.ID); err != nil {
		fmt.Println(err)
	} else {
		s.cron.Remove(cron.EntryID(entryID))
		s.scheduleEntryDel(schedule.ID)
	}
}

func (s *Service) ScheduleUpdate(schedule model.Schedule) {
	s.ScheduleDel(schedule)
	s.ScheduleAdd(schedule)
}
