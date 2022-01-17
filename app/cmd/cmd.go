package cmd

// TODO: Current tmp Launcher

import (
	"net/http"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/luavm"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

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

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	ofs := storage.NewStorage("data")
	worker := luavm.NewWorker(orm, rdb, ofs, []byte{})
	go worker.Run()

	srv := service.NewService(orm, rdb, worker)
	srv.StartSchedule()
	log.Warn("=== RUN ===")

	api := v3.NewApi(orm, srv)
	http.Handle("/gosd/api/v3/", api)
	http.ListenAndServe(":8000", nil)
}
