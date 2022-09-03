package cmd

import (
	"context"
	"io/ioutil"
	"net/http"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/logger"
	"sb.im/gosd/app/luavm"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"
	"sb.im/gosd/app/store"
	"sb.im/gosd/rpc2mqtt"

	"sb.im/gosd/mqttd"
	"sb.im/gosd/state"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewHandler(ctx context.Context, cfg *config.Config) http.Handler {
	log.Debugf("%+v\n", cfg)

	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.NewGorm(),
	})
	if err != nil {
		panic(err)
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(redisOpt)

	ofs := storage.NewStorage(cfg.StorageURL)
	s := store.NewStore(cfg, orm, rdb, ofs)

	store := state.NewState(cfg.RedisURL)

	chI := make(chan mqttd.MqttRpc, 128)
	chO := make(chan mqttd.MqttRpc, 128)

	mqtt := mqttd.NewMqttd(cfg.MqttURL, store, chI, chO)
	go mqtt.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chI, chO)
	go rpcServer.Run(ctx)

	luaFile, err := ioutil.ReadFile(cfg.LuaFilePath)
	if err == nil {
		log.Warn("Use Lua File Path:", cfg.LuaFilePath)
	}
	worker := luavm.NewWorker(luavm.DefaultConfig(), s, rpcServer, luaFile)
	go worker.Run(ctx)

	srv := service.NewService(orm, rdb, worker)
	srv.StartSchedule()

	return api.NewApi(s, srv)
}

func Daemon() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Warn("Launch gosd V3")
	handler := NewHandler(ctx, config.Parse())
	log.Warn("=== RUN ===")
	http.ListenAndServe(":8000", handler)
}
