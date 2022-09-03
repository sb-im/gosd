package v3

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

// @Summary Task Job Index
// @Schemes Job
// @Description get a tasks job index
// @Tags task
// @Accept json
// @Produce json
// @Param id path uint true "Task ID"
// @Param page query uint false "Task Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200
// @Router /tasks/{id}/jobs [get]
func (h *Handler) JobIndex(c *gin.Context) {
	// verify teams
	var count int64
	h.orm.WithContext(c).Find(&model.Task{}, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, nil)
	}

	var jobs []model.Job
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	if err := h.orm.WithContext(c).Order("id desc").Offset((page-1)*size).Limit(size).Find(&jobs, "task_id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

// @Summary Create a Job
// @Schemes Task
// @Description create a new task
// @Tags task
// @Produce json
// @Param id path uint true "Task ID"
// @Success 201
// @Router /tasks/{id}/jobs [post]
func (h *Handler) JobCreate(c *gin.Context) {
	h.TaskRunningCreate(c)
}
