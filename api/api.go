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

	router.PathPrefix("/gosd/api/v1").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})

	sr.HandleFunc("/plans/", handler.plans).Methods(http.MethodGet)
	sr.HandleFunc("/plans/", handler.createPlan).Methods(http.MethodPost)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planByID).Methods(http.MethodGet)

	sr.HandleFunc("/blobs/{blobID:[0-9]+}", handler.blobByID).Methods(http.MethodGet)
}
