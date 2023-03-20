OS=
ARCH=
PROFIX=

# docker, nerdctl
CTR=docker
VERSION=$(shell git describe --tags || git rev-parse --short HEAD || echo "unknown version")
BUILD_DATE=$(shell date +%FT%T%z)
LD_FLAGS='-X "sb.im/gosd/version.Version=$(VERSION)" -X "sb.im/gosd/version.Date=$(BUILD_DATE)"'

# https://pkg.go.dev/time/tzdata
GO_TAGS=timetzdata
GOBUILD=CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
				go build -tags $(GO_TAGS) -ldflags $(LD_FLAGS)

all: build

cli-redis:
	$(CTR) run -it --network=host \
		-e IREDIS_URL=redis://localhost:6379/1 \
		dbcliorg/iredis

cli-pg:
	$(CTR) run -it --network=host \
		-e PGHOST=localhost \
		-e PGPORT=5432 \
		-e PGUSER=postgres \
		-e PGPASSWORD=password \
		-e PGDATABASE=gosd \
		dbcliorg/pgcli

build:
	$(GOBUILD)

.PHONY: run
run:
	DEMO_MODE=true go run main.go server -v

.PHONY: swagger
swagger:
	# go install github.com/swaggo/swag/cmd/swag@leatest
	swag init -g app/app.go -o swag
	@ rm swag/docs.go
	@ rm swag/swagger.yaml

test: test-luavm test-broker test-simulation

test-luavm:
	go clean -testcache && go test ./app/luavm -cover -v

test-broker:
	go test ./mqttd -cover -v

test-simulation:
	go test ./tests/e2e/... -cover -v

# \(statements\)(?:\s+)?(\d+(?:\.\d+)?%)
# https://stackoverflow.com/questions/61246686/go-coverage-over-multiple-package-and-gitlab-coverage-badge
cover:
	go test ./app/... -coverprofile=coverage.txt -covermode=atomic
	go tool cover -func coverage.txt

install:
	install -Dm755 gosd -t ${PROFIX}/usr/bin/gosd
	install -Dm644 gosd.service -t ${PROFIX}/lib/systemd/system/

clean:
	go clean
	@ rm -r swag

