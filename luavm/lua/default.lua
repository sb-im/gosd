function main(plan)
  print("=== START Lua ===")
  sleep("1ms")
  local debug = true

  local drone_id = plan.nodeID
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
    print(debug.traceback())
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
    drone:SyncCall("ncp", {"download", "map", plan:FileUrl("file")})
    drone:SyncCall("loadmap")
    drone:SyncCall("check_gps")

  end,
  function()
    drone:SyncCall("emergency_stop")
  end)

  local depotGPS = depot:GetStatus().status
  print("Depot GPS:", json.encode(depotGPS))

  local droneGPS = drone:GetMsg("position")
  print("Drone GPS:", json.encode(droneGPS))
  print("Lat Distance :", (tonumber(droneGPS.lat) - tonumber(depotGPS.lat)))
  print("Lat Distance :", (tonumber(droneGPS.lng) - tonumber(depotGPS.lng)))

  local distance = 10000000000 * (
    (tonumber(droneGPS.lat) - tonumber(depotGPS.lat)) *
    (tonumber(droneGPS.lat) - tonumber(depotGPS.lat)) +
    (tonumber(droneGPS.lng) - tonumber(depotGPS.lng)) *
    (tonumber(droneGPS.lng) - tonumber(depotGPS.lng)) )

  if debug then
    plan:ToggleDialog({
      name = "最后确认",
      message = "Wow Wow Wow ~",
      level = "success",
      items = {
        {name = "Drone Lat", message = tostring(droneGPS.lat), level = 'success'},
        {name = "Drone Lng", message = tostring(droneGPS.lng), level = 'success'},
        {name = "Depot Lat", message = tostring(depotGPS.lat), level = 'success'},
        {name = "Depot Lng", message = tostring(depotGPS.lng), level = 'success'},
        {name = "Distance ", message = tostring(distance), level = 'warning'},
      },
      buttons = {
        {name = "Cancel" , message = 'cancel', level = 'primary'},
        {name = "Confirm", message = 'confirm', level = 'danger'},
      }
    })
  end

  if plan:Gets() ~= 'confirm' then
    print("cancel")
    plan:CleanDialog()
    drone:SyncCall("emergency_stop")

    return
  end
  plan:CleanDialog()

  xpcall(function()

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
