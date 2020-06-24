package main

//go:generate go run generate.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"sb.im/gosd/config"
	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/luavm"

	"sb.im/gosd/state"

	"sb.im/gosd/database"
	"sb.im/gosd/storage"

	"sb.im/gosd/api"

	"github.com/gorilla/mux"
)

var (
	namespace   = "/gosd"
	api_version = "/api/v1"
	profix      = namespace + api_version
)

var accessGrant *AccessGrant

func main() {
	parse := config.NewParser()
	opts, err := parse.ParseEnvironmentVariables()
	if err != nil {
		fmt.Println(err)
	}

	db, err := database.NewConnectionPool(
		opts.DatabaseURL(),
		opts.DatabaseMinConns(),
		opts.DatabaseMaxConns(),
	)

	if err != nil {
		fmt.Println(err)
	}

	database.Migrate(db)
	store := storage.NewStorage(db)

	fmt.Println("=========")

	uri, err := url.Parse(opts.MqttURL())
	if err != nil {
		log.Fatal(err)
	}

	state := state.NewState()
	mqttClient := state.Connect("cloud.0", uri)
	fmt.Println(mqttClient)

	time.Sleep(1 * time.Second)
	mqttProxy, err := jsonrpc2mqtt.OpenMqttProxy(mqttClient)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mqttProxy)

	// Wait mqtt connected
	time.Sleep(3 * time.Second)

	worker := luavm.NewWorker(state)
	go worker.Run()

	accessGrant = NewAccessGrant()
	r := mux.NewRouter()

	fmt.Println("Start http")
	api.Serve(r, store, worker, opts.BaseURL())

	r.HandleFunc(namespace+"/oauth/token", oauthHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)

	r.HandleFunc(profix+"/{action}/", actionHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(opts.ListenAddr(), r)
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
