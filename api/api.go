package api

import (
	"net/http"
	"net/url"

	"sb.im/gosd/luavm"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	ostore "github.com/go-oauth2/oauth2/v4/store"

	"github.com/gorilla/mux"
)

type ServeConfig struct {
	BaseURL string
	OauthID string
	OauthSecret string
}

// Serve declares API routes for the application.
func Serve(router *mux.Router, cache *state.State, store *storage.Storage, worker *luavm.Worker, opt ServeConfig) {
	u, _ := url.Parse(opt.BaseURL)

	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(ostore.NewMemoryTokenStore())

	// client memory store
	clientStore := ostore.NewClientStore()
	clientStore.Set(opt.OauthID, &models.Client{
		ID:     opt.OauthID,
		Secret: opt.OauthSecret,
		Domain: opt.BaseURL,
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)

	handler := &handler{cache, srv, store, worker, opt.BaseURL}
	handler.oauthInit()
	sr := router.PathPrefix(u.Path + "/api/v1").Subrouter()
	//middleware := newMiddleware(store)

	//sr.Use(middleware.apiKeyAuth)

	sr.Use(CORSOriginMiddleware("*"))
	sr.Use(handler.AuthMiddleware)

	sr2 := router.PathPrefix(u.Path + "/api/v2").Subrouter()
	sr2.Use(CORSOriginMiddleware("*"))
	sr2.Use(handler.AuthMiddleware)


	// new oauth2
	// Only use Oauth password authentication
	// sr2.HandleFunc("/oauth/authorize", handler.authorizeHandler)
	sr2.HandleFunc("/oauth/token", handler.oAuthHandler).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc(u.Path+"/oauth/token", handler.authHandler).Methods(http.MethodPost, http.MethodOptions)

	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")

		w.Header().Set("Access-Control-Allow-Methods",
			http.MethodPut + "," +
			http.MethodPost + "," +
			http.MethodPatch + "," +
			http.MethodDelete)
	})

	// Plan
	plan := sr2.PathPrefix("/plans").Subrouter()
	plan.HandleFunc("/", handler.plans).Methods(http.MethodGet)
	plan.HandleFunc("/", handler.createPlan).Methods(http.MethodPost)
	plan.HandleFunc("/{planID:[0-9]+}", handler.updatePlan).Methods(http.MethodPut)
	plan.HandleFunc("/{planID:[0-9]+}", handler.planDestroy).Methods(http.MethodDelete)

	plan.HandleFunc("/{planID:[0-9]+}/jobs/", handler.planLogs).Methods(http.MethodGet)
	plan.HandleFunc("/{planID:[0-9]+}/jobs/", handler.createPlanLog).Methods(http.MethodPost)
	plan.HandleFunc("/{planID:[0-9]+}/jobs/{logID:[0-9]+}/cancel", handler.planRunningDestroy).Methods(http.MethodPost)

	plan.HandleFunc("/{planID:[0-9]+}/running", handler.createPlanLog).Methods(http.MethodPost)
	plan.HandleFunc("/{planID:[0-9]+}/running", handler.planRunning).Methods(http.MethodGet)
	plan.HandleFunc("/{planID:[0-9]+}/running", handler.planRunningDestroy).Methods(http.MethodDelete)

	sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/", handler.planLogs).Methods(http.MethodGet)
	sr.HandleFunc("/plans/{planID:[0-9]+}/plan_logs/", handler.createPlanLog).Methods(http.MethodPost)

	sr.HandleFunc("/plans/{planID:[0-9]+}/logs/", handler.planLogs).Methods(http.MethodGet)
	sr.HandleFunc("/plans/{planID:[0-9]+}/logs/", handler.createPlanLog).Methods(http.MethodPost)

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

	sr2.PathPrefix("/mqtt/").HandlerFunc(handler.mqttGet).Methods(http.MethodGet)
	sr2.PathPrefix("/mqtt/").HandlerFunc(handler.mqttPut).Methods(http.MethodPost)
}
