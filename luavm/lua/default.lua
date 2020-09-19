local json = require("json")

function main(node_id)
  print("=== START Lua ===")

  local drone_id = node_id
  local drone = NewNode(drone_id)

  local depot_id = drone:GetID()
  local depot = NewNode(depot_id)

  print("Drone Id:", drone.id)
  print("Depot Id:", depot.id)

  print("Drone Status:", json.encode(drone:GetStatus()))
  local ok, data = pcall(drone.GetMsg, drone, "battery")
  if ok then
    print(json.encode(data))
  end

  SD:CleanDialog()
  ask_status = {
    name = "确定要执行任务吗？",
    buttons = {
      {name = "不，手滑点错了", message = 'no', level = 'primary'},
      {name = "是的，我要执行", message = 'yes', level = 'danger'},
    }
  }
  SD:ToggleDialog(ask_status)

  msg, err = SD:IOGets()
  if err ~= nil then
    print(msg)
    print(json.encode(err))
  end

  SD:CleanDialog()

  if msg ~= "yes" then
    print("Task canceled")
    return
  end

  print("=== END Lua END ===")
  return
end
