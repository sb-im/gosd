# gosd

replace RSD-backend

```sh
go build
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
  "attachments": {
    "file":"30"
  }
}
```


