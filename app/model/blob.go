package model

type Blob struct {
	Model

	UXID string `json:"uxid" gorm:"uniqueIndex;column:uxid"`
	Name string `json:"name"`
}
