# gosd

replace RSD-backend

## Build

* golang >= 1.11.x
* make

```sh
make
```

## Run

```sh
MQTT_URL=mqtt://admin:public@localhost:1883 \
DATABASE_URL=postgres://postgres:password@localhost/gosd?sslmode=disable \
./gosd -migrate
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
  "attachments": {
    "file":"30"
  }
}
```

