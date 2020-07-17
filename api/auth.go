package api

import (
	"crypto/rand"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"sb.im/gosd/model"

	"miniflux.app/http/response/json"
)

func (h *handler) authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "content-type,Authorization")
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	item := &model.Login{}
	if strings.HasPrefix(mediaType, "multipart/") {

		params, _, err := h.formData2Blob(r)
		if err != nil {
			json.ServerError(w, r, err)
			return
		}

		item.GrantType = params["grant_type"]
		item.Username = params["username"]
		item.Password = params["password"]
		item.ClientID = params["client_id"]
		item.ClientSecret = params["client_secret"]

	} else {
		item, err = decodeLoginPayload(r.Body)
		if err != nil {
			json.BadRequest(w, r, err)
			return
		}
	}

	if err := h.store.CheckPassword(item.Username, item.Password); err != nil {
		json.ServerError(w, r, err)
		return
	}

	user, err := h.store.UserByUsername(item.Username)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	type Token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		CreatedAt   int64  `json:"created_at"`
	}

	key, err := genToken(32)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	token := &Token{
		AccessToken: key,
		TokenType:   "bearer",
		ExpiresIn:   7200,
		CreatedAt:   time.Now().Unix(),
	}

	h.store.CreateToken(key, user)

	json.OK(w, r, token)
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}