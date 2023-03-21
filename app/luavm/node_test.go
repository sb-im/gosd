package luavm

import (
	"context"
	"testing"

	"sb.im/gosd/app/model"
)

func TestLuaNode(t *testing.T) {
	w := newWorker(t)
	task := helpTestNewTask(t, "Unit Test Lua Node", w)

	uuid := "__luavm_test__node"
	nodeName := "__luavm_test__node_name"

	node := &model.Node{
		UUID:   uuid,
		Name:   nodeName,
		TeamID: task.TeamID,
	}

	if err := w.orm.FirstOrCreate(node, "uuid = ?", uuid).Error; err != nil {
		t.Error(err)
	}

	task.NodeID = node.ID

	if err := w.doRun(context.Background(), task, []byte(`
function main(task)
  print("### RUN Node RUN ###")

  local node = NewNode(task.nodeID)

  if node.id ~= "`+uuid+`" then
    error("node id is: " .. node.id)
  end

  if node.name ~= "`+nodeName+`" then
    error("node name is: " .. node.name)
  end

  local isError = false
  pcall(function()
    local node = NewNode("`+uuid+"not_exist"+`")
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
