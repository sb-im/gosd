package cmd

import (
	"context"
	"net/http"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/luavm"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"
	"sb.im/gosd/rpc2mqtt"

	"sb.im/gosd/mqttd"
	"sb.im/gosd/state"

	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Daemon() {
	log.Warn("Launch gosd V3")

	cfg := config.DefaultConfig()
	if err := env.Parse(cfg); err != nil {
		log.Errorf("%+v\n", err)
	}

	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(redisOpt)

	store := state.NewState(cfg.RedisURL)

	chI := make(chan mqttd.MqttRpc, 128)
	chO := make(chan mqttd.MqttRpc, 128)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqtt := mqttd.NewMqttd(cfg.MqttURL, store, chI, chO)
	go mqtt.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chI, chO)
	go rpcServer.Run(ctx)

	ofs := storage.NewStorage(cfg.StorageURL)
	worker := luavm.NewWorker(orm, rdb, ofs, rpcServer, []byte{})
	go worker.Run()

	srv := service.NewService(orm, rdb, worker)
	srv.StartSchedule()
	log.Warn("=== RUN ===")

	api := v3.NewApi(orm, srv)
	http.Handle("/gosd/api/v3/", api)
	http.ListenAndServe(":8000", nil)
}
