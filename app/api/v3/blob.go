package v3

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"sb.im/gosd/app/model"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// @Summary Blobs Create
// @Schemes Blobs
// @Description create a new blobs
// @Tags blob
// @Accept json, multipart/form-data
// @Produce json
// @Param file body model.Blob true "Blob"
// @Success 200 {object} model.Blob
// @Router /blobs [post]
func (h *Handler) blobCreate(c *gin.Context) {
	bindBlob := make(map[string]string)
	if err := c.ShouldBind(&bindBlob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		for key, value := range form.File {
			for i, file := range value {
				blob := &model.Blob{
					UXID: xid.New().String(),
					Name: filepath.Base(file.Filename),
				}
				log.Infoln(key, i, blob)
				if i == 0 {
					bindBlob[key] = blob.UXID
				} else {
					bindBlob[key+"."+strconv.Itoa(i)] = blob.UXID
				}
				if err := c.SaveUploadedFile(file, h.cfg.StoragePath + blob.UXID); err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
					return
				}
				h.orm.Create(blob)
			}
		}
	c.JSON(http.StatusOK, bindBlob)
}


