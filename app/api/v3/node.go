package v3

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sb.im/gosd/app/model"
)

// @Summary Node Index
// @Schemes Node
// @Description get all nodes index
// @Tags node
// @Accept json
// @Produce json
// @Param page query uint false "Task Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200
// @Failure 500
// @Router /nodes [GET]
func (h *Handler) NodeIndex(c *gin.Context) {
	var nodes []model.Node
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if err := h.orm.Offset((page-1)*size).Limit(size).Find(&nodes, "team_id = ?", h.getCurrent(c).TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// @Summary Create a node
// @Schemes Node
// @Description create a new node
// @Tags node
// @Accept multipart/form-data
// @Produce json
// @Param name   formData string true "Node Name"
// @Param points formData string false "Points"
// @Success 201
// @Failure 500
// @Router /nodes [POST]
func (h *Handler) NodeCreate(c *gin.Context) {
	node := &model.Node{
		TeamID: h.getCurrent(c).TeamID,
	}
	if err := c.Bind(node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.orm.Create(node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, node)
}

// @Summary Node Show
// @Schemes Node
// @Description show a node detail
// @Tags node
// @Accept json
// @Produce json
// @Param id path uint true "Task ID"
// @Success 200
// @failure 404
// @Router /nodes/{id} [GET]
func (h *Handler) NodeShow(c *gin.Context) {
	var node model.Node
	if err := h.orm.First(&node, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, node)
}

// @Summary Node Update
// @Schemes Node
// @Description update a node
// @Tags node
// @Accept json
// @Produce json
// @Param id   path int true "Node ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /nodes/{id} [PUT]
func (h *Handler) NodeUpdate(c *gin.Context) {
	node := model.Node{}
	if err := c.ShouldBind(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orm.Where("id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Updates(&node).Scan(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

// @Summary Node Destroy
// @Schemes Node
// @Description Destroy a node
// @Tags node
// @Accept json
// @Produce json
// @Param id path uint true "Node ID"
// @Success 204
// @Failure 500
// @Router /nodes/{id} [DELETE]
func (h *Handler) NodeDestroy(c *gin.Context) {
	if err := h.orm.Delete(&model.Task{}, "id = ? AND team_id = ?", c.Param("id"), h.getCurrent(c).TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
