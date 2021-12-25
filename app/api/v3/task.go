package v3

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	h.orm.Offset((page - 1) * size).Limit(size).Find(&tasks, "team_id = ?", h.getCurrent(c).TeamID)
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
// @Success 200 {object} model.Task
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
	c.JSON(http.StatusOK, task)
}
