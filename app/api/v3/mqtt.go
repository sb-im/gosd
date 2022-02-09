package v3

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create a mqtt user
// @Schemes Mqtt
// @Description create a new mqtt user
// @Tags mqtt
// @Accept multipart/form-data
// @Produce json
// @Router /mqtt/url [POST]
func (h Handler) MqttUserCreate(c *gin.Context) {

	// TODO: this need config file
	mqttURLFmt := "mqtt://%d:%s@localhost:1883"

	u := h.getCurrent(c)

	// TODO: isSuperUser
	fmt.Println(u)

	passwd := h.srv.MqttAuthUser(strconv.Itoa(int(u.SessID)))
	h.srv.MqttAuthACL(u.TeamID, strconv.Itoa(int(u.SessID)))
	c.JSON(http.StatusOK, fmt.Sprintf(mqttURLFmt, u.SessID, passwd))
}
