package v3

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

func (h *Handler) TaskIndex(c *gin.Context) {
	var tasks []model.Task
	h.orm.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) TaskCreate(c *gin.Context) {
	task := &model.Task{}
	if err := c.BindJSON(task); err != nil {
		fmt.Println(string(task.Extra))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(task)
	c.JSON(http.StatusOK, task)
}
