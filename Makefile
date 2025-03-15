all: build
.PHONY: all build

VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X 'main.BuildVersion=$(VERSION)'
GOLDFLAGS += -X 'main.BuildTime=$(BUILDTIME)'

build:
	@CGO_ENABLED=0 go build -ldflags="-s -w $(GOLDFLAGS)" -o ./bin/etherspace ./cmd/main.go

install:
	cp ./bin/etherspace /usr/local/bin

test:
	@cd pkg/p2p/network; grc go test -v --race
	@cd pkg/p2p/network/mock; grc go test -v --race
	@cd pkg/p2p/network/libp2p; grc go test -v --race
	@cd pkg/p2p/dht; grc go test -v --race
	@cd pkg/p2p/test; grc go test -v --race
	@cd pkg/security/crypto; grc go test -v --race
	@cd internal/crypto; grc go test -v --race
