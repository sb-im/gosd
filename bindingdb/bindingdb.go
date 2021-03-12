package bindingdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"sb.im/gosd/luavm"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	redis "github.com/gomodule/redigo/redis"

	logger "log"
)

func BindingDB(ctx context.Context, s *state.State, db *storage.Storage) {
	keyspace := "__keyspace@0__:%s"
	psc := redis.PubSubConn{Conn: s.Pool.Get()}
	psc.PSubscribe(fmt.Sprintf(keyspace, "plans/*"))
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			logger.Printf("%s: message: %s\n", v.Channel, v.Data)

			if strings.Split(v.Channel, "/")[2] == "running" {
				continue
			}

			topic := strings.Split(v.Channel, ":")[1]
			logger.Println(s.StringGet(topic))

			raw, err := s.BytesGet(topic)
			if err != nil {
				logger.Println(err)
			}

			if string(raw) == "{}" {
				continue
			}

			task := &luavm.Task{}
			if err := json.Unmarshal(raw, task); err != nil {
				logger.Println(err)
			}

			plan, err := db.PlanByID(task.PlanID)
			if err != nil {
				logger.Println(err)
			}
			plan.Attachments = task.Files
			plan.Extra = task.Extra

			db.UpdatePlan(plan)

		case redis.Subscription:
			logger.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			logger.Println(v)
			//return v
		default:
			logger.Println("default")
		}
	}
}
