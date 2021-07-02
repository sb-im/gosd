package api

import (
	"fmt"
	"net/http"
	"strings"

	"miniflux.app/http/response/json"
)

func CORSOriginMiddleware(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			next.ServeHTTP(w, req)
		})
	}
}

func (h handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/gosd/api/v2/ok") {
			next.ServeHTTP(w, req)
			return
		}

		if strings.HasPrefix(req.URL.Path, "/gosd/api/v1/blobs/") {
			next.ServeHTTP(w, req)
			return
		}

		if strings.HasPrefix(req.URL.Path, "/gosd/api/v2/oauth/token") {
			next.ServeHTTP(w, req)
			return
		}

		if h.helpCurrentUser(w, req) == nil {
			json.Unauthorized(w, req)
			return
		}

		next.ServeHTTP(w, req)
	})
}


func (h handler) Oauth2Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userInfo, err := h.oauth.ValidationBearerToken(r); err != nil {
			fmt.Println(userInfo)
			return
		} else {
			fmt.Println("userID: ", userInfo.GetUserID())
		}
		next.ServeHTTP(w, r)
	})
}
