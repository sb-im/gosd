local json = require("json")

function main(node_id)
  node_id = "233"
  drone = NewRPC(node_id)
  ok, result = pcall(drone.SyncCall, drone, "ncp", {"download", "map", SD:FileUrl("file")})
  print(ok)
  print(json.encode(result))


  ok, GetResult = pcall(drone.AsyncCall, drone, "ncp", {"download", "map", SD:FileUrl("file")})

  depot = NewRPC("234")
  ok, GetResult2 = pcall(depot.AsyncCall, depot, "ncp", {"download", "map", SD:FileUrl("file")})

  ok, data = pcall(GetResult)
  print("++++++++++++++++++")
  print(ok)
  print(data)
  ok, data = pcall(GetResult2)
  print("++++++++++++++++++")
  print(ok)
  print(data)

  return node_id
end
