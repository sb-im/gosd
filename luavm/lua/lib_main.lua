json = require("json")

sleep = function(time)
  SD:Sleep(time)
end

function SD_main(node_id)
  print("Running")
  local plan = NewPlan(node_id)

  -- Main
  local ret = main(plan)

  -- Record print log
  plan:SetJobFileContent("luavm", "luavm.log", _SD_printResult)
  return ret
end
