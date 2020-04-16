package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"io"
	"mime"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var (
	profix = "/api/v1"
)

var accessGrant *AccessGrant

func main() {
	DBlink()
	accessGrant = NewAccessGrant()
	r := mux.NewRouter()
	r.HandleFunc("/oauth/token", oauthHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.HandleFunc(profix+"/plans/", planIndexHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(profix+"/plans/", planCreateHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc(profix+"/plans/{id}", planIDHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.HandleFunc(profix+"/{action}/", actionHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8000", r)
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

func planIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var plans []Plan
	db.Where("").Find(&plans)

	b, err := json.Marshal(plans)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}

func planIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	s, err := strconv.Atoi(vars["id"])
	if err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	plan := &Plan{}
	plan.Find(s)

	b, err := json.Marshal(plan)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write(b)
}

func planCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}
	//bb := make([]byte, 4096)
	//r.Body.Read(bb)
	//fmt.Println(string(bb))

	rparams := make(map[string]string)

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}

			if p.FileName() != "" {
				fmt.Println(p.FileName())
				filepath := "data/storage/" + p.FileName()
				err := ioutil.WriteFile(filepath, slurp, 0644)
				if err != nil {
					log.Fatal(err)
				}
				rparams[p.FormName()] = filepath
			} else {
				rparams[p.FormName()] = string(slurp)
			}

		}
	}

	//fmt.Println(rparams)
	plan := &Plan{
		Name:           rparams["name"],
		Description:    rparams["description"],
		File:           rparams["file"],
		Node_id:        rparams["node_id"],
		Cycle_types_id: rparams["cycle_types_id"],
	}

	plan.Create()

	b, err := json.Marshal(plan)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
