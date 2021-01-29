function main(plan)
  print("=== START Lua ===")
  sleep("1ms")

  local drone_id = plan.nodeID
  local drone = NewNode(drone_id)
  local plan = NewPlan()

  print("Drone Id:", drone.id)

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

  -- 正片开始！！！ 开始执行任务

  xpcall(function()
    drone:SyncCall("ncp", {"download", "map", plan:FileUrl("file")})
    drone:SyncCall("loadmap")
    drone:SyncCall("startmission_ready")
    drone:SyncCall("startmission")
  end,
  function()
    print(debug.traceback())
    --drone:SyncCall("emergency_stop")
  end)

  print("=== END Lua END ===")
  return
end
