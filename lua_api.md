# Lua api

## variable

### `node_id`

### `plan_id`

### `plan_log_id`

## function

### `rpc_notify`

```sblua
rpc_notify(id string, rpc table{ "method": string, "params": table})

-> error.string
```

name | type | description
---- | ---- | ----------
<- id   | string | callee id
<- rpc  | table  | `{ "method": string, "params": table}`
-> error| string | `""` Or `"xxxxxxxx"`

About success

```lua
local err = rpc_notify("233", {
  ["method"] = "sbim"
})

if err == ""
then
  print("Notify success")
else
  print("Notify failure")
end
```

### `rpc_call`

```sblua
rpc_call(id string, rpc table{ "method": string, "params": table})

-> table{ "result": table{}, "error": "" }, error.string
```

name | type | description
---- | ---- | ----------
<- id   | string | callee id
<- rpc  | table  | `{ "method": string, "params": table}`
-> res  | table  | `{ "result": table{}, "error": table{} }`
-> error| string | `""` Or `"xxxxxxxx"`

```lua
local res, err = rpc_call("233", {
  ["method"] = "sbim"
})

if err == ""
then
  print("rpc send success")
else
  print("rpc send failure")
end


if res["result"]
then
  print("rpc call success")
else
  print("rpc call failure")
end
```

### `rpc_async`

```sblua
rpc_async(id string, rpc table{ "method": string, "params": table}, LChannel)

-> error.string
```

name | type | description
---- | ---- | ----------
<- id   | string | callee id
<- rpc  | table  | `{ "method": string, "params": table}`
<- ch   | LChannel | look like go channel
-> error| string | `""` Or `"xxxxxxxx"`

```lua
local res, err = rpc_call("233", {
  ["method"] = "sbim"
})


ch = channel.make()
local err = rpc_async(node_id, {
  ["method"] = "test",
  ["params"] = {
    ["a"] = "233",
    ["b"] = "456"
  }
}, ch)

if err == ""
then
  print("rpc send success")
else
  print("rpc send failure")
end





local res = {}

-- Block
channel.select(
{"|<-", ch, function(ok, data)
  print(ok, data)
  print(json.encode(data))

  res = data
end}
)

if res["result"]
then
  print("rpc call success")
else
  print("rpc call failure")
end
```

### Service

#### `get_status`

> Get Status

function:

```golang
get_status() (data tables{}, error string)
```

params:

name | type | description
---- | ---- | ----------
-> data | table  | `table{}`
-> error| string | `""` Or `"xxxxxxxx"`

example:

```lua
local data, err = get_status()

if err == ""
then
  print("get success")
else
  print("get failure")
end

-- Print data
print(json.encode(data))
```

#### `get_msg`

> Get Message

function:

```golang
get_msg(id, msg string) (data tables{}, error string)
```

params:

name | type | description
---- | ---- | ----------
<- id   | string | callee id
<- msg  | string | `weather`, `battery` ...
-> data | table  | `table{}`
-> error| string | `""` Or `"xxxxxxxx"`

example:

```lua
local data, err = get_msg("8", "weather")

if err == ""
then
  print("get success")
else
  print("get failure")
end

-- Print data
print(json.encode(data))
```

#### `get_id`

> Get various types of ID

function:

```golang
get_id(str string) (id string)
```

params:

name | type | description
---- | ---- | ----------
<- str  | string | Only support `link_id`
-> id   | string | id if `id == "0"` is No id

example:

```lua
local id = get_id("link_id")

if id ~= "0"
then
  print("get id" + id)
else
  print("No find id")
end
```

