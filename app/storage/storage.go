package storage

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Storage struct {
	path string
}

func (s *Storage) LocalPath(name string) string {
	return path.Join(s.path, name)
}

func NewStorage(storageURL string) *Storage {
	path := "data/storage"
	p := strings.SplitN(storageURL, "://", 2)
	if len(p) > 1 {
		path = p[1]
	}
	os.MkdirAll(path, 0755)
	return &Storage{path: path}
}

func (s Storage) Set(key string, data []byte) error {
	return ioutil.WriteFile(path.Join(s.path, key), data, 0644)
}

func (s Storage) Get(key string) ([]byte, error) {
	data, err := ioutil.ReadFile(path.Join(s.path, key))
	return data, err
}

func (s Storage) Del(key string) error {
	return os.Remove(path.Join(s.path, key))
}
