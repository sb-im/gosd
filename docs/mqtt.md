# MQTT Topic

## Task

### Running

Topic: `tasks/:id/running`

If This task not running:

```json
{}
```

If This task running:

```json
{
  "files":{},
  "extra":{},
  "job":{
    "job_id":1,
    "files":{},
    "extra":{}
  }
}
```

Field           | Type   | Description
--------------- | ------ | -----------
`files`         | map[string]string | `key/blob_id`
`extra`         | map[string]string | `key/value`
`job`           | Object            | This task run log
`job.job_id`    | int64             | `task_id` + `job_id` == `job.id`
`job.files`     | map[string]string | Job `key/blob_id`
`job.extra`     | map[string]string | Job `key/value`

### Term

Topic: `tasks/:id/term`

> Any data, **NOTE: not json, may be plain**

### ~~Notification~~ (Discard)

Topic: `tasks/:id/notification`

```json
{"time":1565413755,"level":1,"msg":""}
```

Field | Type   | Description
----  | ------ | -----------
time  | uint64 | unix timestamp length `10`
level | uint   | 0-7 `levelEnum`
msg   | string | message body

levelEnum: 0-7

level | Name
----- | ----
0     | Emergency
1     | Alert
2     | Critical
3     | Error
4     | Warn
5     | Notice
6     | Info
7     | Debug

Reference: [RFC5424](https://tools.ietf.org/html/rfc5424#section-6.2.1)

### Dialog

Topic: `tasks/:id/dialog`

Field           | Type   | Description
--------------- | ------ | -----------
name            | string | Opt: A title
message         | string | Opt: A descriptive message
level           | string | Opt: `levelEnum`
items           | array  | Opt: Items
items.name      | string | Name
items.message   | string | Opt: A descriptive message
items.level     | string | Opt: `levelEnum`
buttons         | array  | Opt: Button Group
buttons.name    | string | Button name
buttons.message | string | Send message payload
buttons.level   | string | Opt: `levelEnum`

**levelEnum** : `primary`, `success`, `info`, `warning`, `danger`

Example 0: Clean Dialog

```json
{}
```

Example 1: Check Form

```json
{
  "name": "Checker ~",
  "message": "Not recommended",
  "level": "warning",
  "items": [
    {"name": "Check", "message": "check check", "level": "unkonw"},
    {"name": "Drone", "message": "ok", "level": "primary"},
    {"name": "Depot", "message": "ok", "level": "success"},
    {"name": "Weather: Wind speed", "message": "Strong wind", "level": "warning"},
    {"name": "Weather: Rain forecast", "message": "It is raining", "level": "danger"}
  ],
  "buttons": [
    {"name": "Cancel", "message": "cancel", "level": "primary"},
    {"name": "Confirm","message": "confirm", "level": "danger"}
  ]
}
```

Example 2: Ask Status

```json
{
  "name": "ARE YOU OK ?",
  "buttons": [
    {"name": "Fine, thank you.", "message": "fine", "level": "primary"},
    {"name": "I feel bad.", "message": "bad", "level": "danger"}
  ]
}
```

