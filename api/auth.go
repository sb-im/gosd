package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"sb.im/gosd/model"

	"miniflux.app/http/response/json"
)

const lockTimeOut = 1 * time.Hour
const maxRetry = 5

var lock struct {
	fail map[time.Time]int64
}

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

	user, err := h.store.UserByUsername(item.Username)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	// TODO: This UserByUsername sql.ErrNoRows not return error
	// Maybe this need Change
	if user == nil {
		json.ServerError(w, r, errors.New("Not user found"))
		return
	}

	if len(lock.fail) == 0 {
		lock.fail = make(map[time.Time]int64, 1024)
	}

	count := 0
	t := time.Now()
	for k, v := range lock.fail {
		if k2 := k.Add(lockTimeOut); t.Unix() > k2.Unix() {
			delete(lock.fail, k)
		}
		if v == user.ID {
			count = count + 1
		}
	}
	if user != nil {
		fmt.Println("Lock User: ", user.ID)
	}
	if count > maxRetry {
		return
	}

	if err := h.store.CheckPassword(item.Username, item.Password); err != nil {
		json.ServerError(w, r, err)

		lock.fail[time.Now()] = user.ID
		fmt.Println("ERROR Login user: ", user.ID)
		return
	}

	key, err := genToken(32)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	expire := 7200
	token := &model.Token{
		AccessToken: key,
		TokenType:   "bearer",
		ExpiresIn:   expire,
		CreatedAt:   time.Now().Unix(),
	}

	h.store.CreateToken(key, user)

	uniquekey := fmt.Sprintf("users/token/%s", key)
	h.cache.Do("SET", uniquekey, user.ID)
	h.cache.Do("EXPIRE", uniquekey, expire)

	json.OK(w, r, token)
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
