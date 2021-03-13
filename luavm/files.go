package luavm

import (
	"bytes"
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
