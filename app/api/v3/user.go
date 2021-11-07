package v3

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

type bindUser struct {
	TeamID   int    `json:"team_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

func (h *Handler) UserCreate(c *gin.Context) {
	binduser := &bindUser{}
	if err := c.BindJSON(binduser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := h.userConvert(binduser)
	h.orm.Create(user)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UserUpdate(c *gin.Context) {
	binduser := &bindUser{}
	if err := c.BindJSON(binduser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := h.userConvert(binduser)
	h.orm.Updates(user)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(id)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) userConvert(u *bindUser) *model.User {
	return &model.User{
		TeamID:   u.TeamID,
		Username: u.Username,
		Password: u.Password,
		Language: u.Language,
		Timezone: u.Timezone,
	}
}
