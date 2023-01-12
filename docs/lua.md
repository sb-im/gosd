# Lua api

> Lua version: 5.1

## Main

For Example:

```lua
function main(task)
end
```

The function `main()` params is `task`

## Task

This is Task Object, Only `main()` params

```lua
task:CleanDialog()

ask_status = {
  name = "ARE YOU OK ?",
  buttons = {
    {name = "Fine, thank you.", message = 'fine', level = 'primary'},
    {name = "I feel bad.", message = 'bad', level = 'danger'},
  }
}

task:ToggleDialog(ask_status)

msg, err = task:Gets()
if err ~= nil then
  print(msg)
  print(json.encode(err))
end

task:CleanDialog()

print(task:FileUrl("file"))
```

### Name

Task Name

```lua
print(task:name)
```

### ToggleDialog(table)

Topic: `tasks/:id/dialog`

Send a confirmation form

[params table reference](https://developer.sb.im/#/mqtt?id=dialog)

```lua
ask_status = {
  name = "ARE YOU OK ?",
  buttons = {
    {name = "Fine, thank you.", message = 'fine', level = 'primary'},
    {name = "I feel bad.", message = 'bad', level = 'danger'},
  }
}
task:ToggleDialog(ask_status)
```

### CleanDialog()

Close Dialog

Topic: `tasks/:id/dialog` message `{}`

```lua
task:CleanDialog()
```

### Gets() string

get `tasks/:id/term` message

```lua
print(task:Gets())

msg, err = task:Gets()
if err ~= nil then
  print(msg)
  print(json.encode(err))
end
```

### Puts(string)

put `tasks/:id/term` message

```lua
task:Puts("test")
sleep("1s")
task:Puts("test2")
```

### Notification(msg string, level = 5)

```lua
task:Notification("notification")

task:Notification("notification", 3)
```

### nodeID

> The id of this current task need `node`

### GetExtra() extra map[string]string

Get Plan Extra `key/value`

This `extra` is `map[string]string` of `table`

### SetExtra(extra map[string]string)

Set Plan Extra `key/value`

```lua
print("Extra:", json.encode(task:GetExtra()))

local extra = task:GetExtra()
extra["xxx"] = "aa"
extra["ccc"] = "aaa"
task:SetExtra(extra)
```

### GetJobExtra() extra map[string]string

Get Plan Job Extra `key/value`

### SetJobExtra(extra map[string]string)

Set Plan Job Extra `key/value`

```lua
local job_extra = task:GetJobExtra()
job_extra["ttt"] = "xxx"
print("Job Extra:", json.encode(job_extra))
task:SetJobExtra(job_extra)
```

### GetFileContent(key string) (filename, content string)

Get Plan Files Content

### SetFileContent(key, filename, content string)

Set Plan Files Content

```lua
print("Files:", json.encode(task:GetFiles()))

xpcall(function()
  print(task:GetFileContent("test_files"))
end,
function()
  task:SetFileContent("test_files", "test.txt", "233")
end)

local filename, content = task:GetFileContent("test_files")
if content == "233" then
  task:SetFileContent("test_files", "test2.txt", "456")
else
  task:SetFileContent("test_files", "test.txt", "233")
end
print(task:GetFileContent("test_files"))
```

### GetJobFileContent(key string) (filename, content string)

Get Plan Job Files Content

### SetJonFileContent(key, filename, content string)

Set Plan Job Files Content

```lua
print("Job Files:", json.encode(task:GetJobFiles()))

xpcall(function()
  print(task:GetJobFileContent("test_files"))
end,
function()
  task:SetJobFileContent("test_files", "test.txt", "233")
end)

local filename, content = task:GetJobFileContent("test_files")
if content == "233" then
  task:SetJobFileContent("test_files", "test2.txt", "456")
else
  task:SetJobFileContent("test_files", "test.txt", "233")
end
print(task:GetJobFileContent("test_files"))
```

### FileUrl(key string)

```lua
print("Blobs:", task:FileUrl("test_files"))
print("Blobs:", task:FileUrl("test_blobs"))
```

### JobFileUrl(key string)

```lua
print("Job Blobs:", task:JobFileUrl("test_files"))
print("Job Blobs:", task:JobFileUrl("test_blobs"))
```

## Class node

```lua
local drone_id = task.nodeID
local drone = NewNode(drone_id)

local depot_id = drone:GetID()
local depot = NewNode(depot_id)

xpcall(function()
  local promise = depot:AsyncCall("wait_to_boot_finish")
  depot:SyncCall("power_on_drone")
  depot:SyncCall("power_on_remote")

  -- Block: get rpc result
  local result = promise()

end, function()
  drone:SyncCall("emergency_stop")
end)

local data = drone:GetMsg("battery")

```

### SyncCall(string [, table]) table

> Sync jsonrpc call

### AsyncCall(string [, table]) function() table

> Async jsonrpc call

return a function

### GetID([string]) string

default params: `link_id`

`GetID() == GetID("link_id")`

### GetStatus() table

`nodes/:id/status`

### GetMsg(string) table

`nodes/:id/msg/+`

## Log

### NewLog([fnLine, fnWord function])

```lua
local log = NewLog()

-- Or

local log = NewLog(function(line, nu)
  return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
end)

-- Or

local log = NewLog(function(line, nu)
  return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
end,
function(word, nu)
  if type(word) == "table" then
    return json.encode(word)
  end
  return tostring(word)
end)
```

### fnLine(string, number) string

> Every line hook function

Default: if `fnLine == nil`

```lua
if fnLine == nil then
  fnLine = function(line, nu)
    return tostring(nu) .. ": " .. line
  end
end
```

### fnWord(any, number) string

> Every word hook function

Default: if `fnWord == nil`

```lua
if fnWord == nil then
  fnWord = function(word, nu)
    return tostring(word)
  end
end
```

### Println(...)

> Print auto add `\n`

```lua
log:Println("xxxxx")
log:Println("xxxxx", "xxxxx", "xxxxx", "xxxxx")
log:Println(1, "xxxxx")
print(log:GetContent())
```

```lua
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
```

### GetContent() string

> Get All Content

```lua
_raw_print = print
local log = NewLog(function(line, nu)
  return tostring(nu) .. ": \t" .. os.date("%Y/%m/%d %H:%M:%S") .. " " .. line
end,
function(word, nu)
  if type(word) == "table" then
    return json.encode(word)
  end
  return tostring(word)
end)
print = function(...)
  _raw_print(unpack(arg))
  log:Println(unpack(arg))
end

function main(task)
  print("=== START Lua ===")

  -- Your workflow

  print("=== END Lua END ===")
  task:SetJobFileContent("luavm", "luavm.txt", log:GetContent())
end
```

## Geo

Geography util. a Object

### Distance(aLng, aLat, bLng, bLat number) number

2D Distance measurement

`Geo:Distance(aLng, aLat, bLng, bLat)`

```lua
-- 114.2247765, 22.6857991
-- 114.22475167, 22.68580217
-- = 2.57202994
local distance = Geo:Distance(114.2247765, 22.6857991, 114.22475167, 22.68580217)
if math.floor(distance) == 2 then
  print("Distance:", distance)
else
  error("Distance sum error:", distance)
end
```

## sleep(string)

```lua
sleep("1ms")
```

such as "300ms", "-1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h" [golang time#ParseDuration](https://golang.org/pkg/time/#ParseDuration)

