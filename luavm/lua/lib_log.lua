
--[[
  local log = NewLog(function(line, nu)
    return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
  end)
--]]
function NewLog(fnLine, fnWord)

  if fnLine == nil then
    fnLine = function(line, nu)
      return tostring(nu) .. ": " .. line
    end
  end

  if fnWord == nil then
    fnWord = function(word, nu)
      return tostring(word)
    end
  end

  return {
    _nu = 0,
    _content = "",
    _fnLine = fnLine,
    _fnWord = fnWord,
    Println = function(self, ...)
      self._nu = self._nu + 1
      local line = ""
      for i, v in ipairs(arg) do
        line = line .. self._fnWord(v, i) .. "\t"
      end
      self._content = self._content .. self._fnLine(line, self._nu) .. "\n"
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

function test_logFnLine()
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

function test_logFnWord()
  local log = NewLog(function(line, nu)
    return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
  end,
  function(word, nu)
    return tostring(nu) .. ": " .. word
  end)
  log:Println("xxxxx")
  log:Println("xxxxx", "xxxxx")
  log:Println("xxxxx", "xxxxx", "xxxxx")
  log:Println("xxxxx", "xxxxx", "xxxxx", "xxxxx")
  log:Println(1, "xxxxx")
  print(log:GetContent())
end

test_log()
test_logFnLine()
test_logFnWord()
--]]

