PACKAGES=$(shell go list ./... | grep -v 'tests')

### Testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./... -covermode=atomic

race_test:
	go test ./... --race

### dev tools
devtools:
	@echo "Installing devtools"
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	go install go.uber.org/mock/mockgen@latest
	go install github.com/bufbuild/buf/cmd/buf@latest

### mock
mock:
	mockgen -source=./client/interface.go      -destination=./client/mock.go      -package=client
	mockgen -source=./wallet/interface.go      -destination=./wallet/mock.go      -package=wallet

### proto file generate
proto:
	rm -rf grpc/gen/go
	cd grpc/buf && buf generate --template buf.gen.yaml ../proto

### Formatting, linting, and vetting
fmt:
	gofumpt -l -w .
	go mod tidy

check:
	golangci-lint run --timeout=20m0s

### building
build: build-cli build-discord-mainnet build-discord-testnet build-grpc build-telegram build-http

build-cli:
	go build -o build/pagu-cli     ./cmd/cli

build-grpc:
	go build -o build/pagu-grpc ./cmd/mainnet/grpc

build-discord-mainnet:
	go build -o build/pagu-discord-mainnet ./cmd/mainnet/discord

build-discord-testnet:
	go build -o build/pagu-discord-testnet ./cmd/testnet/discord

build-telegram:
	go build -o build/pagu-telegram ./cmd/mainnet/telegram

build-http:
	go build -o build/pagu-http ./cmd/mainnet/http

### pre commit
pre-commit: mock proto fmt check unit_test
	@echo pre commit commands...

.PHONY: build
