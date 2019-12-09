all: build

BUILD = $(shell git rev-parse HEAD)
BDATE = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p_UTC')
GO_VERSION = $(shell go version|awk '{print $$3}')
VERSION = $(shell cat ./VERSION)

TESTS = test-geoip test-memcache test-mirrorsort test-server test-util/network test-util/system test-workerpool

geo:
	@go get -u github.com/maxmind/geoipupdate/cmd/geoipupdate
	@geoipupdate -d ./ -f ./etc/geoipupdate.cfg

build: geo
	@go mod why
	@go build -tags netgo -installsuffix netgo -v -x -ldflags "-s -w -X main.Build=$(BUILD) -X main.BuildDate=$(BDATE) -X main.GoVersion=$(GO_VERSION) -X main.Version=$(VERSION)" -o searchproxy *.go

build-only:
	@go build -v -x -ldflags "-s -w -X main.Build=$(BUILD) -X main.BuildDate=$(BDATE) -X main.GoVersion=$(GO_VERSION) -X main.Version=$(VERSION)" -o searchproxy *.go

dockerimage:
	@docker build -t tb0hdan/searchproxy .

lint:
	@golangci-lint run --enable-all --disable=gosec

test: $(TESTS)

$(TESTS):
	@go test -bench=. -v -benchmem -race ./$(shell echo $@|awk -F'test-' '{print $$2}')
