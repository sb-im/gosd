# Lua api

## variable

### `node_id`

### `plan_id`

### `plan_log_id`

## function

### notify

```sblua
notify(id string, rpc table{ "method": string, "params": table})

-> error.string
```

name | type | description
---- | ---- | ----------
<- id   | string | callee id
<- rpc  | table  | `{ "method": string, "params": table}`
-> error| string | `""` Or `"xxxxxxxx"`

About success

```lua
local err = notify("233", {
  ["method"] = "sbim"
})

if err == ""
then
  print("Notify success")
else
  print("Notify failure")
end
```

### rpcCall

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

### asyncRpc

```sblua
async_rpc(id string, rpc table{ "method": string, "params": table}, LChannel)

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
local err = async_rpc_call(node_id, {
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

