package storage

import (
	"io/ioutil"
	"os"
	"path"
)

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
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
