package v3

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sb.im/gosd/app/model"
)

// @Summary Index all teams
// @Schemes Team
// @Description get all teams index
// @Tags team
// @Accept json
// @Produce json
// @Param page query uint false "Task Page Num"
// @Param size query uint false "Page Max Count"
// @Success 200
// @Router /teams [GET]
func (h *Handler) TeamIndex(c *gin.Context) {
	var teams []model.Team
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	h.orm.Offset((page - 1) * size).Limit(size).Find(&teams)
	c.JSON(http.StatusOK, teams)
}

// @Summary Create a team
// @Schemes Team
// @Description create a new team
// @Tags team
// @Accept multipart/form-data
// @Produce json
// @Param   name formData string true "Team Name"
// @Success 201 {object} model.Team
// @Router /teams [post]
func (h *Handler) TeamCreate(c *gin.Context) {
	team := &model.Team{}
	if err := c.Bind(team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.orm.Create(team)
	c.JSON(http.StatusCreated, team)
}

// @Summary Team Add user
// @Schemes Team
// @Description add a existing user to the team, id > username
// @Tags team
// @Accept multipart/form-data
// @Produce json
// @Param id formData string false "User ID"
// @Param username formData string false "User Name"
// @Success 200 {object} model.Team
// @Router /teams/users [POST]
func (h Handler) TeamUserAdd(c *gin.Context) {
	current := h.getCurrent(c)

	type bindUser struct {
		ID       uint   `json:"id" form:"id"`
		Username string `json:"username" form:"username"`
	}

	form := &bindUser{}
	if err := c.Bind(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	h.orm.Where("id = ? OR username = ?", form.ID, form.Username).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found User"})
	}

	// Verify Except default user.TeamID
	if user.TeamID != current.TeamID {

		// TODO: need exception to already joined
		userTeam := &model.UserTeam{
			UserID: user.ID,
			TeamID: current.TeamID,
		}
		h.orm.Create(userTeam)
	}

	c.JSON(http.StatusOK, user)
}

func (h Handler) teamIsExist(id uint) bool {
	var count int64
	h.orm.Find(&model.Team{}, id).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
