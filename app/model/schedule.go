package model

type Schedule struct {
	Model

	Name   string `json:"name" example:"Test Schedule"`
	Cron   string `json:"cron"`
	Enable bool   `json:"enable" gorm:"not null;default:false"`
	Target string `json:"target"`
	Params string `json:"params"`
}
