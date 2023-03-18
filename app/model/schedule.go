package model

import (
	"errors"

	"gorm.io/gorm"
)

type Schedule struct {
	Model
	TeamID uint   `json:"-"`
	TaskID uint   `json:"task" gorm:"not null" form:"task_id"`
	Enable bool   `json:"enable" gorm:"not null;default:false" form:"enable"`
	Name   string `json:"name" form:"name"`
	Cron   string `json:"cron" form:"cron"`
}

func (s *Schedule) BeforeSave(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Task{}).Where("team_id = ? AND id = ?", s.TeamID, s.TaskID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Not Found This taks")
	}
	return nil
}
