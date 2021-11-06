package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

func (h *Handler) UserCreate(c *gin.Context) {
	user := &model.User{}
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(user)
	c.JSON(http.StatusOK, user)
}
