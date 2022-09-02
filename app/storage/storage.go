package storage

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

type Storage struct {
	path string
}

func (s *Storage) LocalPath(name string) string {
	return path.Join(s.path, name)
}

func NewStorage(storageURL string) *Storage {
	u, err := url.Parse(storageURL)
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(u.Path, 0755); err != nil {
		panic(err)
	}
	return &Storage{path: u.Path}
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
