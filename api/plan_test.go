package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"sb.im/gosd/database"
	"sb.im/gosd/storage"
)

func Test_CreatePlanjson(t *testing.T) {
	db, _ := database.NewConnectionPool("postgres://postgres:password@localhost/gosd?sslmode=disable", 1, 10)
	database.Migrate(db)
	store := storage.NewStorage(db)

	handler := &handler{store}
	ts := httptest.NewServer(http.HandlerFunc(handler.createPlan))
	defer ts.Close()

	data := `{"name": "233", "description": "test", "node_id": 1}`
	res, err := http.Post(ts.URL, "application/json", bytes.NewReader([]byte(data)))
	if err != nil {
		t.Error(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("%s\n", greeting)
	}
}

func Test_CreatePlanFormData(t *testing.T) {
	db, _ := database.NewConnectionPool("postgres://postgres:password@localhost/gosd?sslmode=disable", 1, 10)
	database.Migrate(db)
	store := storage.NewStorage(db)

	handler := &handler{store}
	ts := httptest.NewServer(http.HandlerFunc(handler.createPlan))
	defer ts.Close()

	file := bytes.NewReader([]byte(`2222222222233333333333333`))
	name := bytes.NewReader([]byte(`test`))

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer

	// Add an file
	fw, err := w.CreateFormFile("file", "233.txt")
	if err != nil {
		t.Error(err)
	}
	if _, err := io.Copy(fw, file); err != nil {
		t.Error(err)
	}

	if fw, err = w.CreateFormField("name"); err != nil {
		t.Error(err)
	}
	if _, err = io.Copy(fw, name); err != nil {
		t.Error(err)
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	req, err := http.NewRequest("PATCH", ts.URL, &b)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Error(err)
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		t.Errorf("%s\n", d)
	}
}
