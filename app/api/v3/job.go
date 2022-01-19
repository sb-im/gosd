package v3

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func (h Handler) JobIndex(c *gin.Context) {
	// verify teams
	var count int64
	h.orm.Find(&model.Task{}, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, nil)
	}

	var jobs []model.Job
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	h.orm.Offset((page-1)*size).Limit(size).Find(&jobs, "task_id = ?", c.Param("id"))
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
func (h Handler) JobCreate(c *gin.Context) {
	var task model.Task
	if err := h.orm.First(&task, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := h.srv.TaskRun(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, nil)
}
