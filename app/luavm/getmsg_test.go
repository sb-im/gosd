package luavm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"sb.im/gosd/app/config"
	"sb.im/gosd/app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLuaGetMsg(t *testing.T) {
	task := newTestTask(t)

	w := newWorker(t)

	cfg := config.Parse()
	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	uuid := "__luavm_test__getmsg"
	node := &model.Node{
		UUID:   uuid,
		TeamID: task.TeamID,
	}

	if err := orm.FirstOrCreate(node, "uuid = ?", uuid).Error; err != nil {
		t.Error(err)
	}

	task.NodeID = strconv.Itoa(int(node.ID))

	rdb := w.rdb

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := rdb.Set(ctx, fmt.Sprintf(topicNodeSys, uuid, "status"), `{"code":0}`, time.Second).Err(); err != nil {
		t.Error(err)
	}

	if err := rdb.Set(ctx, fmt.Sprintf(topicNodeSys, uuid, "network"), `{"code":0}`, time.Second).Err(); err != nil {
		t.Error(err)
	}

	if err := rdb.Set(ctx, fmt.Sprintf(topicNodeMsg, uuid, "weather"), `{"code":0}`, time.Second).Err(); err != nil {
		t.Error(err)
	}

	if err := w.doRun(task, []byte(`
function main(task)
  print("### RUN GetMsg RUN ###")

  local node = NewNode(task.nodeID)

  local status = node:GetStatus()
  if status["code"] ~= 0 then
    error(json.encode(status))
  end

  local network = node:GetNetwork()
  if network["code"] ~= 0 then
    error(json.encode(network))
  end

  local weather = node:GetMsg("weather")
  if weather["code"] ~= 0 then
    error(json.encode(weather))
  end

  print("### END GetMsg END ###")
end
`)); err != nil {
		t.Error(err)
	}
}
