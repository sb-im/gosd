package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

// @Summary Create a team
// @Schemes Team
// @Description create a new team
// @Tags team
// @Accept multipart/form-data
// @Produce json
// @Param   name formData string true "Team Name"
// @Success 200 {object} model.Team
// @Router /teams [post]
func (h *Handler) TeamCreate(c *gin.Context) {
	team := &model.Team{}
	if err := c.Bind(team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(team)
	c.JSON(http.StatusOK, team)
}
