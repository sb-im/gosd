package v3

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Status
// @Schemes Status
// @Description Get Server Status
// @Tags status
// @Accept json
// @Produce json
// @Success 200
// @Router /status [get]
func (h Handler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"time":     time.Now().String(),
		"language": h.cfg.DefaultLanguage,
		"timezone": h.cfg.DefaultTimezone,
		"version":  "Unknown",
	})
}
