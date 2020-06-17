package api

import (
	"net/http"

	"sb.im/gosd/storage"

	"github.com/gorilla/mux"
)

// Serve declares API routes for the application.
func Serve(router *mux.Router, store *storage.Storage) {
	handler := &handler{store}
	sr := router.PathPrefix("/gosd/api/v1").Subrouter()

	//middleware := newMiddleware(store)

	//sr.Use(middleware.apiKeyAuth)

	sr.Use(CORSOriginMiddleware("*"))

	sr.HandleFunc("/plans/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	}).Methods(http.MethodOptions)

	sr.HandleFunc("/plans/", handler.plans).Methods(http.MethodGet)
	sr.HandleFunc("/plans/", handler.createPlan).Methods(http.MethodPost)
}
