package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sb.im/gosd/app/model"
)

func (h *Handler) DatabaseMigrate(c *gin.Context) {
	orm := h.orm

	orm.AutoMigrate(&model.Team{})
	orm.AutoMigrate(&model.User{})
	orm.AutoMigrate(&model.Session{})
	orm.AutoMigrate(&model.UserTeam{})

	orm.AutoMigrate(&model.Schedule{})
	orm.AutoMigrate(&model.Node{})
	orm.AutoMigrate(&model.Task{})
	orm.AutoMigrate(&model.Blob{})
	orm.AutoMigrate(&model.Job{})
	orm.AutoMigrate(&model.Profile{})

	c.JSON(http.StatusNoContent, nil)
}

// seed:
//   - default
//   - test
//   - demo
//   - multi
func (h *Handler) DatabaseSeed(c *gin.Context) {
	h.InitSeed()
	c.JSON(http.StatusNoContent, nil)
}
