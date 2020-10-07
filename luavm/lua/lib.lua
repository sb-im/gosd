json = require("json")

function NewNode(nodeID)
  return {
    id = tostring(nodeID),
    AsyncCall = function(self, method, params)
      rpc = {
        id = SD:GenRpcID(),
        method = tostring(method),
        params = params,
        jsonrpc = "2.0",
      }

      print(self.id, "SEND -->", json.encode(rpc))
      local rpcID, err = SD:RpcSend(self.id, json.encode(rpc))
      if err ~= nil then
        error(err)
      end

      return function()
        print(self.id, "await -->", rpcID)
        data, err = SD:RpcRecv(rpcID)
        if err ~= nil then
          error(err)
        end
        print(self.id, "RECV -->", data)
        return json.decode(data)
      end
    end,
    SyncCall = function(self, method, params)
      local result = self.AsyncCall(self, method, params)()
      if result["result"] then
        return result["result"]
      end

      if result["error"] then
        error(result["error"])
      end

      error("Result Error")
    end,
    GetMsg = function(self, msg)
      local raw, err = SD:GetMsg(self.id, msg)
      if err ~= nil then
        error(err)
      end
      return json.decode(raw)
    end,
    GetStatus = function(self)
      local raw, err = SD:GetStatus(self.id)
      if err ~= nil then
        error(err)
      end
      return json.decode(raw)
    end,
    GetID = function(self, str)
      if str == nil then
        str = "link_id"
      end
      return (self.GetStatus(self)).status[str]
    end,
  }
end

function NewPlan(nodeID)
  return {
    nodeID = nodeID,
    ToggleDialog = function(self, dialog)
      local err = SD:ToggleDialog(dialog)
      if err ~= nil then
        error(err)
      end
    end,
    CleanDialog = function(self)
      SD:CleanDialog()
    end,
    Gets = function(self)
      local data, err = SD:IOGets()
      if err ~= nil then
        error(err)
      end
      return data
    end,
    Puts = function(self, data)
      local err = SD:IOPuts(data)
      if err ~= nil then
        error(err)
      end
    end,
    FileUrl = function(self, key)
      local data = SD:FileUrl(key)
      if data == "" then
        error("key not found")
      end
      return data
    end
  }
end

sleep = function(time)
  SD:Sleep(time)
end

function SD_main(node_id)
  return main(NewPlan(node_id))
end
