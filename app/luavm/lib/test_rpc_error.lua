function main(plan)
  print("=== START Lua ===")
  sleep("1ms")

  local node = NewNode(plan.nodeID)
  xpcall(function()
    local rfn = node:AsyncCall("__luavm_test__no_result")

    xpcall(function()
      rfn()
    end,
    function()
      print("error is success")
    end)
  end,
  function()
    print(debug.traceback())
  end)
end
