OS=
ARCH=
PROFIX=
GO_TEST=./state ./cache
VERSION=$(shell git describe --tags || git rev-parse --short HEAD || echo "unknown version")
BUILD_DATE=$(shell date +%FT%T%z)
LD_FLAGS='-X "miniflux.app/version.Version=$(VERSION)" -X "miniflux.app/version.BuildDate=$(BUILD_DATE)"'
GOBUILD=CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
				go build -ldflags $(LD_FLAGS)

all: build

generate:
	@ go generate

build: generate
	$(GOBUILD)

run: generate
	@ go run main.go --debug --noauth

test: generate
	go test ${GO_TEST} -cover -v

test-broker: generate
	go test ./mqttd -cover -v

test-simulation:
	go test ./luavm ./integration -cover -v

# \(statements\)(?:\s+)?(\d+(?:\.\d+)?%)
# https://stackoverflow.com/questions/61246686/go-coverage-over-multiple-package-and-gitlab-coverage-badge
cover: generate
	go test ${GO_TEST} -coverprofile profile.cov
	go tool cover -func profile.cov
	@rm profile.cov

install:
	install -Dm755 gosd -t ${PROFIX}/usr/bin/gosd
	install -Dm644 gosd.service -t ${PROFIX}/lib/systemd/system/

clean:
	go clean

