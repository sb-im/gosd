local json = require("json")

function NewNode(nodeID)
  return {
    id = nodeID,
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
      local raw, err = SD:GetID(self.id, tostring(str))
      if err ~= nil then
        error(err)
      end
      return raw
    end,
  }
end