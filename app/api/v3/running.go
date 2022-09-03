package v3

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sb.im/gosd/app/model"
)

// @Summary Task Run
// @Schemes Running
// @Description Start run a task
// @Tags running
// @Accept json
// @Produce json
// @Param id path uint true "Task ID"
// @Success 201
// @Failure 500
// @Router /tasks/{id}/running [POST]
func (h *Handler) TaskRunningCreate(c *gin.Context) {
	var task model.Task
	if err := h.orm.WithContext(c).Model(&task).Where("id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).UpdateColumn("index", gorm.Expr("index + ?", 1)).Scan(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	job := model.Job{
		TaskID: task.ID,
		Index:  task.Index,
	}
	if err := h.orm.WithContext(c).Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	task.Job = &job

	if err := h.srv.TaskRun(c, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// @Summary Task Cancel
// @Schemes Running
// @Description Cancel a running task
// @Tags running
// @Accept json
// @Produce json
// @Param id path uint true "Task ID"
// @Success 204
// @Failure 404
// @Failure 500
// @Router /tasks/{id}/running [DELETE]
func (h *Handler) TaskRunningDestroy(c *gin.Context) {
	if !h.taskTeamIsExist(c.Param("id"), h.getCurrent(c).TeamID) {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	if err := h.srv.TaskKill(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
