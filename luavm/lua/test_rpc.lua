function main(plan)
  print("=== START Lua ===")
  sleep("1ms")

  local drone_id = plan.nodeID
  local drone = NewNode(drone_id)
  xpcall(function()
    drone:SyncCall("mode_brake")
  end,
  function()
    print(debug.traceback())
  end)

  -- local depot_id = drone:GetID()
  -- local depot = NewNode(depot_id)

  -- local plan = NewPlan()

  -- print("Drone Id:", drone.id)
  -- print("Depot Id:", depot.id)

  -- print("Drone Status:", json.encode(drone:GetStatus()))
  -- local ok, data = pcall(drone.GetMsg, drone, "battery")
  -- if ok then
  --   print(json.encode(data))
  -- end
end
