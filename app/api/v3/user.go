package v3

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

// @Summary Create a user
// @Schemes User
// @Description create a new user
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param team_id  formData uint true "Team ID"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param language formData string false "Language"
// @Param timezone formData string false "Timezone"
// @Success 200 {object} model.User
// @Router /users [post]
func (h *Handler) UserCreate(c *gin.Context) {
	type bindUser struct {
		TeamID   uint   `json:"team_id" form:"team_id" binding:"required"`
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Language string `json:"language" form:"language"`
		Timezone string `json:"timezone" form:"timezone"`
	}

	u := &bindUser{}
	if err := c.Bind(u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		TeamID:   u.TeamID,
		Username: u.Username,
		Password: u.Password,
		Language: u.Language,
		Timezone: u.Timezone,
	}
	h.orm.Create(user)
	c.JSON(http.StatusOK, user)
}

// @Summary Update a user
// @Schemes User
// @Description update a new user
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "User ID"
// @Param team_id  formData uint false "Team ID"
// @Param username formData string false "Username"
// @Param password formData string false "Password"
// @Param language formData string false "Language"
// @Param timezone formData string false "Timezone"
// @Success 200 {object} model.User
// @Router /users/{id} [patch]
func (h *Handler) UserUpdate(c *gin.Context) {
	type bindUser struct {
		TeamID   uint   `json:"team_id" form:"team_id"`
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
		Language string `json:"language" form:"language"`
		Timezone string `json:"timezone" form:"timezone"`
	}

	u := &bindUser{}
	if err := c.Bind(u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		TeamID:   u.TeamID,
		Username: u.Username,
		Password: u.Password,
		Language: u.Language,
		Timezone: u.Timezone,
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(id)
	h.orm.Updates(user)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) userOverride() {
	for id, user := range h.cfg.Auth.SuperAdmin {
		user.ID = id
		h.orm.Updates(user)
	}
}
