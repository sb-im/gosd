package luavm

import (
	"testing"

	"sb.im/gosd/app/config"
	"sb.im/gosd/app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLuaNode(t *testing.T) {
	task := newTestTask(t)
	cfg := config.Parse()

	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	nodeID := "__luavm_test__node"
	nodeName := "__luavm_test__node_name"
	task.NodeID = nodeID

	node := &model.Node{
		UUID:   nodeID,
		Name:   nodeName,
		TeamID: task.TeamID,
	}

	if err := orm.Save(node).Error; err != nil {
		t.Error(err)
	}

	w := newWorker(t)

	if err := w.doRun(task, []byte(`
function main(task)
  print("### RUN Node RUN ###")

  local node = NewNode(task.nodeID)

  if node.id ~= "`+nodeID+`" then
    error("node id is: " .. node.id)
  end

  if node.name ~= "`+nodeName+`" then
    error("node name is: " .. node.name)
  end

  local isError = false
  pcall(function()
    local node = NewNode("`+nodeID+"not_exist"+`")
    isError = true
  end)

  if isError then
    error("This should is not found node id")
  end

  print("### END Node END ###")
end
`)); err != nil {
		t.Error(err)
	}
}
