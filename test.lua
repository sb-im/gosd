local json = require("json")

function info(n)
  local param = {
    ['id'] = n,
    ['name'] = 'jyjiiiiii'
  }
  ret = call_service(filepoolservice,"getuserinfo",param)
  print("22222222222222222222222")
  print(plan_id)
  print(plan_log_id)
  print(node_id)

  ch = channel.make()
  async_rpc_call("10", json.encode({
    ["method"] = "test",
    ["params"] = {
      ["a"] = "233",
      ["b"] = "456"
    }
  }), ch)
  print("5555555555555555555555")

  local res = {}
  channel.select(
  {"|<-", ch, function(ok, data)
    print(ok, data)
    if json.decode(data)["result"]
      then
        res = json.decode(data)
        print("SSSSSSSSSSSSSSssss")
      end
  end}
  )
  print(res["result"])

  print(rpc_call("10", json.encode({
    ["method"] = "sync1"
  })))
  print(rpc_call("10", json.encode({
    ["method"] = "sync2"
  })))

  print("33333333333333333333333")
  return ret['data']
end
