package v3

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"

	log "github.com/sirupsen/logrus"
)

// @Summary Schedule Index
// @Schemes Schedule
// @Description get all schedules index
// @Tags schedule
// @Accept json
// @Produce json
// @Param page query uint false "Schedule Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200 {object} model.Schedule
// @Router /schedules [get]
func (h Handler) ScheduleIndex(c *gin.Context) {
	var schedules []model.Schedule
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	h.orm.Offset((page - 1) * size).Limit(size).Find(&schedules)
	c.JSON(http.StatusOK, schedules)
}

// @Summary Schedule Create
// @Schemes Schedule
// @Description create a new schedules
// @Tags schedule
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name" default(Test Schedule)
// @Param cron formData string true "cron expression https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format" default(@hourly)
// @Param enable formData bool true "Enable" default(false)
// @Param method formData string true "Method" default(cowSay)
// @Param params formData string false "Params" default(Hello, world!)
// @Success 200 {object} model.Schedule
// @Router /schedules [post]
func (h Handler) ScheduleCreate(c *gin.Context) {
	schedule := &model.Schedule{}
	if err := c.ShouldBind(schedule); err != nil {
		log.Error(schedule, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(schedule)
	c.JSON(http.StatusOK, schedule)
}

// @Summary Schedule Update
// @Schemes Schedule
// @Description update a new schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param id   path int            true "Schedule ID"
// @Param data body model.Schedule true "Schedule"
// @Success 200 {object} model.Schedule
// @Router /schedules/{id} [patch]
func (h Handler) ScheduleUpdate(c *gin.Context) {
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

// @Summary Schedule Destroy
// @Schemes Schedule
// @Description destroy a new schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param   id     path    int     true        "Schedule ID"
// @Success 200 {object} model.Schedule
// @Router /schedules/{id} [delete]
func (h Handler) ScheduleDestroy(c *gin.Context) {
	h.orm.Delete(&model.Schedule{}, c.Param("id"))
}

// @Summary Schedule Toggle
// @Schemes Schedule
// @Description toggle a schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param   id     path    int     true        "Schedule ID"
// @Success 200 {object} model.Schedule
// @Router /schedules/{id}/toggle [POST]
func (h Handler) ScheduleToggle(c *gin.Context) {
	schedule := &model.Schedule{}
	h.orm.First(schedule, c.Param("id"))
	if err := h.srv.JSON.Call(schedule.Method, []byte(schedule.Params)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, schedule)
}
