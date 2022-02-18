package v3

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

// @Summary Index all users
// @Schemes User
// @Description get all users index
// @Tags user
// @Accept json
// @Produce json
// @Param page query uint false "Task Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200
// @Router /users [GET]
func (h *Handler) UserIndex(c *gin.Context) {
	var users []model.User
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	h.orm.Offset((page - 1) * size).Limit(size).Find(&users)
	c.JSON(http.StatusOK, users)
}

// @Summary Create a user
// @Schemes User
// @Description create a new user
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param team_id  formData uint true "Team ID"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param language formData string false "Language Codes, such as 'zh_CN'" default(en_US)
// @Param timezone formData string false "IANA Time Zone database, such as 'America/New_York'" default(Asia/Shanghai)
// @Success 201 {object} model.User
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
	c.JSON(http.StatusCreated, user)
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

func (h *Handler) UserOverride() {
	// TODO:
	//for id, user := range h.cfg.SuperAdmin {
	//	user.ID = id
	//	h.orm.Updates(user)
	//}
}

func (h *Handler) userIsExist(id uint) bool {
	var count int64
	h.orm.Find(&model.User{}, id).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
