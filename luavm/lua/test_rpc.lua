function main(plan)
  print("=== START Lua ===")
  sleep("1ms")

  local drone_id = plan.nodeID
  local drone = NewNode(drone_id)
  xpcall(function()
    drone:SyncCall("get")
  end,
  function()
    print(debug.traceback())
  end)
end
