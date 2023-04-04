# gosd

SuperDock Cloud Server

[![Build Status](https://github.com/sb-im/gosd/workflows/ci/badge.svg)](https://github.com/sb-im/gosd/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/sb-im/gosd/branch/master/graph/badge.svg)](https://codecov.io/gh/sb-im/gosd)
[![GitHub release](https://img.shields.io/github/tag/sb-im/gosd.svg?label=release)](https://github.com/sb-im/gosd/releases)
[![license](https://img.shields.io/github/license/sb-im/gosd.svg?maxAge=2592000)](https://github.com/sb-im/gosd/blob/master/LICENSE)

## Build

* golang >= 1.18.x
* make

```bash
make
```

### Swagger

```bash
# Docs swagger
go get -u github.com/swaggo/swag/cmd/swag
make swagger
```

## Dependencies

* PostgreSQL >= 13
* Redis
* Mqtt broker (Mosquitto Or Emqx)
  * Mosquitto >= 1.6

## Introduction

`gosd` is a `gosd server` and `gosd client`

* server (up gosd daemon)
  * database
* client (http client)
  * node
  * user
  * team

### Run as Demo

This Demo Mode auto inserts the demo data

```bash
DEMO_MODE=true ./gosd server
```

### Run as Production

```bash
cp dot.env .env

# Create database tables
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd database migrate

# Create database init seed
# Default:
# - TeamId: 1
# - UserId: 1
# - SessId: 1
./gosd database seed

# Create User
# gosd users add <username> <password>
./gosd users add demo demodemo

# Sync node data
./gosd node sync ./data

# Run
MQTT_URL=mqtt://admin:public@localhost:1883 \
REDIS_URL=redis://localhost:6379/0 \
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd server
```

## Environment Variables

Read Environment Variables, Support `.env` file

```bash
cp dot.env .env
```

### Auth

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`SECRET`       | JWT Secret                                               | `falling-cats-and-dogs`
`API_KEY`      | Admin Api. http header `X-Api-Key`, unset is disable     | `the-elephant-in-the-room`
`BASIC_AUTH`   | A auth plugin. support http base auth                    | `ture`

### Service URL

Dependent services

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`MQTT_URL`     | MQTT broker Server address                               | `mqtt://admin:public@localhost:1883`
`REDIS_URL`    | Redis Server URL                                         | `redis://localhost:6379/0`
`STORAGE_URL`  | File storage path                                        | `data/storage`
`DATABASE_URL` | Postgresql connection parameters                         | `postgres://postgres:password@localhost/gosd?sslmode=disable`

### Public URL

Provide services, example: `nginx` gateway need change `BASE_URL`

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`BASE_URL`     | Base URL to generate API links and base path             | `http://localhost:8000/gosd/api/v3`
`CLIENT_URL`   | Only `gosd client` use `BASE_URL`                        | `http://localhost:8000/gosd/api/v3`
`API_MQTT_WS`  | MQTT broker Websocket server address                     | `ws://admin:public@localhost:1883`
`LISTEN_ADDR`  | Address to listen on (use absolute path for Unix socket) | `0.0.0.0:8000`

### Feature Flags

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`SCHEDULE`     | Only Single Node. **Not support cluster**                | `true`
`EMQX_AUTH`    | Use Emqx redis auth plugin. If Mosquitto, Set `false`    | `false`
`LUA_FILE`     | Task lua > `LUA_FILE` > Default(`luavm/lua/default.lua`) | `default.lua`

### Custom Config

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`INSTANCE`     | Instance name                                            | `gosd`
`LANGUAGE`     | Language                                                 | `en_US`
`TIMEZONE`     | Timezone                                                 | `Asia/Shanghai`
`LOG_LEVEL`| `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`  | `info`

### Developer Flags

This System Default User have a `user.Id == 1`, `team.Id == 1` and `sess.Id == 1`.
The `SINGLE_USER` flag enable System Default User

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`DEBUG`        | Set the value to `true` to enable debug logs             | `false`
`DEMO_MODE`    | Auto Run `database migrate`, `database seed`, `node sync`| `false`
`SINGLE_USER`  | System Only One User. All belong user.Id == `1`          | `false`

