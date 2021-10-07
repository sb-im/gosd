package v3

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

func (h *Handler) scheduleIndex(c *gin.Context) {
	var schedules []model.Schedule
	h.orm.Find(&schedules)
	c.JSON(http.StatusOK, schedules)
}

func (h *Handler) scheduleCreate(c *gin.Context) {
	schedule := &model.Schedule{}
	if err := c.ShouldBind(schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(schedule)
	c.JSON(http.StatusOK, schedule)
}

func (h *Handler) scheduleUpdate(c *gin.Context) {
	schedule := &model.Schedule{}
	h.orm.Find(schedule, c.Param("id"))
	if err := c.ShouldBind(schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// NOTE: this because golang default bool value: false
	// No way to judge whether to change
	// Field 'enable' will Force Updated
	// If not 'enable' key, status default to false
	h.orm.Updates(schedule).Select("enable").Updates(schedule)
	c.JSON(http.StatusOK, schedule)
}

func (h *Handler) scheduleDestroy(c *gin.Context) {
	h.orm.Delete(&model.Schedule{}, c.Param("id"))
}

func (h *Handler) scheduleToggle(c *gin.Context) {
	schedule := &model.Schedule{}
	if err := c.ShouldBind(schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(schedule)
	fmt.Println(c.Param("id"))
	h.srv.PlanTask("1")
	//h.orm.Create(schedule)
	c.JSON(http.StatusOK, schedule)
}
