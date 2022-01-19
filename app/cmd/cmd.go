package cmd

// TODO: Current tmp Launcher

import (
	"context"
	"net/http"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/luavm"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"
	"sb.im/gosd/rpc2mqtt"

	"sb.im/gosd/mqttd"
	"sb.im/gosd/state"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Execute() {
	log.Warn("Launch gosd V3")

	dsn := "host=localhost user=postgres password=password dbname=gosd port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	orm.AutoMigrate(&model.Team{})
	orm.AutoMigrate(&model.User{})
	orm.AutoMigrate(&model.Session{})
	orm.AutoMigrate(&model.UserTeam{})

	orm.AutoMigrate(&model.Schedule{})
	orm.AutoMigrate(&model.Task{})
	orm.AutoMigrate(&model.Blob{})
	orm.AutoMigrate(&model.Job{})

	redisURL := "redis://localhost:6379/1"
	redisOpt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(redisOpt)

	store := state.NewState(redisURL)

	mqttURL := "mqtt://admin:public@localhost:1883"

	chI := make(chan mqttd.MqttRpc, 128)
	chO := make(chan mqttd.MqttRpc, 128)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqtt := mqttd.NewMqttd(mqttURL, store, chI, chO)
	go mqtt.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chI, chO)
	go rpcServer.Run(ctx)

	ofs := storage.NewStorage("data")
	worker := luavm.NewWorker(orm, rdb, ofs, rpcServer, []byte{})
	go worker.Run()

	srv := service.NewService(orm, rdb, worker)
	srv.StartSchedule()
	log.Warn("=== RUN ===")

	api := v3.NewApi(orm, srv)
	http.Handle("/gosd/api/v3/", api)
	http.ListenAndServe(":8000", nil)
}
