local json = require("json")

function run(node_id)
  dialog = {
    name = "Checker ~",
    message = "Wow Wow Wow ~",
    level = "success",
    items = {
      {name = "nn", message = 'mm', level = 'info'},
      {name = "n2", message = 'ok', level = 'success'},
      {name = "n3", message = 'Not ok', level = 'danger'},
      {name = "n4", message = '...', level = 'warning'},
    },
    buttons = {
      {name = "Cancel", message = 'cancel', level = 'primary'},
      {name = "Confirm", message = 'confirm', level = 'danger'},
    }
  }

  if SD:ToggleDialog(dialog) ~= nil then
    print(json.encode(err))
  end

  msg, err = SD:IOGets()
  if err ~= nil then
    print(msg)
    print(json.encode(err))
  end

  if msg ~= 'confirm' then
    return
  end


  ask_status = {
    name = "Input ~",
    inputs = {
      {name = "Height", message = '10', level = 'primary'},
      {name = "Speed", message = '2', level = 'danger'},
    }
  }
  SD:ToggleDialog(ask_status)

  return node_id
end
