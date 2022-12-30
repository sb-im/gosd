# gosd

SuperDock Cloud Server

[![Build Status](https://github.com/sb-im/gosd/workflows/ci/badge.svg)](https://github.com/sb-im/gosd/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/sb-im/gosd/branch/master/graph/badge.svg)](https://codecov.io/gh/sb-im/gosd)
[![GitHub release](https://img.shields.io/github/tag/sb-im/gosd.svg?label=release)](https://github.com/sb-im/gosd/releases)
[![license](https://img.shields.io/github/license/sb-im/gosd.svg?maxAge=2592000)](https://github.com/sb-im/gosd/blob/master/LICENSE)

## Build

* golang >= 1.17.x
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

Variable Name  | Description                                              | Default
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`LISTEN_ADDR`  | Address to listen on (use absolute path for Unix socket) | `0.0.0.0:8000`
`DEBUG`        | Set the value to `1` to enable debug logs                | `false`
`DEMO_MODE`    | Auto Run `database migrate`, `database seed`, `node sync`| `false`
`DATABASE_URL` | Postgresql connection parameters                         | `postgres://postgres:password@localhost/gosd?sslmode=disable`
`STORAGE_URL`  | File storage path                                        | `data/storage`
`BASE_URL`     | Base URL to generate API links and base path             | `http://localhost:8000/gosd/api/v3`
`MQTT_URL`     | MQTT broker Server address                               | `mqtt://admin:public@localhost:1883`
`REDIS_URL`    | Redis Server URL                                         | `redis://localhost:6379/0`
`CLIENT_URL`   | Only `gosd client` use `BASE_URL`                        | `http://localhost:8000/gosd/api/v3`
`LOG_FILE`     | Log File                                                 | `STDOUT`
`LOG_LEVEL`    | Log Level: `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | `info`
`LUA_FILE` | Task lua > `LUA_FILE` > System Default (`luavm/lua/default.lua`) | `default.lua`

