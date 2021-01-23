package cli

import (
	"net/http"
	"net/url"
	"runtime"
	"time"

	"sb.im/gosd/api"
	"sb.im/gosd/config"
	"sb.im/gosd/luavm"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	"miniflux.app/logger"

	"github.com/gorilla/mux"
)

func startDaemon(store *storage.Storage, opts *config.Options) {
	logger.Info("Starting gosd...")

	//go showProcessStatistics()

	uri, err := url.Parse(opts.MqttURL())
	if err != nil {
		panic(err)
	}

	state := state.NewState()
	state.Connect(opts.MqttClientID(), uri)

	// Wait mqtt connected
	time.Sleep(3 * time.Second)

	worker := luavm.NewWorker(state, store)
	go worker.Run()

	r := mux.NewRouter()

	logger.Info("=========")
	api.Serve(r, store, worker, opts.BaseURL())
	http.ListenAndServe(opts.ListenAddr(), r)

	logger.Info("Process gracefully stopped")
}

func showProcessStatistics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		logger.Debug("Sys=%vK, InUse=%vK, HeapInUse=%vK, StackSys=%vK, StackInUse=%vK, GoRoutines=%d, NumCPU=%d",
			m.Sys/1024, (m.Sys-m.HeapReleased)/1024, m.HeapInuse/1024, m.StackSys/1024, m.StackInuse/1024,
			runtime.NumGoroutine(), runtime.NumCPU())
		time.Sleep(30 * time.Second)
	}
}
