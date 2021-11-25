package v3

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

type Current struct {
	TeamID uint
	UserID uint
}

func (c Current) isUser() bool {
	return c.UserID != 0
}

func (h Handler) singleUserMode() bool {
	return h.cfg.SingleUserMode
}

func (h Handler) getCurrent(c *gin.Context) *Current {
	current, ok := c.Get(identityGinKey)
	if ok {
		return current.(*Current)
	}

	if h.singleUserMode() {
		return &Current{
			TeamID: 1,
			UserID: 1,
		}
	}
	return nil
}

// @Summary Get Current User info
// @Schemes Auth
// @Description user login
// @Tags auth
// @Produce json
// @Success 200 {object} model.Schedule
// @Router /current [GET]
func (h Handler) current(c *gin.Context) {
	current := h.getCurrent(c)
	if current.isUser() {
		var user model.User
		h.orm.First(&user, current.UserID)

		fmt.Println(user.Teams)
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current User Error"})
	}
}
