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
	sr.Use(handler.AuthMiddleware)

	sr2 := router.PathPrefix(u.Path + "/api/v2").Subrouter()
	sr2.Use(CORSOriginMiddleware("*"))
	sr2.Use(handler.AuthMiddleware)

	//router.Use(mux.CORSMethodMiddleware(sr))
	router.HandleFunc(u.Path+"/oauth/token", handler.authHandler).Methods(http.MethodPost, http.MethodOptions)

	router.PathPrefix(u.Path + "/api").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		allMethods := []string{
			http.MethodPatch,
			http.MethodDelete,
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allMethods, ","))
	})

	sr2.HandleFunc("/plans/", handler.createPlan2).Methods(http.MethodPost)
	sr2.HandleFunc("/plans/", handler.plan2s).Methods(http.MethodGet)
	sr2.HandleFunc("/plans/{planID:[0-9]+}", handler.updatePlan2).Methods(http.MethodPatch, http.MethodPut)
	sr2.HandleFunc("/plans/{planID:[0-9]+}", handler.planDestroy).Methods(http.MethodDelete)

	sr.HandleFunc("/plans/", handler.plans).Methods(http.MethodGet)
	sr.HandleFunc("/plans/", handler.createPlan).Methods(http.MethodPost)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planByID).Methods(http.MethodGet)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planUpdate).Methods(http.MethodPatch, http.MethodPut)
	sr.HandleFunc("/plans/{planID:[0-9]+}", handler.planDestroy).Methods(http.MethodDelete)

	sr.HandleFunc("/plans/{planID:[0-9]+}/mission_queues/", handler.missionQueue).Methods(http.MethodGet)

	// How is this API designed WTF ???
	sr.HandleFunc("/mission_queues/plan/{planID:[0-9]+}", handler.missionStop).Methods(http.MethodDelete)

	sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/", handler.planLogs).Methods(http.MethodGet)
	sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/", handler.createPlanLog).Methods(http.MethodPost)

	sr.HandleFunc("/plans/{planID:[0-9]+}/logs/", handler.planLogs).Methods(http.MethodGet)
	sr.HandleFunc("/plans/{planID:[0-9]+}/logs/", handler.createPlanLog).Methods(http.MethodPost)

	sr2.HandleFunc("/plans/{planID:[0-9]+}/jobs/", handler.planLogs).Methods(http.MethodGet)
	sr2.HandleFunc("/plans/{planID:[0-9]+}/jobs/", handler.createPlanLog).Methods(http.MethodPost)
	sr2.HandleFunc("/plans/{planID:[0-9]+}/jobs/{logID:[0-9]+}", handler.createPlanLog).Methods(http.MethodGet)
	sr2.HandleFunc("/plans/{planID:[0-9]+}/jobs/{logID:[0-9]+}/cancel", handler.missionStop).Methods(http.MethodPost)

	// Debug use
	sr2.HandleFunc("/plans/{planID:[0-9]+}/jobs/running", handler.missionStop).Methods(http.MethodDelete)
	//sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/{logID:[0-9]+}/", handler.createPlanLog).Methods(http.MethodPost)


	sr2.HandleFunc("/blobs/", handler.createBlob).Methods(http.MethodPost)
	sr2.HandleFunc("/blobs/{blobID:[0-9]+}", handler.blobByID).Methods(http.MethodGet)
	sr2.HandleFunc("/blobs/{blobID:[0-9]+}", handler.updateBlob).Methods(http.MethodPatch, http.MethodPut)
	sr2.HandleFunc("/blobs/{blobID:[0-9]+}", handler.destroyBlob).Methods(http.MethodDelete)

	sr.HandleFunc("/blobs/", handler.createBlob).Methods(http.MethodPost)
	sr.HandleFunc("/blobs/{blobID:[0-9]+}", handler.blobByID).Methods(http.MethodGet)
	sr.HandleFunc("/blobs/{blobID:[0-9]+}", handler.updateBlob).Methods(http.MethodPatch, http.MethodPut)
	sr.HandleFunc("/blobs/{blobID:[0-9]+}", handler.destroyBlob).Methods(http.MethodDelete)

	sr.HandleFunc("/user/", handler.currentUser).Methods(http.MethodGet)
	sr.HandleFunc("/{action}/", handler.actionHandler).Methods(http.MethodGet)

	sr2.HandleFunc("/ok", handler.ok).Methods(http.MethodGet)
}
