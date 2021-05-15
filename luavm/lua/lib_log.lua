
--[[
  local log = NewLog(function(line, nu)
    return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
  end)
--]]
function NewLog(fn)
  if fn == nil then
    fn = function(line, nu)
      return tostring(nu) .. ": " .. line
    end
  end
  return {
    _nu = 0,
    _call = fn,
    _content = "",
    Println = function(self, ...)
      self._nu = self._nu + 1
      local line = ""
      for i, v in ipairs(arg) do
        line = line .. tostring(v) .. "\t"
      end
      self._content = self._content .. self._call(line, self._nu) .. "\n"
    end,
    GetContent = function(self)
      return self._content
    end
  }
end

--[[
function example_save(plan)
  local log = NewLog()
  log:Println("xxxxx")
  log:Println("xxxxx", "xxxxx")
  plan:SetJobFileContent("luavm", "luavm.txt", log:GetContent())
end
--]]

--[[
function test_log()
  local log = NewLog()
  log:Println("xxxxx")
  log:Println("xxxxx")
  log:Println("xxxxx", "xxxxx")
  log:Println("xxxxx", "xxxxx", "xxxxx")
  log:Println("xxxxx", "xxxxx", "xxxxx", "xxxxx")
  log:Println(1, "xxxxx")
  print(log:GetContent())
end

function test_logfn()
  local log = NewLog(function(line, nu)
    return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
  end)
  log:Println("xxxxx")
  log:Println(1, "xxxxx")
  for i=10,1,-1 do
    log:Println(i, "ccccc")
  end
  print(log:GetContent())
end

test_log()
test_logfn()
--]]

