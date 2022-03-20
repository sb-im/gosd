package luavm

import (
	"fmt"
	"os"
	"time"

	"sb.im/gosd/app/helper"
	"sb.im/gosd/app/model"
)

const (
	blob_url = "/blobs/%s?token=%s"
)

// Reader
func (s Service) BlobReader(id string) (string, string) {
	blob := model.Blob{}
	if err := s.orm.First(&blob, "uxid = ?", id); err != nil {
		return "", ""
	}
	data, _ := s.ofs.Get(blob.UXID)
	return blob.Name, string(data)
}

// Update
func (s Service) BlobUpdate(id, filename, content string) {
	blob := model.Blob{}
	if err := s.orm.First(&blob, "uxid = ?", id); err != nil {
		return
	}

	if blob.Name != filename {
		s.orm.Updates(&blob)
	}

	s.ofs.Set(blob.UXID, []byte(content))
}

// Create
func (s Service) BlobCreate(filename, content string) string {
	blob := model.Blob{
		Name: filename,
	}
	s.orm.Save(&blob)

	s.ofs.Set(blob.UXID, []byte(content))
	return blob.UXID
}

// Delete
// TODO: need storage blobs delete

// TODO: test token
// {BASE_URL}/blobs/{UUID}?token={token}
func (s Service) BlobUrl(blobID string) string {
	token := helper.GenSecret(16)

	// TODO: need common lib
	s.rdb.Set(s.ctx, fmt.Sprintf("token/%s", token), s.Task.TeamID, 2*time.Hour)
	return fmt.Sprintf(os.Getenv("BASE_URL")+blob_url, blobID, token)
}
