package api

import (
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"sb.im/gosd/model"

	"miniflux.app/http/response/json"
)

func (h *handler) createBlob(w http.ResponseWriter, r *http.Request) {
	plan, err := decodePlanCreationPayload(r.Body)
	if err != nil {
		json.BadRequest(w, r, err)
		return
	}

	err = h.store.CreatePlan(plan)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, plan)
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
