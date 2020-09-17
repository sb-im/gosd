local json = require("json")

function NewRPC(nodeID)
  return {
    id = nodeID,
    SyncCall = function(self, method, params)
      data, err = SD2:SyncCall(self.id, tostring(method), json.encode(params))
      if err ~= nil then
        error(err)
      end

      result = json.decode(data)
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
