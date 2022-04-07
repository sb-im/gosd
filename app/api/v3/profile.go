package v3

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sb.im/gosd/app/model"
)

// @Summary Get User Profile
// @Schemes Profile
// @Description USer Profile
// @Tags profile
// @Security APIKeyHeader
// @Accept multipart/form-data
// @Produce json
// @Param key path string true "Profile Key"
// @Success 200
// @Failure 500
// @Router /profiles/{key} [GET]
func (h *Handler) ProfileGet(c *gin.Context) {
	var profile model.Profile
	if err := h.orm.First(&profile, "key = ? AND user_id = ?", c.Param("key"), h.getCurrent(c).UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			content, err := defaultProfile(c.Param("key"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, json.RawMessage(content))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, profile.Data)
}

// @Summary Set User Profile
// @Schemes Profile
// @Description USer Profile
// @Tags profile
// @Security APIKeyHeader
// @Accept multipart/form-data
// @Produce json
// @Param key path string true "Profile Key"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /profiles/{key} [PUT]
func (h *Handler) ProfileSet(c *gin.Context) {
	var profile model.Profile
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	profile = model.Profile{
		UserID: h.getCurrent(c).UserID,
		Key:    c.Param("key"),
		Data:   data,
	}

	if err := h.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"data"}),
	}).Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusNoContent, nil)
}

func defaultProfile(profile string) ([]byte, error) {
	// TODO:
	return ioutil.ReadFile("data/" + profile + ".json")
}
