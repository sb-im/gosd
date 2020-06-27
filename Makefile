OS=
ARCH=
VERSION=$(shell git describe --tags || git rev-parse --short HEAD || echo "unknown version")
BUILD_DATE=$(shell date +%FT%T%z)
LD_FLAGS='-X "miniflux.app/version.Version=$(VERSION)" -X "miniflux.app/version.BuildDate=$(BUILD_DATE)"'
GOBUILD=GOOS=$(OS) GOARCH=$(ARCH) \
				go build -ldflags $(LD_FLAGS)

all: build

generate:
	@ go generate

build: generate
	$(GOBUILD)

run: generate
	@ go run `ls *.go | grep -v _test.go | grep -v generate.go`

test: generate
	go test ./jsonrpc2mqtt ./state ./luavm -cover

clean:
	go clean

