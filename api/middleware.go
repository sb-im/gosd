package api

import (
	"net/http"
)

func CORSOriginMiddleware(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			next.ServeHTTP(w, req)
		})
	}
}
