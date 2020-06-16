package api

import (
	"sb.im/gosd/storage"

	"github.com/gorilla/mux"
)

// Serve declares API routes for the application.
func Serve(router *mux.Router, store *storage.Storage) {
	handler := &handler{store}
	sr := router.PathPrefix("/v1").Subrouter()

	//middleware := newMiddleware(store)

	//sr.Use(middleware.apiKeyAuth)

	//sr.Use(middleware.basicAuth)

	sr.HandleFunc("/plans", handler.createPlan).Methods("POST")
}
