package luavm

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"

	"sb.im/gosd/model"
)

const (
	blob_url = "api/v1/blobs/%s"
)

// Reader
func (s *Service) BlobReader(id int64) (string, string) {
	blob, _ := s.Store.BlobByID(id)
	content, _ := ioutil.ReadAll(blob.Reader)
	return blob.FileName, string(content)
}

// Update
func (s *Service) BlobUpdate(id int64, filename, content string) {
	blob, _ := s.Store.BlobByID(id)
	blob.FileName = filename
	blob.Reader = bytes.NewReader([]byte(content))
	s.Store.UpdateBlob(blob)
}

// Create
func (s *Service) BlobCreate(filename, content string) int64 {
	blob := model.NewBlob(filename, bytes.NewReader([]byte(content)))
	s.Store.CreateBlob(blob)
	return blob.ID
}

// Delete
// TODO: need storage blobs delete

// TODO: test token
// api/v1/plans/{planID}?token={token}
func (s *Service) BlobUrl(blobID string) string {
	token, _ := genToken(16)
	s.State.Record(fmt.Sprintf("token/%s", token), nil)
	return fmt.Sprintf(os.Getenv("BASE_URL")+blob_url, blobID)
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
