package luavm

import (
	"crypto/rand"
	"fmt"
	"os"

	"sb.im/gosd/app/model"
)

const (
	blob_url = "api/v1/blobs/%s"
)

// Reader
func (s Service) BlobReader(id int64) (string, string) {
	blob := model.Blob{}
	if err := s.orm.First(&blob, "uxid = ?", id); err != nil {
		return "", ""
	}
	data, _ := s.ofs.Get(blob.UXID)
	return blob.Name, string(data)
}

// Update
func (s Service) BlobUpdate(id int64, filename, content string) {
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
func (s Service) BlobCreate(filename, content string) int64 {
	blob := model.Blob{
		Name: filename,
	}
	s.orm.Save(&blob)

	s.ofs.Set(blob.UXID, []byte(content))
	return int64(blob.ID)
}

// Delete
// TODO: need storage blobs delete

// TODO: test token
// api/v1/plans/{planID}?token={token}
func (s Service) BlobUrl(blobID string) string {
	token, _ := genToken(16)
	s.rdb.Set(s.ctx, fmt.Sprintf("token/%s", token), nil, 0)
	return fmt.Sprintf(os.Getenv("BASE_URL")+blob_url, blobID)
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
