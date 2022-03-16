package v3

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	mqttClientIdUserPrefix = "user."
	mqttClientIdNodePrefix = "node."
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
	user := h.getCurrent(c)
	username := mqttClientIdUserPrefix + strconv.Itoa(int(user.SessID))
	password := h.srv.MqttAuthUser(username)

	u.User = url.UserPassword(username, password)

	// TODO: isSuperUser
	fmt.Println(user)

	h.srv.MqttAuthACL(user.TeamID, username)
	c.JSON(http.StatusOK, u.String())
}
