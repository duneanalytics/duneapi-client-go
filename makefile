.PHONY: all setup dunecli build lint yamllint

all: lint test build

setup: bin/golangci-lint
	go mod download

dunecli: lint
	go build -o dunecli cmd/main.go

build: dunecli

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.0

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...
	bin/golangci-lint -c .golangci.yml run ./...
	go mod tidy

test:
	go mod tidy
	go test -timeout=10s -race -benchmem ./...
