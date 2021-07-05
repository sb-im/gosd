# gosd

replace RSD-backend

## Build

* golang >= 1.12.x
* make

```sh
make
```

## Dependencies

* PostgreSQL >= 13
* Redis
* Mqtt broker (Mosquitto Or Emqx)
  * Mosquitto >= 1.6

## Run

```sh
# Create database tables
MQTT_URL=mqtt://admin:public@localhost:1883 \
REDIS_URL=redis://localhost:6379/0 \
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd database migrate


# Create User
# gosd users add <username> <password>
MQTT_URL=mqtt://admin:public@localhost:1883 \
REDIS_URL=redis://localhost:6379/0 \
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd users add demo demodemo

# Run
MQTT_URL=mqtt://admin:public@localhost:1883 \
REDIS_URL=redis://localhost:6379/0 \
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd
```

## Environment Variables

Variable Name        | Description                                               | Default Value
-------------------- | --------------------------------------------------------- | -------------------------------------------------------------
`DEBUG`              | Set the value to `1` to enable debug logs                 | Off
`DATABASE_URL`       | Postgresql connection parameters                          | `postgres://postgres:password@localhost/gosd?sslmode=disable`
`DATABASE_MAX_CONNS` | Maximum number of database connections                    | 20
`DATABASE_MIN_CONNS` | Minimum number of database connections                    | 1
`LISTEN_ADDR`        | Address to listen on (use absolute path for Unix socket)  | `127.0.0.1:8000`
`PORT`               | Override `LISTEN_ADDR` to `0.0.0.0:$PORT` (PaaS)          | None
`BASE_URL`           | Base URL to generate HTML links and base path for cookies | `http://localhost/`
`MQTT_URL`           | MQTT broker Server address                                | `mqtt://admin:public@localhost:1883`
`MQTT_CLIENT_ID`     | MQTT Client ID                                            | `cloud.0`
`REDIS_URL`          | Redis Server URL, **Only use db `0`**                     | `redis://localhost:6379/0`
`LOG_FILE`           | Log File                                                  | `STDOUT`
`LOG_LEVEL`          | Log Level: `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | `info`
`OAUTH_CLIENT_ID`    | OAuth id | `000000`
`OAUTH_CLIENT_SECRET` | OAuthSecret | `999999`


### Auth

```sh
curl -X POST 'http://localhost:8000/gosd/api/v2/oauth/token' \
-F grant_type=password \
-F client_id=000000 \
-F client_secret=999999 \
-F scope=read \
-F username=demo \
-F password=demodemo
```


### New Plan

FormData

```sh
curl -X POST -F "name=233" \
-F "description=test"
-F "node_id"=1 \
-F "file=@go.mod" \
localhost:8000/v1/plans
```

Or

json

```sh
curl -X POST localhost:8000/v1/plans \
-d '{"name": "233", "description": "test", "node_id": 1}'
```

```json
{
  "id":35,
  "name":"233",
  "description":"test",
  "node_id":1,
  "files": {
    "file":"30"
  }
}
```

