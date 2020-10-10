function main(plan)
  print("=== START Lua ===")
  sleep("1ms")

  -- DEBUG Mode
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

  -- Check UPS
  local ups_status = drone:GetMsg("ups_status")
  if ups_status.status ~= "online" then
    plan:ToggleDialog({
      name = "无法执行此任务",
      message = json.encode(ups_status),
      level = "danger",
      buttons = {
        {name = "朕，知道了", message = 'know', level = 'danger'},
      }
    })

    if plan:Gets() ~= "know" then
      print("Task canceled")
      plan:CleanDialog()
      return
    end

    plan:CleanDialog()
    return
  end

  xpcall(function()
    depot:SyncCall("dooropen")
    depot:SyncCall("power_on_remote")
    depot:SyncCall("drone_switch")
    drone:SyncCall("wait_to_boot_finish")
  end, function()
    print(debug.traceback())
    drone:SyncCall("emergency_stop")
  end)

  local battery = drone:GetMsg("battery")
  print(json.encode(data))

  if debug then
    plan:ToggleDialog({
      name = "请继续",
      items = {
        {name = "剩余电量", message = battery.remain .. '%', level = 'info'},
      },
      buttons = {
        {name = "Cancel", message = 'cancel', level = 'primary'},
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

  -- 正片开始！！！ 开始执行任务

  xpcall(function()

    drone:SyncCall("ncp", {"download", "map", plan:FileUrl("file")})
    drone:SyncCall("loadmap")
    drone:SyncCall("check_rtk")

    depot:SyncCall("freedrone")

    drone:SyncCall("startmission_ready")
    drone:SyncCall("startmission")
    drone:SyncCall("check_land")
    drone:SyncCall("mission_preupload_cloud")

    local rfn1 = drone:AsyncCall("mission_upload_nas")

    depot:SyncCall("fixdrone")
    -- depot:SyncCall("power_chargedrone_on")
    depot:SyncCall("doorclose")

    rfn1()
    --depot:SyncCall("power_off_drone")
    depot:SyncCall("drone_switch")


  end,
  function()
    drone:SyncCall("emergency_stop")
  end)

  print("=== END Lua END ===")
  return
end
