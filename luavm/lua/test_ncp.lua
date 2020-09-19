local json = require("json")

function main(node_id)
  local drone = NewNode(node_id)
  drone:SyncCall("ncp", {"download", "map", SD:FileUrl("file")})

  local ok, result = pcall(drone.SyncCall, drone, "ncp", {"download", "map", SD:FileUrl("file")})
  if not ok then
    print("ERROE")
  end
  json.encode(result)

  local ok, GetResult = pcall(drone.AsyncCall, drone, "ncp", {"download", "map", SD:FileUrl("file")})
  if not ok then
    print("ERROE")
  end

  local ok, GetResult2 = pcall(depot.AsyncCall, depot, "ncp", {"download", "map", SD:FileUrl("file")})
  if not ok then
    print("ERROE")
  end


  local ok, data = pcall(GetResult)
  if not ok then
    print("ERROE")
  end
  json.encode(data)

  local ok, data = pcall(GetResult2)
  json.encode(data)

  return node_id
end
