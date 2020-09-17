local json = require("json")

function run(node_id)
  print("node_id --> ", node_id)
  drone = NewRPC(node_id)
  ok, result = pcall(drone.SyncCall, drone, "ncp", {"download", "map", SD:FileUrl("file")})
  if ok then
    print(result)
    print(json.encode(result))
  end
  return node_id
end
