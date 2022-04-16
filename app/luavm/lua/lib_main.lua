json = require("json")

sleep = function(time)
  SD:Sleep(time)
end

function SD_main(node_id)
  return main(NewTask(node_id))
end
