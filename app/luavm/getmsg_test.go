package luavm

import (
	"context"
	"fmt"
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

	nodeID := "1"
	task.NodeID = nodeID

	node := &model.Node{
		ID:     nodeID,
		TeamID: task.TeamID,
	}

	orm.Save(node)
	rdb := w.rdb

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := rdb.Set(ctx, fmt.Sprintf(topicNodeSys, task.NodeID, "status"), `{"code":0}`, time.Second).Err(); err != nil {
		t.Error(err)
	}

	if err := rdb.Set(ctx, fmt.Sprintf(topicNodeSys, task.NodeID, "network"), `{"code":0}`, time.Second).Err(); err != nil {
		t.Error(err)
	}

	if err := rdb.Set(ctx, fmt.Sprintf(topicNodeMsg, task.NodeID, "weather"), `{"code":0}`, time.Second).Err(); err != nil {
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
