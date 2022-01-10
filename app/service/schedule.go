package service

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"sb.im/gosd/app/model"
)

func (s Service) StartSchedule() {
	schedules := []model.Schedule{}
	s.orm.Find(&schedules, "enable = TRUE")

	for _, i := range schedules {
		s.ScheduleAdd(i)
	}

	s.cron.Start()
}

func (s Service) scheduleEntrySet(id uint, entryID int) {
	s.rdb.Set(context.Background(), fmt.Sprintf("schedule.%d", id), entryID, 0)
}

func (s Service) scheduleEntryGet(id uint) (int, error) {
	return s.rdb.Get(context.Background(), fmt.Sprintf("schedule.%d", id)).Int()
}

func (s Service) scheduleEntryDel(id uint) {
	s.rdb.Del(context.Background(), fmt.Sprintf("schedule.%d", id))
}

func (s Service) ScheduleAdd(schedule model.Schedule) {
	if schedule.Enable {
		if entryID, err := s.cron.AddFunc(schedule.Cron, func() {
			s.JSON.Call(schedule.Method, []byte(schedule.Params))
		}); err != nil {
			fmt.Println(err)
		} else {
			s.scheduleEntrySet(schedule.ID, int(entryID))
		}
	}
}

func (s Service) ScheduleDel(schedule model.Schedule) {
	if entryID, err := s.scheduleEntryGet(schedule.ID); err != nil {
		fmt.Println(err)
	} else {
		s.cron.Remove(cron.EntryID(entryID))
		s.scheduleEntryDel(schedule.ID)
	}
}

func (s Service) ScheduleUpdate(schedule model.Schedule) {
	s.ScheduleDel(schedule)
	s.ScheduleAdd(schedule)
}
