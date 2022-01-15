package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	key := "gosd_storage_test_tmp"
	path := "/tmp"
	name := "test.txt"
	data := []byte(`xxxxxxx`)
	s := NewStorage(path)
	if err := s.Set(key, name, data); err != nil {
		t.Error(err)
	}

	if fName, fData, err := s.Get(key); err != nil {
		t.Error(err)
	} else {
		if name != fName {
			//t.Error(fName)
		}

		if string(data) != string(fData) {
			t.Errorf("%s\n", fData)
		}
	}

	if err := s.Del(key); err != nil {
		t.Error(err)
	}

	if _, _, err := s.Get(key); err == nil {
		t.Error("Should Not File")
	}
}
