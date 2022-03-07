package v3

import (
	"net/http"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"sb.im/gosd/app/model"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// @Summary Blobs Create
// @Schemes Blobs
// @Description create a new blobs, return map[string]string, key is key, value is blobID
// @Tags blob
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "this is a file"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /blobs [POST]
func (h *Handler) BlobCreate(c *gin.Context) {
	bindBlob := make(map[string]string)
	if err := c.ShouldBind(&bindBlob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for key, value := range form.File {
		for i, file := range value {
			uxid, err := uuid.NewV4()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			blob := &model.Blob{
				UXID: uxid.String(),
				Name: filepath.Base(file.Filename),
			}
			log.Infoln(key, i, blob)
			if i == 0 {
				bindBlob[key] = blob.UXID
			} else {
				bindBlob[key+"."+strconv.Itoa(i)] = blob.UXID
			}

			if err := c.SaveUploadedFile(file, h.ofs.LocalPath(blob.UXID)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if err := h.orm.Create(blob).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
	c.JSON(http.StatusOK, bindBlob)
}

// @Summary Blobs Update
// @Schemes Blobs
// @Description create a new blobs
// @Tags blob
// @Accept multipart/form-data
// @Produce json
// @Param blobID path string true "blob ID"
// @Param file formData file true "this is a file"
// @Success 200 {object} model.Blob
// @Failure 400
// @Failure 500
// @Router /blobs/{blobID} [PUT]

func (h *Handler) BlobUpdate(c *gin.Context) {
	bindBlob := make(map[string]string)
	if err := c.ShouldBind(&bindBlob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// blobID: url, key
	// squash key / value
	updateBlob := make(map[string]interface{})

	blobID := c.Param("blobID")
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

// @Summary Blobs Get
// @Schemes Blobs
// @Description get a blob content
// @Tags blob
// @Accept multipart/form-data
// @Produce json
// @Param blobID path string true "blob ID"
// @Success 200 {object} model.Blob
// @Success 404
// @Router /blobs/{blobID} [GET]
func (h *Handler) BlobShow(c *gin.Context) {
	blob := model.Blob{}
	h.orm.Take(&blob, "uxid = ?", c.Param("blobID"))
	if blob.ID != 0 {
		c.FileAttachment(h.ofs.LocalPath(blob.UXID), blob.Name)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "NotFound this blob"})
		return
	}
}
