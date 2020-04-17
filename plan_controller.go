package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

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

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	}

	plan := &Plan{}
	plan.Find(id)

	b, err := json.Marshal(plan)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write(b)
}

func file2path(r *http.Request) map[string]string {
	params := make(map[string]string)

	mediaType, mimeParams, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, mimeParams["boundary"])
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
				params[p.FormName()] = filepath
			} else {
				params[p.FormName()] = string(slurp)
			}

		}
	}

	return params
}

func planCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	params := file2path(r)

	node_id, err := strconv.Atoi(params["node_id"])
	if err != nil {
		fmt.Println(err)
	}

	cycle_types_id, err := strconv.Atoi(params["cycle_types_id"])
	if err != nil {
		fmt.Println(err)
	}

	plan := &Plan{
		Name:           params["name"],
		Description:    params["description"],
		File:           params["file"],
		Node_id:        node_id,
		Cycle_types_id: cycle_types_id,
	}

	plan.Create()

	b, err := json.Marshal(plan)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func planUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	}

	plan := &Plan{}
	plan.Find(id)

	params := file2path(r)

	for _, key := range []string{"name", "description", "file"} {
		db.Model(&plan).Update(key, params[key])
	}

	for _, key := range []string{"node_id", "cycle_types_id"} {
		id, err := strconv.Atoi(params[key])
		if err != nil {
			fmt.Println(err)
		}
		db.Model(&plan).Update(key, id)
	}

	b, err := json.Marshal(plan)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func planDestroyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	}

	plan := &Plan{}
	db.Delete(plan, id)

	b, err := json.Marshal(plan)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write(b)
}
