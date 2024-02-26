.PHONY: all setup dunecli build lint yamllint

all: lint test build

setup: bin/golangci-lint
	go mod download

dunecli: lint
	go build -o dunecli cmd/main.go

build: dunecli

bin:
	mkdir -p bin

bin/golangci-lint: bin
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...
	bin/golangci-lint -c .golangci.yml run ./...
	go mod tidy

test:
	go test -timeout=10s -race -cover -bench=. -benchmem ./...
