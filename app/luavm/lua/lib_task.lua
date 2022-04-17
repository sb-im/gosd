function NewTask(nodeID)
  local task = SD:GetTask()
  return {
    name = task.name,
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
    Notification = function(self, msg, level)
      return SD:Notification({
        time = os.time(),
        level = tonumber(level or 5),
        msg = tostring(msg),
      })
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
      return self.GetAttach(self).extra or {}
    end,
    SetExtra = function(self, extra)
      local data = self.GetAttach(self)
      data.extra = extra
      self.SetAttach(self, data)
    end,
    GetJobExtra = function(self)
      return self.GetAttach(self).job.extra or {}
    end,
    SetJobExtra = function(self, extra)
      local data = self.GetAttach(self)
      data.job.extra = extra
      self.SetAttach(self, data)
    end,
    GetFiles = function(self)
      return self.GetAttach(self).files or {}
    end,
    SetFiles = function(self, files)
      local data = self.GetAttach(self)
      data.files = files
      self.SetAttach(self, data)
    end,
    GetJobFiles = function(self)
      return self.GetAttach(self).job.files or {}
    end,
    SetJobFiles = function(self, files)
      local data = self.GetAttach(self)
      data.job.files = files
      self.SetAttach(self, data)
    end,
    GetFileContent = function(self, key)
      local id = self.GetFiles(self)[key]
      return SD:BlobReader(id)
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
      local blobID = self.GetFiles(self)[key]
      if not blobID then
        blobID = tostring(SD:BlobCreate("filename", ""))
        local files = self:GetFiles(self)
        files[key] = blobID
        self.SetFiles(self, files)
      end
      return SD:BlobUrl(blobID)
    end,
    JobFileUrl = function(self, key)
      local blobID = self.GetJobFiles(self)[key]
      if not blobID then
        blobID = tostring(SD:BlobCreate("filename", ""))
        local files = self:GetJobFiles(self)
        files[key] = blobID
        self.SetJobFiles(self, files)
      end
      return SD:BlobUrl(blobID)
    end,
    LogFileUrl = function(self, key)
      return self.JobFileUrl(self, key)
    end
  }
end
