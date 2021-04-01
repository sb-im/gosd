package api

import (
	//"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	"miniflux.app/http/response/json"
)

func (h *handler) mqttPut(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.Path, "/mqtt/")[1]

	// TODO: need to more than 4096
	raw := make([]byte, 4096)
	defer r.Body.Close()
	n, _ := r.Body.Read(raw)
	log.Println("SET: ", key, string(raw))
	data, err := redis.Bytes(h.cache.Do("GET", key))
	if err != nil {
		json.ServerError(w, r, err)
	}
	h.cache.Do("SET", key, raw[:n])
	w.Header().Set("Content-Type", `application/json`)
	w.Write(data)
}

func (h *handler) mqttGet(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.Path, "/mqtt/")[1]
	log.Println("GET: ", key)
	data, err := redis.Bytes(h.cache.Do("GET", key))
	if err != nil {
		json.ServerError(w, r, err)
	}
	w.Header().Set("Content-Type", `application/json`)
	w.Write(data)
}
