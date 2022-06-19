# gosd

SuperDock Cloud Server

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

## Run

```sh
# Create database tables
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd database migrate

# Create database init seed
# Default:
# - TeamId: 1
# - UserId: 1
# - SessId: 1
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd database seed

# Create User
# gosd users add <username> <password>
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd users add demo demodemo

# Run
MQTT_URL=mqtt://admin:public@localhost:1883 \
REDIS_URL=redis://localhost:6379/0 \
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd
```

## Environment Variables

Variable Name  | Description                                              | Default Value
-------------- | -------------------------------------------------------- | -------------------------------------------------------------
`DEBUG`        | Set the value to `1` to enable debug logs                | Off
`DATABASE_URL` | Postgresql connection parameters                         | `postgres://postgres:password@localhost/gosd?sslmode=disable`
`LISTEN_ADDR`  | Address to listen on (use absolute path for Unix socket) | `127.0.0.1:8000`
`PORT`         | Override `LISTEN_ADDR` to `0.0.0.0:$PORT` (PaaS)         | None
`BASE_URL`     | Base URL to generate API links and base path             | `http://localhost:8000/gosd/api/v3`
`MQTT_URL`     | MQTT broker Server address                               | `mqtt://admin:public@localhost:1883`
`REDIS_URL`    | Redis Server URL                                         | `redis://localhost:6379/0`
`LOG_FILE`     | Log File                                                 | `STDOUT`
`LOG_LEVEL`    | Log Level: `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | `info`
`LUA_FILE` | Task lua > `LUA_FILE` > System Default (`luavm/lua/default.lua`) | `default.lua`

