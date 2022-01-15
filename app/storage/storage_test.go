package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	key := "gosd_storage_test_tmp"
	path := "/tmp"
	data := []byte(`xxxxxxx`)
	s := NewStorage(path)
	if err := s.Set(key, data); err != nil {
		t.Error(err)
	}

	if fData, err := s.Get(key); err != nil {
		t.Error(err)
	} else {
		if string(data) != string(fData) {
			t.Errorf("%s\n", fData)
		}
	}

	if err := s.Del(key); err != nil {
		t.Error(err)
	}

	if _, err := s.Get(key); err == nil {
		t.Error("Should Not File")
	}
}
