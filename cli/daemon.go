package cli

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"sb.im/gosd/api"
	"sb.im/gosd/config"
	"sb.im/gosd/luavm"
	"sb.im/gosd/mqttd"
	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/state"
	"sb.im/gosd/model"
	"sb.im/gosd/storage"

	logger "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func StartDaemon(store *storage.Storage, opts *config.Options) {
	logger.Info("Starting gosd...")

	if opts.LogFile() != "STDOUT" {
		file, err := os.Create(opts.LogFile())
		if err != nil {
			panic(err)
		}
		logger.SetOutput(file)
	}
	lvl, err := logger.ParseLevel(opts.LogLevel())
	if err != nil {
		panic(err)
	}
	logger.SetLevel(lvl)
	logger.Warn("log level: ", logger.GetLevel().String())

	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	//dsn := "host=localhost user=postgres password=password dbname=gosd port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := opts.DatabaseURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&model.Schedule{})

	//go showProcessStatistics()

	state := state.NewState(opts.RedisURL())

	chI := make(chan mqttd.MqttRpc, 128)
	chO := make(chan mqttd.MqttRpc, 128)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqttd := mqttd.NewMqttd(opts.MqttURL(), state, chI, chO)
	go mqttd.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chI, chO)
	go rpcServer.Run(ctx)

	luaStr := []byte{}
	if luaFile, err := os.Open(opts.LuaFile()); err == nil {
		if luaStr, err = ioutil.ReadAll(luaFile); err == nil {
			logger.Warn("lua file: ", opts.LuaFile())
		}
	}

	worker := luavm.NewWorker(state, store, rpcServer, luaStr)
	go worker.Run()

	go func() {
		r := mux.NewRouter()
		logger.Info("=========")
		api.Serve(r, db, state, store, worker, api.ServeConfig{
			BaseURL: opts.BaseURL(),
			OauthID: opts.OauthID(),
			OauthSecret: opts.OauthSecret(),
		})

		logger.Warn(opts.ListenAddr())
		http.ListenAndServe(opts.ListenAddr(), r)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	// This Worker Need Wait data sync
	worker.Close()
	time.Sleep(1 * time.Second)

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
