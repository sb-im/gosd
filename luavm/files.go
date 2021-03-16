package luavm

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"

	"sb.im/gosd/model"
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

// FilesUrl(key)
// api/v1/plans/{planID}?files={filesKey}&token={token}
func (s *Service) FilesUrl(key string) string {
	token, _ := genToken(16)
	s.State.Record(fmt.Sprintf("token/%s", token), nil)
	return fmt.Sprintf(s.Task.URL, s.Task.PlanID, key, token)
}

// JobFilesUrl(key)
// api/v1/plans/{planID}/jobs/{jobID}?files={fileKey}&token={token}
func (s *Service) JobFilesUrl(key string) string {
	token, _ := genToken(16)
	s.State.Record(fmt.Sprintf("token/%s", token), nil)
	return fmt.Sprintf(s.Task.JobURL, s.Task.PlanID, s.Task.Job.JobID, key, token)
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
