package v3

import (
	"net/http"
	"time"

	"sb.im/gosd/version"

	"github.com/gin-gonic/gin"
)

// @Summary Status
// @Schemes Status
// @Description Get Server Status, Current Time, language, timezone, version, build_date
// @Tags status
// @Accept json
// @Produce json
// @Success 200
// @Router /status [GET]
func (h *Handler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"time":        time.Now().String(),
		"language":    h.cfg.Language,
		"timezone":    h.cfg.Timezone,
		"version":     version.Version,
		"build_date:": version.Date,
	})
}
