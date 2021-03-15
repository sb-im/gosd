package luavm

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"strconv"
)

func (s *Service) getBlob(id string) (string, string) {
	blobID, _ := strconv.ParseInt(id, 10 ,64)
	blob, _ := s.Store.BlobByID(blobID)
  content, _ := ioutil.ReadAll(blob.Reader)
	return blob.FileName, string(content)
}

func (s *Service) setBlob(id, filename, content string) {
	blobID, _ := strconv.ParseInt(id, 10 ,64)
	blob, _ := s.Store.BlobByID(blobID)
	blob.FileName = filename
	blob.Reader =  bytes.NewReader([]byte(content))
	s.Store.UpdateBlob(blob)
}

// TODO: Need Fix s.Task.Files == nil

// GetFilesContent(key) filename, content
func (s *Service) GetFilesContent(key string) (string, string) {
	return s.getBlob(s.Task.Files[key])
}

// SetFilesContent(key, filename, content)
func (s *Service) SetFilesContent(key, filename, content string) {
	s.setBlob(s.Task.Files[key], filename, content)
}

// GetJobFilesContent(key) filename, content
func (s *Service) GetJobFilesContent(key string) (string, string) {
	return s.getBlob(s.Task.Job.Files[key])
}

// SetJobFilesContent(key, filename, content)
func (s *Service) SetJobFilesContent(key, filename, content string) {
	s.setBlob(s.Task.Job.Files[key], filename, content)
}

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
