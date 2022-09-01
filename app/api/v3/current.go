package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/logger"
	"sb.im/gosd/app/model"
)

type Current struct {
	TeamID uint
	UserID uint
	SessID uint
}

func (c Current) isUser() bool {
	return c.UserID != 0
}

func (h Handler) singleUserMode() bool {
	return h.cfg.SingleUser
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
			SessID: 1,
		}
	}
	return nil
}

// @Summary Get Current User info
// @Schemes Auth
// @Description user login
// @Tags auth
// @Security BasicAuth
// @Security JWTSecret
// @Produce json
// @Success 200 {object} model.Schedule
// @Router /current [GET]
func (h *Handler) Current(c *gin.Context) {
	current := h.getCurrent(c)
	logger.WithContext(c).
		WithField("teamId", current.TeamID).
		WithField("userId", current.UserID).
		WithField("sessId", current.SessID).
		Infof("%+v", *current)

	if current.isUser() {
		var user model.User
		if err := h.orm.WithContext(c).Preload("Teams").First(&user, current.UserID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current User Error"})
	}
}
