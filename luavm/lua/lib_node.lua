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
        local result = json.decode(data)
        if result["result"] then
          return result["result"]
        end

        if result["error"] then
          error(result["error"])
        end

        error("Result Error")
      end
    end,
    SyncCall = function(self, method, params)
      return self.AsyncCall(self, method, params)()
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
    GetNetwork = function(self)
      local raw, err = SD:GetNetwork(self.id)
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
