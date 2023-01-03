function main(task)
  print("=== RUN LUA RUN ===")
  sleep("1ms")

  print("Node Status:", json.encode(NewNode(task.nodeID):GetStatus()))

  task:CleanDialog()
  local ask_status = {
    name = "Are You OK?",
    message = "Need set default lua workflow",
    level = "error",
    buttons = {
      {name = "Back", message = 'no', level = 'primary'},
    }
  }
  task:ToggleDialog(ask_status)

  if task:Gets() == "no" then
    print("Task canceled")
    return
  end
  task:CleanDialog()

  print("=== END LUA END ===")
end
