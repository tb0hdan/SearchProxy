all: slow-lint lint test build

BUILD = $(shell git rev-parse HEAD)
BDATE = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p_UTC')
GO_VERSION = $(shell go version|awk '{print $$3}')
VERSION = $(shell cat ./VERSION)
LINTS = lint-main.go lint-geoip lint-mirrorsearch lint-mirrorsort lint-server lint-util/miscellaneous lint-util/network lint-util/system lint-workerpool
TESTS = test-geoip test-mirrorsearch test-mirrorsort test-server test-util/network test-util/system test-workerpool
COVERAGE = coverage-main.go coverage-geoip coverage-mirrorsearch coverage-mirrorsort coverage-server coverage-util/miscellaneous coverage-util/network coverage-util/system coverage-workerpool

geo:
	@go get -u github.com/maxmind/geoipupdate/cmd/geoipupdate
	@geoipupdate -d ./ -f ./etc/geoipupdate.cfg

build: mod test geo build-only

mod:
	@go mod why

build-only:
	@go build -a -trimpath -tags netgo -installsuffix netgo -v -x -ldflags "-s -w -X main.Build=$(BUILD) -X main.BuildDate=$(BDATE) -X main.GoVersion=$(GO_VERSION) -X main.Version=$(VERSION)" -o searchproxy *.go
	@strip -S -x searchproxy


dockerimage:
	@docker build -t tb0hdan/searchproxy .

deps:
	@go get -u golang.org/x/lint/golint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.21.0

lint:
	@golangci-lint run --enable-all --disable=gosec

test: $(TESTS)

$(TESTS):
	@go test -bench=. -v -benchmem -race ./$(shell echo $@|awk -F'test-' '{print $$2}')

slow-lint: $(LINTS)

$(LINTS):
	@golint -set_exit_status=1 $(shell echo $@|awk -F'lint-' '{print $$2}')

codecov: $(COVERAGE)

$(COVERAGE):
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./$(shell echo $@|awk -F'coverage-' '{print $$2}')

tag:
	@git tag -a v$(VERSION) -m v$(VERSION)
	@git push --tags
