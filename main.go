package main

//go:generate go run generate.go

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"sb.im/gosd/api"
	"sb.im/gosd/cli"
	"sb.im/gosd/database"
	"sb.im/gosd/luavm"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	"github.com/gorilla/mux"
)

var (
	namespace   = "/gosd"
	api_version = "/api/v1"
	profix      = namespace + api_version
)

func main() {
	opts, err := cli.Parse()
	if err != nil {
		panic(err)
	}

	db, err := database.NewConnectionPool(
		opts.DatabaseURL(),
		opts.DatabaseMinConns(),
		opts.DatabaseMaxConns(),
	)

	if err != nil {
		panic(err)
	}

	store := storage.NewStorage(db)

	uri, err := url.Parse(opts.MqttURL())
	if err != nil {
		panic(err)
	}

	state := state.NewState()
	state.Connect("cloud.0", uri)

	// Wait mqtt connected
	time.Sleep(3 * time.Second)

	worker := luavm.NewWorker(state)
	go worker.Run()

	r := mux.NewRouter()

	fmt.Println("=========")
	api.Serve(r, store, worker, opts.BaseURL())

	r.HandleFunc(profix+"/{action}/", actionHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(opts.ListenAddr(), r)
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
