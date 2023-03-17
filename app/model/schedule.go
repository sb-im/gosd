package model

type Schedule struct {
	Model
	TeamID uint   `json:"-"`
	TaskID uint   `json:"task" gorm:"not null" form:"task"`
	Name   string `json:"name" form:"name"`
	Cron   string `json:"cron" form:"cron"`
	Enable bool   `json:"enable" gorm:"not null;default:false" form:"enable"`
	Method string `json:"method" form:"method"`
	Params string `json:"params" form:"params"`
}
