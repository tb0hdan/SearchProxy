all: build


build:
	@go mod why
	@go build -v -x -ldflags "-s -w" -o searchproxy *.go

dockerimage:
	@cp -r /usr/local/etc/openssl ./ssl
	@docker build -t tb0hdan/searchproxy .

lint:
	@golangci-lint run --enable-all --disable=gosec
