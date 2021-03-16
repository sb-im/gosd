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
    GetAttach = function(self)
      local raw, err = SD:GetAttach()
      if err ~= nil then
        error(err)
      end
      return json.decode(raw)
    end,
    SetAttach = function(self, attach)
      local err = SD:SetAttach(json.encode(attach))
      if err ~= nil then
        error(err)
      end
    end,
    GetExtra = function(self)
      return self.GetAttach(self).extra
    end,
    SetExtra = function(self, extra)
      local data = self.GetAttach(self)
      data.extra = extra
      self.SetAttach(self, data)
    end,
    GetJobExtra = function(self)
      return self.GetAttach(self).job.extra
    end,
    SetJobExtra = function(self, extra)
      local data = self.GetAttach(self)
      data.job.extra = extra
      self.SetAttach(self, data)
    end,
    GetFiles = function(self)
      return self.GetAttach(self).files
    end,
    SetFiles = function(self, files)
      local data = self.GetAttach(self)
      data.files = files
      self.SetAttach(self, data)
    end,
    GetJobFiles = function(self)
      return self.GetAttach(self).job.files
    end,
    SetJobFiles = function(self, files)
      local data = self.GetAttach(self)
      data.job.files = files
      self.SetAttach(self, data)
    end,
    GetFileContent = function(self, key)
      local id = self.GetFiles(self)[key]
      return SD:BlobReader(tonumber(id))
    end,
    SetFileContent = function(self, key, filename, content)
      local id = tonumber(self.GetFiles(self)[key])
      if filename == "" then
        if id then
          -- TODO:
          -- SD:BlobDelete(id)
        end
        return
      end

      if id then
        -- Update Blob
        SD:BlobUpdate(id, filename, content)
      else
        -- Create Blob
        local id = SD:BlobCreate(filename, content)
        local files = self:GetFiles(self)
        files[key] = tostring(id)
        self.SetFiles(self, files)
      end
    end,
    GetJobFileContent = function(self, key)
      local id = self.GetJobFiles(self)[key]
      return SD:BlobReader(tonumber(id))
    end,
    SetJobFileContent = function(self, key, filename, content)
      local id = tonumber(self.GetJobFiles(self)[key])
      if filename == "" then
        if id then
          -- TODO:
          -- SD:BlobDelete(id)
        end
        return
      end

      if id then
        -- Update Blob
        SD:BlobUpdate(id, filename, content)
      else
        -- Create Blob
        local id = SD:BlobCreate(filename, content)
        local files = self:GetJobFiles(self)
        files[key] = tostring(id)
        self.SetJobFiles(self, files)
      end
    end,
    FileUrl = function(self, key)
      local data = SD:FileUrl(key)
      if data == "" then
        error("key not found")
      end
      return data
    end,
    LogFileUrl = function(self, key)
      return SD:LogFileUrl(key)
    end
  }
end

sleep = function(time)
  SD:Sleep(time)
end

--[[
6371004*ACOS(
        (
                SIN(RADIANS(C2))*SIN(RADIANS(F2))
                +
                COS(RADIANS(C2))*COS(RADIANS(F2))
                *
                COS(RADIANS(E2-B2))
        )
)
--]]

function GetDistance(aLng, aLat, bLng, bLat)
  -- Earth Radius: 6371004
  return 6371004 * math.acos(
      math.sin(math.rad(aLat)) * math.sin(math.rad(bLat))
      +
      math.cos(math.rad(aLat)) * math.cos(math.rad(bLat))
      *
      math.cos(math.rad(bLng - aLng))
    )
end

function SD_main(node_id)
  return main(NewPlan(node_id))
end
