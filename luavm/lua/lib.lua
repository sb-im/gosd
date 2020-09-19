local json = require("json")

function NewRPC(nodeID)
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
  }
end
