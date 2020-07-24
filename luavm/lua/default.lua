local json = require("json")

function run(node_id)
  print("=== START Lua ===")

  local drone_id = node_id
  local depot_id = get_id("link_id")

  print("Drone Id:", drone_id)
  print("Depot Id:", depot_id)

  print(json.encode(get_status()))

  print(node_id)

  --local data, err = get_msg("8", "weather")
  --print(json.encode(data))

  local data = rpc_call(depot_id, {
    ["method"] = "power_chargedrone_on",
  })

  if data["result"] then
    print("power_chargedrone_on success")
  else
    print("EEEEEEEEEEE")
    print(json.encode(data))
    return
  end

  ch_onDrone = channel.make()
  rpc_async(drone_id, {
    ["method"] = "power_on_drone",
  }, ch_onDrone)

  ch_onDoor = channel.make()
  rpc_async(depot_id, {
    ["method"] = "dooropen",
  }, ch_onDoor)


  print("asyncCall send")

  ch_onFree = channel.make()
  ch_onLoad = channel.make()
  local res = {}

  err = 10
  step_1 = 0
  while( step_1 < 2 ) do
    channel.select(
    {"|<-", ch_onDrone, function(ok, data)
      if data["result"] then
        print("Drone Battery started")


        local data = rpc_call(drone_id, {
          ["method"] = "ncp",
          ["params"] = {"download", "map", SD:FileUrl("file")},
        })

        if not data["result"] then
          print("EEEEEEEEEEE")
          print(json.encode(data))
          step_1 = err
          return
        end

        rpc_async(drone_id, {
          ["method"] = "loadmap",
        }, ch_onLoad)
      else
        print("EEEEEEEEEEE")
        print(json.encode(data))
        step_1 = err
      end
    end},
    {"|<-", ch_onLoad, function(ok, data)
      if data["result"] then
        print("Load Map")
        step_1 = step_1 + 1
      else
        print("EEEEEEEEEEE")
        print(json.encode(data))
        step_1 = err
      end
    end},
    {"|<-", ch_onDoor, function(ok, data)
      if data["result"] then
        rpc_async(depot_id, {
          ["method"] = "freedrone",
        }, ch_onFree)
      else
        print("EEEEEEEEEEE")
        print(json.encode(data))
        step_1 = err
      end
    end},
    {"|<-", ch_onFree, function(ok, data)
      if data["result"] then
        local data = rpc_call(depot_id, {
          ["method"] = "check_ready",
        })
        if data["result"] then
          print("check_ready success")
          step_1 = step_1 + 1
        else
          print("EEEEEEEEEEE")
          print(json.encode(data))
          step_1 = err
        end
      else
        print("EEEEEEEEEEE")
        print(json.encode(data))
        step_1 = err
      end
    end})
  end


  if step_1 >= 10 then
    return
  end

  local data = rpc_call(drone_id, {
    ["method"] = "startmission_ready",
  })

  if not data["result"] then
    print("EEEEEEEEEEE")
    print(json.encode(data))
    return
  end


  local data = rpc_call(drone_id, {
    ["method"] = "startmission",
  })

  if not data["result"] then
    print("EEEEEEEEEEE")
    print(json.encode(data))
    return
  end

  print("=== END Lua END ===")
  return
end
