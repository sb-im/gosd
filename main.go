package main

//go:generate go run generate.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/luavm"
	"sb.im/gosd/state"

	"sb.im/gosd/database"
	"sb.im/gosd/storage"
	"sb.im/gosd/model"

	"github.com/gorilla/mux"
)

var (
	namespace   = "/gosd"
	api_version = "/api/v1"
	profix      = namespace + api_version
)

var accessGrant *AccessGrant

func main() {


	db ,_ := database.NewConnectionPool("postgres://postgres:password@localhost/gosd?sslmode=disable", 1, 10)
	database.Migrate(db)
	s := storage.NewStorage(db)
	p := model.NewPlan()

	p.Name = "TTTTTTTTTTTTTTT"
	p.Description = "DDDD"
	p.Attachments["233"] = "789"
	p.Attachments["23"] = "789"
	s.CreatePlan(p)
	fmt.Println(p)




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

	state := state.NewState()
	mqttClient := state.Connect("cloud.0", uri)
	fmt.Println(mqttClient)

	time.Sleep(1 * time.Second)
	mqttProxy, err := jsonrpc2mqtt.OpenMqttProxy(mqttClient)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mqttProxy)

	//req := []byte(`{"jsonrpc":"2.0","id":"gosd.0","method":"check_ready"}`)
	//ch_recv := make(chan []byte)
	//mqttProxy.AsyncRpc("10", req, ch_recv)
	//fmt.Println("000000000000000")
	//aa := <-ch_recv
	//fmt.Println(string(aa))
	//res, err := mqttProxy.SyncRpc("10", req)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	//res, err := jsonrpc2mqtt.SyncMqttRpc(mq, 10, req)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(res))

	// Wait mqtt connected
	time.Sleep(3 * time.Second)
	path := "test.lua"
	luavm.Run(state, path)

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
