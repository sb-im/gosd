package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

func (h *Handler) TeamCreate(c *gin.Context) {
	team := &model.Team{}
	if err := c.BindJSON(team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(team)
	c.JSON(http.StatusOK, team)
}
