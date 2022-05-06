package v3

import (
	"errors"
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
// @Success 200
// @Failure 500
// @Router /tasks [GET]
func (h *Handler) TaskIndex(c *gin.Context) {
	var tasks []model.Task
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if err := h.orm.Offset((page-1)*size).Limit(size).Find(&tasks, "team_id = ?", h.getCurrent(c).TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// @Summary Create a task
// @Schemes Task
// @Description create a new task
// @Tags task
// @Accept multipart/form-data
// @Produce json
// @Param name    formData string true "Task Name"
// @Param node_id formData string true "Node ID"
// @Success 201
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /tasks [POST]
func (h *Handler) TaskCreate(c *gin.Context) {
	task := &model.Task{
		TeamID: h.getCurrent(c).TeamID,
	}
	if err := c.Bind(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify: `node_id`
	if err := h.orm.First(&model.Node{}, "id = ? AND team_id = ?", task.NodeID, h.getCurrent(c).TeamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := h.orm.Create(task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Need Sync mqtt acl new task, this team
	h.srv.MqttAuthAclTeam(h.getCurrent(c).TeamID)

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
// @Success 200
// @failure 404
// @failure 500
// @Router /tasks/{id} [GET]
func (h *Handler) TaskShow(c *gin.Context) {
	var task model.Task
	if err := h.orm.First(&task, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	page := mustStringToInt(c.Query("page"))
	size := mustStringToInt(c.Query("size"))

	if err := h.orm.Order("index desc").Offset((page-1)*size).Limit(size).Find(&task.Jobs, "task_id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary Task Update
// @Schemes Task
// @Description update a task
// @Tags task
// @Accept json
// @Produce json
// @Param id   path int true "Task ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /tasks/{id} [PUT]
func (h *Handler) TaskUpdate(c *gin.Context) {
	if id, _ := h.store.LockTaskGet(c.Param("id")); id != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "This Task is Running"})
		return
	}

	task := model.Task{}
	if err := c.ShouldBind(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify: `node_id`
	if err := h.orm.First(&model.Node{}, "id = ? AND team_id = ?", task.NodeID, h.getCurrent(c).TeamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := h.orm.Where("id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Updates(&task).Scan(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
// @Failure 500
// @Router /tasks/{id} [DELETE]
func (h *Handler) TaskDestroy(c *gin.Context) {
	if id, _ := h.store.LockTaskGet(c.Param("id")); id != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "This Task is Running"})
		return
	}

	if err := h.orm.Delete(&model.Task{}, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) taskTeamIsExist(taskId, teamId interface{}) bool {
	var count int64
	h.orm.Find(&model.Task{}, "id = ? AND team_id = ?", taskId, teamId).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
