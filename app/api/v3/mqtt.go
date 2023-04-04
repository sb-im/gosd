package v3

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// @Summary Create a mqtt user
// @Schemes Mqtt
// @Description create a new mqtt user
// @Tags mqtt
// @Accept multipart/form-data
// @Produce json
// @Router /mqtt/url [POST]
func (h *Handler) MqttUserCreate(c *gin.Context) {
	u, err := url.Parse(h.cfg.ApiMqttWs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if h.cfg.EmqxAuth {
		user := h.getCurrent(c)
		username, password, err := h.srv.MqttAuthReqTeam(user.TeamID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		u.User = url.UserPassword(username, password)
		h.srv.MqttAuthAclTeam(user.TeamID)
	}
	c.JSON(http.StatusOK, u.String())
}
