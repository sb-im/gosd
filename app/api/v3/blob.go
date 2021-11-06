package v3

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"sb.im/gosd/app/model"

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
			if err := c.SaveUploadedFile(file, h.cfg.StoragePath+blob.UXID); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
			h.orm.Create(blob)
		}
	}
	c.JSON(http.StatusOK, bindBlob)
}

func (h *Handler) blobUpdate(c *gin.Context) {
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

	// blobID: url, key
	// squash key / value
	updateBlob := make(map[string]interface{})

	blobID := c.Query("blobID")
	for key, value := range form.File {
		if !h.blobIsExist(key) {
			if blobID != "" {
				// TODO: url param
				// if url param h.blobIsExist
				if h.blobIsExist(blobID) {
					updateBlob[blobID] = value[0]

					// This blobID only use once
					blobID = ""
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "blobID error"})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "blobID error"})
				return
			}
		}
		updateBlob[key] = value[0]
	}

	for key, value := range updateBlob {
		bindBlob[key] = "ok"
		// TODO: update blob
		log.Warnln("needs implement update blob", value)
	}

	c.JSON(http.StatusOK, bindBlob)
}

func (h *Handler) blobIsExist(id string) bool {
	var count int64
	h.orm.Model(&model.Blob{}).Where("uxid = ?", id).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func (h *Handler) blobShow(c *gin.Context) {
	blob := model.Blob{}
	h.orm.Take(&blob, "uxid = ?", c.Param("blobID"))
	if blob.ID != 0 {
		c.FileAttachment(h.cfg.StoragePath+blob.UXID, blob.Name)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NotFound this blob"})
		return
	}
}
