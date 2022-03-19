package model

type Job struct {
	Model

	Index  uint `json:"index" gorm:"index:idx_index,unique;not null"`
	TaskID uint `json:"-" gorm:"index:idx_index,unique;not null"`
	Task   Task `json:"-"`
	Files  JSON `json:"files"`
	Extra  JSON `json:"extra"`
}
