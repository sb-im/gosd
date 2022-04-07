package model

type Profile struct {
	Model

	Key    string `json:"key" gorm:"index:profile_index,unique;not null"`
	UserID uint   `json:"-" gorm:"index:profile_index,unique;not null"`
	User   User   `json:"-"`
	Data   JSON   `json:"data"`
}
