local json = require("json")

function main(plan)
  local node_id = plan.nodeID
  local node = NewNode(node_id)
  print(plan:FileUrl("file"))
  print(plan:LogUrl("file"))


  return node_id
end
