function run(node_id)
  print("=== START Lua ===")

  local drone_id = node_id

  local data = rpc_call(drone_id, {
    ["method"] = "ncp",
    ["params"] = {"download", "map", SD:FileUrl("file")}
  })

  if data["result"] then
    print("download success")
  else
    print("EEEEEEEEEEE")
    print(json.encode(data))
    return
  end

  print("=== END Lua END ===")

  return node_id
end
