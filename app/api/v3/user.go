package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"sb.im/gosd/app/model"
)

type bindUser struct {
	TeamID int `json:"team_id"`

	Username string `json:"username"`
	Password string `json:"password"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *Handler) UserCreate(c *gin.Context) {
	binduser := &bindUser{}
	if err := c.BindJSON(binduser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password, err := hashPassword(binduser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := model.User{
		TeamID:   binduser.TeamID,
		Username: binduser.Username,
		Password: password,
		Language: binduser.Language,
		Timezone: binduser.Timezone,
	}
	h.orm.Create(&user)
	c.JSON(http.StatusOK, user)
}
