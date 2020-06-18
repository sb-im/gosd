package api

import (
	"net/http"
	"strings"

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

	//router.Use(mux.CORSMethodMiddleware(sr))

	router.PathPrefix("/gosd/api/v1").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		allMethods := []string{
			http.MethodPatch,
			http.MethodDelete,
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allMethods, ","))
	})

	sr.HandleFunc("/plans/", handler.plans).Methods(http.MethodGet)
	sr.HandleFunc("/plans/", handler.createPlan).Methods(http.MethodPost)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planByID).Methods(http.MethodGet)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planUpdate).Methods(http.MethodPatch, http.MethodPut)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planDestroy).Methods(http.MethodDelete)

	sr.HandleFunc("/blobs/{blobID:[0-9]+}", handler.blobByID).Methods(http.MethodGet)
}
