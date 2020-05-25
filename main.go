package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"fmt"
	"net/url"

	"sb.im/gosd/state"

	"github.com/gorilla/mux"
)

var (
	namespace   = "/gosd"
	api_version = "/api/v1"
	profix      = namespace + api_version
)

var accessGrant *AccessGrant

func main() {
	config_path := "./config.yaml"
	if os.Getenv("GOSD_CONF") != "" {
		config_path = os.Getenv("GOSD_CONF")
	}
	fmt.Println("load config: " + config_path)

	config, err := getConfig(config_path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("=========")

	uri, err := url.Parse(config.Mqtt)
	if err != nil {
		log.Fatal(err)
	}

	state := &state.State{
		Node: make(map[int]state.NodeState),
	}

	mq := state.Connect("cloud.0", uri)
	fmt.Println(mq)

	DBlink(config.Database)
	accessGrant = NewAccessGrant()
	r := mux.NewRouter()
	r.HandleFunc(namespace+"/oauth/token", oauthHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.HandleFunc(profix+"/plans/", planIndexHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(profix+"/plans/", planCreateHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc(profix+"/plans/{id}", planIDHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(profix+"/plans/{id}", planUpdateHandler).Methods(http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.HandleFunc(profix+"/plans/{id}", planDestroyHandler).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc(profix+"/plans/{id}/get_map", planMapFileDownloadHandler).Methods(http.MethodGet, http.MethodOptions)

	//r.Handle("/storage/{name}", http.StripPrefix("/storage/", http.FileServer(http.Dir("data/storage"))))

	r.HandleFunc(profix+"/{action}/", actionHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(config.Server, r)
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "content-type,Authorization")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	type Token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		CreatedAt   int64  `json:"created_at"`
	}

	key, _ := accessGrant.Grant(nil)
	log.Println(key)

	token := &Token{
		AccessToken: key,
		TokenType:   "bearer",
		ExpiresIn:   7200,
		CreatedAt:   time.Now().Unix(),
	}
	b, err := json.Marshal(token)
	if err != nil {
		log.Println("error:", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	w.Write(b)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	content, err := ioutil.ReadFile("data/" + vars["action"] + ".json")
	if err != nil {
		log.Fatal(err)
	}

	w.Write(content)
}
