.PHONY: install server test

## install: fetches go modules
install:
	go mod tidy; \
	go mod download \

## server: runs the server with -race
server:
	go run -race cmd/server/main.go

## test: runs tests
test:
	go test -race ./...

## help: prints help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'