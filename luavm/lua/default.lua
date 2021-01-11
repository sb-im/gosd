
STOP = false

function running(plan, drone)
  xpcall(function()
    drone:SyncCall("check_land")
  end,
  function()
    plan:ToggleDialog({
      name = "终止了任务",
      message = "后续动作该如何执行？",
      level = "danger",
      buttons = {
        {name = "悬停", message = 'mode_brake', level = 'danger'},
        {name = "返航", message = 'mode_rtl', level = 'warning'},
        {name = "取消后续动作", message = 'cancel', level = 'primary'},
      }
    })

    -- Need Wait Stop reset
    sleep("1ms")

    local input = plan:Gets()
    print("Gets:", input)
    plan:CleanDialog()

    if input == 'cancel' then
      STOP = true
      return
    elseif input == 'mode_brake' then
      xpcall(function()
        drone:SyncCall("mode_brake")
      end,
      function()
        print(debug.traceback())
      end)
    elseif input == 'mode_rtl' then
      xpcall(function()
        drone:SyncCall("mode_rtl")
      end,
      function()
        print(debug.traceback())
      end)
    else
      sleep("1ms")
    end
    return running(plan, drone)
  end)
end

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
  -- local ok, data = pcall(drone.GetMsg, drone, "battery")
  -- if ok then
  --   print(json.encode(data))
  -- end

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

  local drone_battery
  local depot_weather
  xpcall(function()
    depot:SyncCall("power_on_drone_and_remote")
    sleep("1s")
    drone_battery = depot:SyncCall("get_drone_battery")
    depot_weather = depot:SyncCall("get_weather")
    print(drone_battery)
    print(depot_weather)
  end, function()
    print(debug.traceback())
    drone:SyncCall("emergency_stop")
  end)

  dialog = {
    name = "起飞环境确认",
    message = "确认剩余电量、风速、降雨概率等情况 ~",
    level = "success",
    items = {
      {name = "剩余电量", message = drone_battery .. '%', level = 'info'},
      -- {name = "电池温度", message = data.temp .. '°C', level = 'success'},
      -- {name = "风速", message = '0 m/s', level = 'danger'},
      -- {name = "降水", message = '可能有降水', level = 'warning'},
    },
    buttons = {
      {name = "Cancel", message = 'cancel', level = 'primary'},
      {name = "Confirm", message = 'confirm', level = 'danger'},
    }
  }

  print("prepare depot_weather")
  dialog.items = depot_weather
  table.insert(dialog.items, 1, {name = "剩余电量", message = drone_battery .. '%', level = 'info'})
  print(dialog.items)

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
    drone:SyncCall("wait_to_boot_finish")
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

  local retrycount = 0
  local distance
  while true do
    if retrycount > 3 then

      plan:ToggleDialog({
        name = "是否继续检查 GPS",
        message = "Wow Wow Wow ~",
        level = "danger",
        items = {
          {name = "Distance ", message = string.format("%.2f",  distance) .. 'M', level = 'danger'},
        },
        buttons = {
          {name = "Cancel" , message = 'cancel', level = 'primary'},
          {name = "Confirm", message = 'confirm', level = 'danger'},
        }
      })

      if plan:Gets() ~= 'confirm' then
        print("cancel")
        plan:CleanDialog()
        drone:SyncCall("emergency_stop")
        print("=== Distance:", distance, "END ===")

        return
      else
        plan:CleanDialog()
        retrycount = 0
      end
    end
    retrycount = retrycount + 1

    local droneGPS = drone:GetMsg("position")
    print("Drone GPS:", json.encode(droneGPS))

    distance = GetDistance(
      tonumber(droneGPS.lng),
      tonumber(droneGPS.lat),
      tonumber(depotGPS.lng),
      tonumber(depotGPS.lat)
    )

    if distance < 5 then
      print("Distance:", distance, "continue")
      break
    end

    sleep("5s")
  end

  local distanceLevel = 'warning'
  if distance < 3 then
    distanceLevel = 'success'
  end

  local data = drone:GetMsg("battery")
  print(json.encode(data))
  print(data.vol_cell)
  print(data["vol_cell"])

  if debug then
    plan:ToggleDialog({
      name = "最后确认",
      message = "注意起飞点、电池温度是否正常 ~",
      level = distanceLevel,
      items = {
        {name = "距离", message = string.format("%.2f",  distance) .. 'M', level = distanceLevel},
        {name = "剩余电量", message = data.remain .. '%', level = 'info'},
        {name = "电池温度", message = data.temp .. '°C', level = 'info'},
      },
      buttons = {
        {name = "取消" , message = 'cancel', level = 'primary'},
        {name = "↑确认起飞↑", message = 'confirm', level = 'danger'},
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
    depot:AsyncCall("freedrone")
    drone:SyncCall("startmission")
  end,
  function()
    drone:SyncCall("emergency_stop")
  end)

  -- mission running
  running(plan, drone)
  if STOP then
    return
  end

  xpcall(function()
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
