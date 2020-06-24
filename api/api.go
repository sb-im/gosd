package api

import (
	"net/http"
	"net/url"
	"strings"

	"sb.im/gosd/luavm"
	"sb.im/gosd/storage"

	"github.com/gorilla/mux"
)

// Serve declares API routes for the application.
func Serve(router *mux.Router, store *storage.Storage, worker *luavm.Worker, baseURL string) {
	u, _ := url.Parse(baseURL)

	handler := &handler{store, worker, baseURL}
	sr := router.PathPrefix(u.Path + "/api/v1").Subrouter()

	//middleware := newMiddleware(store)

	//sr.Use(middleware.apiKeyAuth)

	sr.Use(CORSOriginMiddleware("*"))

	//router.Use(mux.CORSMethodMiddleware(sr))

	router.PathPrefix(u.Path + "/api/v1").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	sr.HandleFunc("/plans/{planID:[0-9]+}/mission_queues/", handler.missionQueue).Methods(http.MethodGet)

	sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/", handler.createPlanLog).Methods(http.MethodPost)
	sr.HandleFunc("/plans/{planID:[0-9]+}/logs/", handler.createPlanLog).Methods(http.MethodPost)
	//sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/{logID:[0-9]+}/", handler.createPlanLog).Methods(http.MethodPost)

	sr.HandleFunc("/blobs/{blobID:[0-9]+}", handler.blobByID).Methods(http.MethodGet)
}
