package model

type Schedule struct {
	Model
	TeamID uint   `json:"-"`
	TaskID uint   `json:"task" gorm:"not null" form:"task_id"`
	Enable bool   `json:"enable" gorm:"not null;default:false" form:"enable"`
	Name   string `json:"name" form:"name"`
	Cron   string `json:"cron" form:"cron"`
}
