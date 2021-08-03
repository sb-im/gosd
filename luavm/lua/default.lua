
_raw_print = print
local log = NewLog(function(line, nu)
  return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
end,
function(word, nu)
  if type(word) == "table" then
    return json.encode(word)
  end
  return tostring(word)
end)
print = function(...)
  _raw_print(unpack(arg))
  log:Println(unpack(arg))
end

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

  plan:CleanDialog()
  local ask_status = {
    name = "错误",
    message = "请联系管理员在后端指定默认流程，或者上传自定义流程",
    level = "error",
    buttons = {
      {name = "返回", message = 'no', level = 'primary'},
    }
  }
  plan:ToggleDialog(ask_status)

  if plan:Gets() == "no" then
    print("Task canceled")
    return
  end
  plan:CleanDialog()

  -- drone:SyncCall("test")

  print("=== END Lua END ===")
  plan:SetJobFileContent("luavm", "luavm.txt", log:GetContent())
  return
end
