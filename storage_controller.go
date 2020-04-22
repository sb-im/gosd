package main

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func newBlob(name string, contents []byte) (string, error) {
	blob := &StorageBlob{
		Filename: name,
	}
	db.Create(blob)
	//fmt.Printf("%s\n", strconv.FormatUint(uint64(blob.ID), 10))

	//filepath := "data/storage/" + name
	filepath := "data/storage/" + strconv.FormatUint(uint64(blob.ID), 10)
	err := ioutil.WriteFile(filepath, contents, 0644)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

func file2path(r *http.Request) map[string]string {
	params := make(map[string]string)

	mediaType, mimeParams, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, mimeParams["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}

			if p.FileName() != "" {
				//fmt.Println(p.FileName())
				//filepath := "data/storage/" + p.FileName()
				//err := ioutil.WriteFile(filepath, slurp, 0644)
				filepath, err := newBlob(p.FileName(), slurp)
				if err != nil {
					log.Fatal(err)
				}
				params[p.FormName()] = filepath
			} else {
				params[p.FormName()] = string(slurp)
			}

		}
	}

	return params
}

func handleDownload(filename, filepath string, w http.ResponseWriter) {
	file, err := os.Open(filepath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	defer file.Close()

	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}
