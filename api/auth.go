package api

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"miniflux.app/http/response/json"
)

func (h *handler) authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "content-type,Authorization")
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

	h.store.CreateToken(key, nil)

	json.OK(w, r, token)
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
