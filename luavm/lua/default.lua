local json = require("json")

function run(node_id)
  print("=== START Lua ===")

  print(get_id("link_id"))
  print(json.encode(get_status()))

  print(node_id)

  local data, err = get_msg("8", "weather")
  print(json.encode(data))

  local err = rpc_notify(node_id, {
    ["method"] = "test",
    ["params"] = {
      ["a"] = "233",
      ["b"] = "456"
    }
  })

  print("=============")
  if err ~= ""
    then
      print(err)
    end
  print("=============")

  ch = channel.make()
  rpc_async(node_id, {
    ["method"] = "test",
    ["params"] = {
      ["a"] = "233",
      ["b"] = "456"
    }
  }, ch)
  print("asyncCall send")

  local res = {}
  channel.select(
  {"|<-", ch, function(ok, data)
    print(ok, data)
    print(json.encode(data))
    if data["result"]
      then
        print("asyncCall success")
      end
  end}
  )
  print(res["result"])

  print(json.encode(rpc_call(node_id, {
    ["method"] = "sync1"
  })))
  print(json.encode(rpc_call(node_id, {
    ["method"] = "sync2"
  })))

  print("=== END Lua END ===")
  return ret['data']
end
