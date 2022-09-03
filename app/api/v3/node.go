package v3

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sb.im/gosd/app/helper"
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
	if err := h.orm.WithContext(c).Offset((page-1)*size).Limit(size).Find(&nodes, "team_id = ?", h.getCurrent(c).TeamID).Error; err != nil {
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
		Secret: helper.GenSecret(16),
	}
	err := c.Bind(node)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If not set uuid, need to create a uuid = id
	if node.UUID == "" {
		err = h.orm.WithContext(c).Create(node).Update("uuid", strconv.Itoa(int(node.ID))).Error
	} else {
		err = h.orm.WithContext(c).Create(node).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.srv.MqttAuthReqNode(node.UUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.srv.MqttAuthAclNode(node.UUID); err != nil {
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
// @Param id path string true "Node ID"
// @Success 200
// @failure 404
// @Router /nodes/{uuid} [GET]
func (h *Handler) NodeShow(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("uuid"))
	if id != 0 {
		h.nodeShowByID(c)
	} else {
		h.nodeShowByUUID(c)
	}
}

func (h *Handler) nodeShowByID(c *gin.Context) {
	var node model.Node
	if err := h.orm.WithContext(c).First(&node, "id = ? AND team_id = ?", c.Param("uuid"), h.getCurrent(c).TeamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, node)
}

func (h *Handler) nodeShowByUUID(c *gin.Context) {
	var node model.Node
	if err := h.orm.WithContext(c).First(&node, "uuid = ? AND team_id = ?", c.Param("uuid"), h.getCurrent(c).TeamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, node)
}

// @Summary Node Update
// @Schemes Node
// @Description update a node
// @Tags node
// @Accept json
// @Produce json
// @Param id path string true "Node ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /nodes/{uuid} [PUT]
func (h *Handler) NodeUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("uuid"))

	// TODO: Running uuid disable change
	if id != 0 {
		h.nodeUpdateByID(c)
	} else {
		h.nodeUpdateByUUID(c)
	}
}

func (h *Handler) nodeUpdateByID(c *gin.Context) {
	node := model.Node{}
	if err := c.ShouldBind(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orm.WithContext(c).Where("id = ? AND team_id = ?", c.Param("uuid"), h.getCurrent(c).TeamID).Updates(&node).Scan(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

func (h *Handler) nodeUpdateByUUID(c *gin.Context) {
	node := model.Node{}
	if err := c.ShouldBind(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orm.WithContext(c).Where("uuid = ? AND team_id = ?", c.Param("uuid"), h.getCurrent(c).TeamID).Updates(&node).Scan(&node).Error; err != nil {
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
// @Param id path string true "Node ID"
// @Success 204
// @Failure 500
// @Router /nodes/{uuid} [DELETE]
func (h *Handler) NodeDestroy(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("uuid"))

	// TODO: Running uuid disable change
	if id != 0 {
		h.nodeDestroyByID(c)
	} else {
		h.nodeDestroyByUUID(c)
	}
}

func (h *Handler) nodeDestroyByID(c *gin.Context) {
	if err := h.orm.WithContext(c).
		Model(&model.Node{}).
		Where("id = ? AND team_id = ?", c.Param("uuid"), h.getCurrent(c).TeamID).
		Update("uuid", nil).
		Delete(&model.Node{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) nodeDestroyByUUID(c *gin.Context) {
	if err := h.orm.WithContext(c).
		Model(&model.Node{}).
		Where("uuid = ? AND team_id = ?", c.Param("uuid"), h.getCurrent(c).TeamID).
		Update("uuid", nil).
		Delete(&model.Node{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
