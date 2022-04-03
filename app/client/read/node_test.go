package read

import (
	"os"
	"testing"
)

const (
	testPathPrefix = "/tmp/gosd_test"
)

func TestGetPoints(t *testing.T) {
	os.MkdirAll(testPathPrefix+"/points", 0755)
	defer os.Remove(testPathPrefix)
	pointfile1 := testPathPrefix + "/points/test-1.json"
	pointfile2 := testPathPrefix + "/points/test-2.json"

	os.Create(pointfile1)
	if err := os.WriteFile(pointfile1, []byte(`{"type": "test-1"}`), 0644); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(pointfile2, []byte(`{"type": "test-2"}`), 0644); err != nil {
		t.Error(err)
	}

	points := getPoints(testPathPrefix + "/points")
	if len(points) != 2 {
		t.Errorf("This point number : %d", len(points))
	}
}
