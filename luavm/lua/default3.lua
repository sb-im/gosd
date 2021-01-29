local json = require("json")

function sync_call(id, raw)
  method = tostring(raw)

  local data = rpc_call(id, {
    ["method"] = method,
  })

  if data["result"] then
    print(method .. " success")

    err = SD:IOPuts(method .. " success")
    if err ~= nil then
      print(json.encode(err))
    end

  else
    print("ERROE: " .. method)
    print(json.encode(data))
    return data
  end
  return
end

local sss = 0

function onClose()
  if sync_call(sss, "emergency_stop") then
    return
  end

  print("=== END Lua END ===")
end

function run(node_id)
  print("=== START Lua ===")

  local drone_id = node_id
  sss = drone_id
  local depot_id = get_id("link_id")

  print("Drone Id:", drone_id)
  print("Depot Id:", depot_id)

  print(json.encode(get_status()))

  print(node_id)


  ask_status = {
    name = "确定要执行任务吗？",
    buttons = {
      {name = "不，手滑点错了", message = 'no', level = 'primary'},
      {name = "是的，我要执行", message = 'yes', level = 'danger'},
    }
  }
  SD:ToggleDialog(ask_status)

  msg, err = SD:IOGets()
  if err ~= nil then
    print(msg)
    print(json.encode(err))
  end

  SD:CleanDialog()

  if msg ~= "yes" then
    print("Task canceled")
    onClose()
    return
  end

  if sync_call(depot_id, "power_on_drone") then
    onClose()
    return
  end

  if sync_call(depot_id, "power_on_remote") then
    onClose()
    return
  end

  if sync_call(drone_id, "wait_to_boot_finish") then
    onClose()
    return
  end

  local data, err = get_msg(drone_id, "battery")
  print(json.encode(data))
  --bat = json.encode(data)
  --print(type(bat))
  print(data["remain"])
  print(data["temp"])

  dialog = {
    name = "Checker ~",
    message = "Wow Wow Wow ~",
    level = "success",
    items = {
      {name = "剩余电量", message = data["remain"] .. '%', level = 'info'},
      {name = "电池温度", message = data['temp'] .. '°C', level = 'success'},
      {name = "风速", message = '0 m/s', level = 'danger'},
      {name = "降水", message = '可能有降水', level = 'warning'},
    },
    buttons = {
      {name = "Cancel", message = 'cancel', level = 'primary'},
      {name = "Confirm", message = 'confirm', level = 'danger'},
    }
  }

  err = SD:ToggleDialog(dialog)
  if err ~= nil then
    print(json.encode(err))
  end

  msg, err = SD:IOGets()
  if err ~= nil then
    print(msg)
    print(json.encode(err))
  end

  if msg ~= 'confirm' then

    ch_onDrone = channel.make()
    rpc_async(depot_id, {
      ["method"] = "power_off_drone",
    }, ch_onDrone)

    ch_onDoor = channel.make()
    rpc_async(drone_id, {
      ["method"] = "power_off_remote_smart",
    }, ch_onDoor)


    print("asyncCall send")

    err = 10
    step_1 = 0
    while( step_1 < 2 ) do
      channel.select(
      {"|<-", ch_onDrone, function(ok, data)
        if data["result"] then
          step_1 = step_1 + 1
          print("Drone power_off_remote started")
        else
          print("EEEEEEEEEEE")
          print(json.encode(data))
          step_1 = err
        end
      end},
      {"|<-", ch_onDoor, function(ok, data)
        if data["result"] then
          step_1 = step_1 + 1
          print("Drone power_off_drone started")
        else
          print("EEEEEEEEEEE")
          print(json.encode(data))
          step_1 = err
        end
      end})
    end

    if step_1 >= 10 then

      ask_status = {
        name = "设备故障，请联系厂家！！！",
        message = "Tel: xxxxxxxxxx",
        level = 'danger',
        buttons = {
          {name = "我知道了", message = 'yes', level = 'danger'},
        }
      }
      SD:ToggleDialog(ask_status)

      msg, err = SD:IOGets()
      if err ~= nil then
        print(msg)
        print(json.encode(err))
      end

      SD:CleanDialog()

      return
    end

    return
  end
  SD:CleanDialog()


  if sync_call(depot_id, "dooropen") then
    onClose()
    return
  end

  if sync_call(depot_id, "freedrone") then
    onClose()
    return
  end

  if sync_call(drone_id, "check_drone_ready") then
    onClose()
    return
  end

  local data = rpc_call(drone_id, {
    ["method"] = "ncp",
    ["params"] = {"download", "map", SD:FileUrl("file")},
  })

  if not data["result"] then
    print("EEEEEEEEEEE")
    print(json.encode(data))
    onClose()
    return
  end

  if sync_call(drone_id, "loadmap") then
    onClose()
    return
  end

  if sync_call(drone_id, "check_gps") then
    onClose()
    return
  end

  if sync_call(drone_id, "startmission_ready") then
    onClose()
    return
  end

  if sync_call(drone_id, "startmission") then
    onClose()
    return
  end

  if sync_call(drone_id, "check_land") then
    onClose()
    return
  end

  if sync_call(drone_id, "mission_preupload_cloud") then
    onClose()
    return
  end

  if sync_call(depot_id, "fixdrone") then
    onClose()
    return
  end

  if sync_call(depot_id, "power_chargedrone_on") then
    onClose()
    return
  end

  if sync_call(depot_id, "doorclose") then
    onClose()
    return
  end

  if sync_call(drone_id, "mission_upload_nas") then
    onClose()
    return
  end

  if sync_call(depot_id, "power_off_drone") then
    onClose()
    return
  end

  -- if sync_call(depot_id, "power_off_remote_smart") then
  --   onClose()
  --   return
  -- end

  print("=== END Lua END ===")
  return
end
