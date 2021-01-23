package api

import (
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"sb.im/gosd/model"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

func handleDownload(filename string, reader io.Reader, w http.ResponseWriter) {
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	_, err := io.Copy(w, reader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}

func (h *handler) createBlob(w http.ResponseWriter, r *http.Request) {
	_, files, err := h.formData2Blob2(r)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	mapBlob := make(map[string]string)
	for name, file := range files {
		blob := model.NewBlob(file.Name, file.Reader)
		if err := h.store.CreateBlob(blob); err != nil {
			json.ServerError(w, r, err)
		}
		mapBlob[name] = strconv.FormatInt(blob.ID, 10)
	}

	json.OK(w, r, mapBlob)
}

func (h *handler) blobByID(w http.ResponseWriter, r *http.Request) {
	id := request.RouteInt64Param(r, "blobID")
	item, err := h.store.BlobByID(id)
	if err != nil {
		json.BadRequest(w, r, errors.New("Unable to fetch this blob from the database"))
		return
	}

	if item == nil {
		json.NotFound(w, r)
		return
	}

	//json.OK(w, r, item)
	handleDownload(item.FileName, item.Reader, w)
}

func (h *handler) updateBlob(w http.ResponseWriter, r *http.Request) {
	_, files, err := h.formData2Blob2(r)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}
	for _, file := range files {
		id := request.RouteInt64Param(r, "blobID")
		blob, err := h.store.BlobByID(id)
		if err != nil {
			json.BadRequest(w, r, errors.New("Unable to fetch this blob from the database"))
			return
		}

		if blob == nil {
			json.NotFound(w, r)
			return
		}

		blob.FileName = file.Name
		blob.Reader = file.Reader

		if err := h.store.UpdateBlob(blob); err != nil {
			json.ServerError(w, r, err)
		}

		// Only First: range is random
		json.OK(w, r, blob)
		return
	}
	json.ServerError(w, r, errors.New("not found content"))
}

// -> extra, files, error
func (h *handler) formData2Blob2(r *http.Request) (map[string]string, map[string]struct {
	Name   string
	Reader io.Reader
}, error) {
	extra := make(map[string]string)
	files := make(map[string]struct {
		Name   string
		Reader io.Reader
	})
	mediaType, mimeParams, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return extra, files, err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, mimeParams["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}

			if p.FileName() != "" {
				files[p.FormName()] = struct {
					Name   string
					Reader io.Reader
				}{
					Name:   p.FileName(),
					Reader: p,
				}
			} else {
				slurp, err := ioutil.ReadAll(p)
				if err != nil {
					return extra, files, err
				}
				extra[p.FormName()] = string(slurp)
			}

		}
	}
	return extra, files, err
}

// -> params, blobFileID, error
func (h *handler) formData2Blob(r *http.Request) (map[string]string, map[string]string, error) {
	params := make(map[string]string)
	file := make(map[string]string)

	mediaType, mimeParams, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return params, file, err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, mimeParams["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}

			if p.FileName() != "" {
				blob := model.NewBlob(p.FileName(), p)

				err = h.store.CreateBlob(blob)
				if err != nil {
					return params, file, err
				}
				file[p.FormName()] = strconv.FormatInt(blob.ID, 10)
			} else {
				slurp, err := ioutil.ReadAll(p)
				if err != nil {
					return params, file, err
				}
				params[p.FormName()] = string(slurp)
			}

		}
	}

	return params, file, nil
}

func (h *handler) destroyBlob(w http.ResponseWriter, r *http.Request) {
	blobID := request.RouteInt64Param(r, "blobID")
	// TODO: Need Add destroyBlob
	json.OK(w, r, blobID)
}
