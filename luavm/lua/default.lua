local json = require("json")

function main(node_id)
  print("=== START Lua ===")

  local drone_id = node_id
  local drone = NewNode(drone_id)

  local depot_id = drone:GetID()
  local depot = NewNode(depot_id)

  local plan = NewPlan()

  print("Drone Id:", drone.id)
  print("Depot Id:", depot.id)

  print("Drone Status:", json.encode(drone:GetStatus()))
  local ok, data = pcall(drone.GetMsg, drone, "battery")
  if ok then
    print(json.encode(data))
  end

  plan:CleanDialog()
  ask_status = {
    name = "确定要执行任务吗？",
    buttons = {
      {name = "不，手滑点错了", message = 'no', level = 'primary'},
      {name = "是的，我要执行", message = 'yes', level = 'danger'},
    }
  }
  plan:ToggleDialog(ask_status)

  if plan:Gets() ~= "yes" then
    print("Task canceled")
    return
  end
  plan:CleanDialog()

  -- drone:SyncCall("test")

  xpcall(function()
    depot:SyncCall("power_on_drone")
    depot:SyncCall("power_on_remote")
    drone:SyncCall("wait_to_boot_finish")
  end, function()
    drone:SyncCall("emergency_stop")
  end)

  local data = drone:GetMsg("battery")
  print(json.encode(data))
  print(data.vol_cell)
  print(data["vol_cell"])

  dialog = {
    name = "Checker ~",
    message = "Wow Wow Wow ~",
    level = "success",
    items = {
      {name = "剩余电量", message = data.remain .. '%', level = 'info'},
      {name = "电池温度", message = data.temp .. '°C', level = 'success'},
      {name = "风速", message = '0 m/s', level = 'danger'},
      {name = "降水", message = '可能有降水', level = 'warning'},
    },
    buttons = {
      {name = "Cancel", message = 'cancel', level = 'primary'},
      {name = "Confirm", message = 'confirm', level = 'danger'},
    }
  }

  plan:ToggleDialog(dialog)

  if plan:Gets() ~= 'confirm' then
    print("cancel")
    plan:CleanDialog()

    xpcall(function()
      local rfn1 = depot:AsyncCall("power_off_drone")
      local rfn2 = depot:AsyncCall("power_off_remote_smart")

      rfn1()
      rfn2()
    end,
    function()
      drone:SyncCall("emergency_stop")

    end)
    return
  end
  plan:CleanDialog()

  -- 正片开始！！！ 开始执行任务

  xpcall(function()
    depot:SyncCall("dooropen")
    depot:SyncCall("freedrone")

    drone:SyncCall("check_drone_ready")
    drone:SyncCall("ncp", {"download", "map", SD:FileUrl("file")})
    drone:SyncCall("loadmap")
    drone:SyncCall("check_gps")
    drone:SyncCall("startmission_ready")
    drone:SyncCall("startmission")
    drone:SyncCall("check_land")
    drone:SyncCall("mission_preupload_cloud")

    local rfn1 = drone:AsyncCall("mission_upload_nas")

    depot:SyncCall("fixdrone")
    depot:SyncCall("power_chargedrone_on")
    depot:SyncCall("doorclose")

    rfn1()
    depot:SyncCall("power_off_drone")


  end,
  function()
    drone:SyncCall("emergency_stop")
  end)

  print("=== END Lua END ===")
  return
end
