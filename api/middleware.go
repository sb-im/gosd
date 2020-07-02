package api

import (
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

		key := strings.Split(req.Header.Get("Authorization"), " ")[1]
		user := h.store.GetCurrentUser(key)

		if user == nil {
			json.Unauthorized(w, req)
			return
		}

		next.ServeHTTP(w, req)
	})
}
