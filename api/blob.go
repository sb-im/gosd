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

func (h *handler) blobByID(w http.ResponseWriter, r *http.Request) {
	id := request.RouteInt64Param(r, "blobID")
	item, err := h.store.BlobByID(id)
	if err != nil {
		json.BadRequest(w, r, errors.New("Unable to fetch this plan from the database"))
		return
	}

	if item == nil {
		json.NotFound(w, r)
		return
	}

	//json.OK(w, r, item)
	handleDownload(item.FileName, item.Reader, w)
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
