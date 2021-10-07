package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/model"
)

func (h *Handler) ScheduleIndex(c *gin.Context) {
	var schedules []model.Schedule
	h.orm.Find(&schedules)
	c.JSON(http.StatusOK, schedules)
}

func (h *Handler) ScheduleCreate(c *gin.Context) {
	schedule := &model.Schedule{}
	if err := c.ShouldBind(schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(schedule)
	c.JSON(http.StatusOK, schedule)
}
