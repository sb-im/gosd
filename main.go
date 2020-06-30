package main

//go:generate go run generate.go

import (
	"fmt"
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
	http.ListenAndServe(opts.ListenAddr(), r)
}
