package api

import (
	"errors"
	//"fmt"
	"net/http"
	"strconv"
	//"strings"

	"sb.im/gosd/auth"
	"sb.im/gosd/model"

	//"github.com/gomodule/redigo/redis"

	"miniflux.app/http/response/json"
)

func (h *handler) currentUser(w http.ResponseWriter, r *http.Request) {
	user := h.helpCurrentUser(w, r)
	if user == nil {
		json.ServerError(w, r, errors.New("No login"))
		return
	}

	json.OK(w, r, user)
}

func (h *handler) helpCurrentUser(w http.ResponseWriter, r *http.Request) *model.User {
	// Single user mode
	// User.ID == 1
	if auth.AuthMethod == auth.NoAuth {
		user, _ := h.store.UserByID(1)
		return user
	}

	info, _ := h.oauth.ValidationBearerToken(r)
	userID, _ := strconv.ParseInt(info.GetUserID(), 10, 0)
	user, _ := h.store.UserByID(userID)
	return user

	//token := strings.Split(r.Header.Get("Authorization"), " ")
	//if len(token) == 2 {
	//	key := token[1]
	//	if user := h.store.GetCurrentUser(key); user != nil {
	//		return user
	//	} else {
	//		uniquekey := fmt.Sprintf("users/token/%s", key)
	//		if userID, err := redis.Int64(h.cache.Do("GET", uniquekey)); err != nil {
	//			return user
	//		} else {
	//			user, _ := h.store.UserByID(userID)
	//			if user != nil {
	//				h.store.CreateToken(key, user)
	//			}
	//			return user
	//		}
	//	}
	//}
	return nil
}
