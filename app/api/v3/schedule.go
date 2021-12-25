package v3

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
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
// @Accept json, multipart/form-data
// @Produce json
// @Param data body model.Schedule true "Schedule"
// @Success 200 {object} model.Schedule
// @Router /schedules [post]
func (h Handler) ScheduleCreate(c *gin.Context) {
	schedule := &model.Schedule{}
	if err := c.ShouldBind(schedule); err != nil {
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

func (h Handler) ScheduleToggle(c *gin.Context) {
	schedule := &model.Schedule{}
	if err := c.ShouldBind(schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(schedule)
	fmt.Println(c.Param("id"))
	if err := h.srv.JSON.Call(schedule.Target, []byte(schedule.Params)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//h.orm.Create(schedule)
	c.JSON(http.StatusOK, schedule)
}