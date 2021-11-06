package model

type User struct {
	Model

	TeamID int  `json:"team_id"`
	Team   Team `json:"team"`

	Username string `json:"username"`
	Password string `json:"-"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}
