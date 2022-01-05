package v3

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sb.im/gosd/app/model"
)

// @Summary Task Index
// @Schemes Task
// @Description get all tasks index
// @Tags task
// @Accept json
// @Produce json
// @Param page query uint false "Task Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200 {object} model.Task
// @Router /tasks [get]
func (h *Handler) TaskIndex(c *gin.Context) {
	var tasks []model.Task
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	h.orm.Offset((page-1)*size).Limit(size).Find(&tasks, "team_id = ?", h.getCurrent(c).TeamID)
	c.JSON(http.StatusOK, tasks)
}

// @Summary Create a task
// @Schemes Task
// @Description create a new task
// @Tags task
// @Accept multipart/form-data
// @Produce json
// @Param name    formData string true "Task Name"
// @Param node_id formData uint true "Node ID"
// @Success 201 {object} model.Task
// @Router /tasks [post]
func (h *Handler) TaskCreate(c *gin.Context) {
	task := &model.Task{
		TeamID: h.getCurrent(c).TeamID,
	}
	if err := c.Bind(task); err != nil {
		fmt.Println(string(task.Extra))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(task)
	c.JSON(http.StatusCreated, task)
}

// @Summary Task Show
// @Schemes Task
// @Description show a tasks detail
// @Tags task
// @Accept json
// @Produce json
// @Param id path uint true "Task ID"
// @Param page query uint false "Task Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200 {object} model.Task
// @failure 404
// @Router /tasks/{id} [get]
func (h Handler) TaskShow(c *gin.Context) {
	var task model.Task
	if err := h.orm.First(&task, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// TODO:
	//page, _ := strconv.Atoi(c.Query("page"))
	//size, _ := strconv.Atoi(c.Query("size"))
	//h.orm.Offset((page - 1) * size).Limit(size).Find(&task, "team_id = ?", h.getCurrent(c).TeamID)
	c.JSON(http.StatusOK, task)
}

// @Summary Task Update
// @Schemes Task
// @Description update a task
// @Tags task
// @Accept json
// @Produce json
// @Param id   path int true "Task ID"
// @Param data body model.Task true "Task"
// @Success 200 {object} model.Task
// @Router /tasks/{id} [put]
func (h Handler) TaskUpdate(c *gin.Context) {
	task := model.Task{}
	if err := c.ShouldBind(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.orm.Where("id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Updates(&task).Scan(&task)
	c.JSON(http.StatusOK, task)
}

// @Summary Task Destroy
// @Schemes Task
// @Description Destroy a task
// @Tags task
// @Accept json
// @Produce json
// @Param id path uint true "Task ID"
// @Success 204
// @Router /tasks/{id} [delete]
func (h Handler) TaskDestroy(c *gin.Context) {
	h.orm.Delete(&model.Task{}, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID)
	c.JSON(http.StatusNoContent, nil)
}
