
STOP = false

--合并两个table
function MergeTables(...)
  local tabs = {...}
  if not tabs then
    return {}
  end
  local origin = tabs[1]
  for i = 2,#tabs do
    if origin then
      if tabs[i] then
        for k,v in pairs(tabs[i]) do
          table.insert(origin,v)
        end
      end
    else
      origin = tabs[i]
    end
  end
  return origin
end

function Running(rfn, plan, drone)
  xpcall(function()
    rfn()
  end,
  function()
    -- Need Wait Stop reset
    sleep("1s")

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
      sleep("1s")
    end

    -- Recursion
    return Running(rfn, plan, drone)
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

  -- 提前获取天气
  local depot_weather = depot:AsyncCall("get_weather")

  plan:CleanDialog()
  local ask_status = {
    name = "[1准备开始]=>2环境确认=>3航线确认=>↑起飞↑",
    message = "确定要开始执行任务吗？",
    level = "info",
    buttons = {
      {name = "不，返回", message = 'no', level = 'primary'},
      {name = "是的，下一步", message = 'yes', level = 'danger'},
    }
  }
  plan:ToggleDialog(ask_status)

  if plan:Gets() ~= "yes" then
    print("Task canceled")
    return
  end
  plan:CleanDialog()

  -- drone:SyncCall("test")

  -- 给机场设置保护状态，避免中途用户操作
  depot:SyncCall("set_protectstatus_true")
  depot:SyncCall("power_off_fan")

  local drone_battery
  local depot_weather_result
  xpcall(function()
    depot:SyncCall("power_on_drone_and_remote")
    sleep("1s")
    depot_weather_result = depot_weather()
    drone_battery = depot:SyncCall("get_drone_battery")
    print(drone_battery)
    print(depot_weather_result)
  end, function()
    print(debug.traceback())
    depot:SyncCall("set_protectstatus_false")
    drone:SyncCall("emergency_stop")
  end)

  local dialog_environment_check = {
    name = "1准备开始=>[2环境确认]=>3航线确认=>↑起飞↑",
    message = "请确认剩余电量、风速、降雨概率等情况",
    level = "info",
    items = {},
    buttons = {
      {name = "取消任务", message = 'cancel', level = 'primary'},
      {name = "下一步", message = 'confirm', level = 'danger'},
    }
  }

  local drone_predict_battery_level = 'success'
  if tonumber(drone_battery) < 30 then
    drone_predict_battery_level = 'danger'
  elseif tonumber(drone_battery) < 60 then
    drone_predict_battery_level = 'warning'
  end

  print("prepare depot_weather")
  dialog_environment_check.items = depot_weather_result
  table.insert(dialog_environment_check.items, 1, {name = "预计剩余电量", message = drone_battery .. '%', level = drone_predict_battery_level})
  print(dialog_environment_check.items)

  plan:ToggleDialog(dialog_environment_check)

  if plan:Gets() ~= 'confirm' then
    print("cancel")
    plan:CleanDialog()

    xpcall(function()
      depot:SyncCall("power_off_drone_and_remote")
      depot:SyncCall("doorclose")
      depot:SyncCall("set_protectstatus_false")
    end,
    function()
      depot:SyncCall("set_protectstatus_false")
      drone:SyncCall("emergency_stop")

    end)
    return
  end
  plan:CleanDialog()

  -- 正片开始！！！ 开始执行任务

  local check_waypoints
  local check_drone
  xpcall(function()
    local download_map = drone:AsyncCall("ncp", {"download", "map", plan:FileUrl("file")})
    depot:SyncCall("dooropen")
    drone:SyncCall("wait_to_boot_finish")
    -- drone:SyncCall("ncp", {"download", "map", plan:FileUrl("file")})
    download_map()
    drone:SyncCall("loadmap")
    check_drone = drone:SyncCall("check_drone")
    print(check_drone)
    check_waypoints = depot:SyncCall("check_waypoints")
    print(check_waypoints)
    drone:SyncCall("check_gps")

  end,
  function()
    depot:SyncCall("set_protectstatus_false")
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
        message = "当前GPS信号欠佳，建议继续等待后再起飞",
        level = "warning",
        items = {
          {name = "飞机与机场距离", message = string.format("%.2f",  distance) .. 'M', level = 'danger'},
        },
        buttons = {
          {name = "取消任务" , message = 'cancel', level = 'primary'},
          {name = "继续等待", message = 'confirm', level = 'danger'},
        }
      })

      if plan:Gets() ~= 'confirm' then
        print("cancel")
        plan:CleanDialog()
        depot:SyncCall("set_protectstatus_false")
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

    if distance < 6 then
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

  local battery_level = 'success'
  if tonumber(data.remain) < 30 or tonumber(data.temp) < 5 then
    battery_level = 'danger'
  elseif tonumber(data.remain) < 60 or tonumber(data.temp) < 10 then
    battery_level = 'warning'
  end

  local dialog_last_check_level = 'success'
  if battery_level == 'danger' then
    dialog_last_check_level = 'danger'
  elseif battery_level == 'warning' or distanceLevel == 'warning' then
    dialog_last_check_level = 'warning'
  end

  local dialog_last_check = {
    name = "1准备开始=>2环境确认=>[3航线确认]=>↑起飞↑",
    message = "请核对剩余电量、航线是否正常",
    level = dialog_last_check_level,
    items = {},
    buttons = {
      {name = "取消任务" , message = 'cancel', level = 'primary'},
      {name = "↑开始起飞↑", message = 'confirm', level = 'danger'},
    }
  }

  dialog_last_check.items = MergeTables(check_waypoints, check_drone)
  table.insert(dialog_last_check.items, 1, {name = "电池状态", message = '电量 ' .. data.remain .. '% ,温度 ' .. data.temp .. '°C' , level = battery_level})
  -- table.insert(dialog_last_check.items, 1, {name = "距离", message = string.format("%.2f",  distance) .. 'M', level = distanceLevel})
  local final_check_result = drone:SyncCall("final_check")
  if final_check_result ~= 'ok' then
    table.insert(dialog_last_check.items, 1, {name = "自检错误", message = final_check_result , level = 'danger'})
    dialog_last_check.buttons = {
      {name = "取消任务" , message = 'cancel', level = 'primary'},
    }
    dialog_last_check.message = "自检错误，无法起飞！"
    dialog_last_check.level = "danger"
  end

  if debug then
    plan:ToggleDialog(dialog_last_check)
  end

  if plan:Gets() ~= 'confirm' then
    print("cancel")
    plan:CleanDialog()
    xpcall(function()
      depot:SyncCall("power_off_drone_and_remote")
      depot:SyncCall("doorclose")
      depot:SyncCall("set_protectstatus_false")
    end,
    function()
      depot:SyncCall("set_protectstatus_false")
      drone:SyncCall("emergency_stop")

    end)

    return
  end
  plan:CleanDialog()

  xpcall(function()
    drone:SyncCall("startmission_ready")
    depot:AsyncCall("freedrone")
    drone:SyncCall("startmission")
  end,
  function()
    depot:SyncCall("set_protectstatus_false")
    drone:SyncCall("emergency_stop")
  end)

  -- mission running
  Running(drone:AsyncCall("check_land"), plan, drone)
  if STOP then
    return
  end

  xpcall(function()
    depot:SyncCall("power_off_fan")
    drone:SyncCall("mission_preupload_cloud")
    local rfn1 = drone:AsyncCall("mission_upload_nas")
    drone:SyncCall("ncp", {"upload", "view_url", plan:LogFileUrl("查看原始数据")})
    drone:SyncCall("ncp", {"upload", "view_login_url", plan:LogFileUrl("(登录)")})
    depot:SyncCall("fixdrone")
    depot:SyncCall("power_chargedrone_on")
    depot:SyncCall("doorclose")
    -- 复位机场的保护状态
    depot:SyncCall("set_protectstatus_false")

    rfn1()
    depot:SyncCall("power_off_drone_and_remote")


  end,
  function()
    depot:SyncCall("set_protectstatus_false")
    drone:SyncCall("emergency_stop")
  end)

  print("=== END Lua END ===")
  return
end
