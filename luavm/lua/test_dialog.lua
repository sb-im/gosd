local json = require("json")

function run(node_id)
  SD:CleanDialog()

  err = SD:IOPuts("checking")
  if err ~= nil then
    print(json.encode(err))
  end

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

  err = SD:ToggleDialog(dialog)
  if err ~= nil then
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
    name = "ARE YOU OK ?",
    buttons = {
      {name = "Fine, thank you.", message = 'fine', level = 'primary'},
      {name = "I feel bad.", message = 'bad', level = 'danger'},
    }
  }
  SD:ToggleDialog(ask_status)

  msg, err = SD:IOGets()
  if err ~= nil then
    print(msg)
    print(json.encode(err))
  end

  SD:CleanDialog()

  err = SD:IOPuts("checked")
  if err ~= nil then
    print(json.encode(err))
  end

  return node_id
end
