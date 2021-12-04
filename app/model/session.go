package model

type Session struct {
	Model

	TeamID uint   `json:"team_id" gorm:"primaryKey"`
	Team   Team   `json:"team"`
	UserID uint   `json:"user_id" gorm:"primaryKey"`
	User   User   `json:"user"`
	IP     string `json:"ip"`

	// TODO
	// sources: login, token
}
