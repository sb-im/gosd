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
	if err := h.orm.WithContext(c).Offset((page-1)*size).Limit(size).Find(&schedules, "team_id = ?", h.getCurrent(c).TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// @Summary Schedule Create
// @Schemes Schedule
// @Description create a new schedules
// @Tags schedule
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name" default(Test Schedule)
// @Param cron formData string true "cron expression https://pkg.go.dev/github.com/robfig/cron/v3#hdr-Usage" default(@every 1h30m)
// @Param enable formData bool true "Enable" default(false)
// @Param method formData string true "Method" default(cowSay)
// @Param params formData string false "Params" default(Hello, world!)
// @Success 201 {object} model.Schedule
// @Router /schedules [post]
func (h Handler) ScheduleCreate(c *gin.Context) {
	schedule := model.Schedule{
		TeamID: h.getCurrent(c).TeamID,
	}
	if err := c.ShouldBind(&schedule); err != nil {
		log.Error(schedule, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.orm.WithContext(c).Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		h.srv.ScheduleAdd(schedule)
	}

	c.JSON(http.StatusCreated, schedule)
}

// @Summary Schedule Update
// @Schemes Schedule
// @Description update a new schedules
// @Tags schedule
// @Accept multipart/form-data
// @Produce json
// @Param id   path int            true "Schedule ID"
// @Param name formData string false "Name" default(Test Schedule)
// @Param cron formData string false "cron expression https://pkg.go.dev/github.com/robfig/cron/v3#hdr-Usage" default(@every 1h30m)
// @Param enable formData bool   false "Enable" default(false)
// @Param method formData string false "Method" default(cowSay)
// @Param params formData string false "Params" default(Hello, world!)
// @Success 200 {object} model.Schedule
// @Router /schedules/{id} [patch]
func (h Handler) ScheduleUpdate(c *gin.Context) {
	schedule := model.Schedule{}
	h.orm.WithContext(c).First(&schedule, c.Param("id"))
	if err := c.ShouldBind(&schedule); err != nil {
		log.Error(schedule, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// NOTE: this because golang default bool value: false
	// No way to judge whether to change
	// Field 'enable' will Force Updated
	// If not 'enable' key, status default to false
	h.orm.WithContext(c).Updates(&schedule).Select("enable").Updates(&schedule)
	h.srv.ScheduleUpdate(schedule)
	c.JSON(http.StatusOK, schedule)
}

// @Summary Schedule Destroy
// @Schemes Schedule
// @Description destroy a new schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 204 {object} model.Schedule
// @Router /schedules/{id} [delete]
func (h Handler) ScheduleDestroy(c *gin.Context) {
	schedule := model.Schedule{}
	h.orm.WithContext(c).First(&schedule, c.Param("id"))

	// Need Destroy cron
	// h.orm.Delete(&model.Schedule{}, c.Param("id"))
	h.orm.WithContext(c).Delete(&schedule)
	h.srv.ScheduleDel(schedule)
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Schedule Toggle
// @Schemes Schedule
// @Description toggle a schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} model.Schedule
// @Router /schedules/{id}/trigger [POST]
func (h *Handler) ScheduleTrigger(c *gin.Context) {
	schedule := &model.Schedule{}
	h.orm.WithContext(c).First(schedule, c.Param("id"))
	if err := h.srv.JSON.Call(schedule.Method, []byte(schedule.Params)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, schedule)
}
