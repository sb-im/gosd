package luavm_test

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"sb.im/gosd/app/config"
	"sb.im/gosd/app/luavm"
	"sb.im/gosd/app/luavm/lib"
	"sb.im/gosd/app/model"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/tests/help"

	"sb.im/gosd/mqttd"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LuaVM Rpc", func() {
	uuid := "__e2e_test__luavm_rpc"
	ctx := context.TODO()
	cfg := config.Parse()

	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	node := &model.Node{
		UUID:   uuid,
		TeamID: 1,
	}
	if err := orm.FirstOrCreate(node, "uuid = ?", uuid).Error; err != nil {
		log.Error(err)
	}

	task := model.Task{
		Name:   "E2E Test",
		NodeID: node.ID,
		TeamID: node.TeamID,
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(redisOpt)
	rdb.ConfigSet(context.Background(), "notify-keyspace-events", "$KEx")

	chI := make(chan mqttd.MqttRpc, 128)
	chO := make(chan mqttd.MqttRpc, 128)

	mqtt := mqttd.NewMqttd(cfg.MqttURL, rdb, chI, chO)
	go mqtt.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chI, chO)
	go rpcServer.Run(ctx)

	ofs := storage.NewStorage(cfg.StorageURL)

	luaFile, err := ioutil.ReadFile(cfg.LuaFilePath)
	if err == nil {
		log.Warn("Use Lua File Path:", cfg.LuaFilePath)
	}
	worker := luavm.NewWorker(luavm.Config{
		Instance: cfg.Instance,
		BaseURL:  cfg.BaseURL,
	}, service.NewService(nil, orm, rdb, ofs), rpcServer, luaFile)
	go worker.Run(ctx)

	go func() {
		if err := help.StartNcp(ctx, os.TempDir(), config.Parse().MqttURL, uuid); err != nil {
			panic(err)
		}
	}()

	// Wait mqttd server startup && sub topic on broker
	time.Sleep(100 * time.Millisecond)

	orm.Create(&task)

	job := model.Job{
		Task: task,
	}
	orm.Create(&job)
	task.Job = &job

	Context("Test Context", func() {
		It("run luaFile", func() {
			luaFile, err := lib.File.ReadFile("test_rpc.lua")
			Expect(err).NotTo(HaveOccurred())

			err = worker.RunTask(&task, luaFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Context Cancel", func() {
		It("run luaFile no result", func() {
			luaFile, err := lib.File.ReadFile("test_rpc_error.lua")
			Expect(err).NotTo(HaveOccurred())

			ch := make(chan error)
			go func() {
				ch <- worker.RunTask(&task, luaFile)
			}()

			time.Sleep(100 * time.Millisecond)
			err = worker.Kill(strconv.Itoa(int(task.ID)))
			Expect(err).NotTo(HaveOccurred())
			Expect(<-ch).NotTo(HaveOccurred())
		})
	})
})
