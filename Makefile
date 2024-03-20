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
build: build-cli build-dc build-grpc

build-cli:
	go build -o build/robopac-cli     ./cmd/cli

build-dc:
	go build -o build/robopac-discord ./cmd/discord

build-grpc:
	go build -o build/robopac-grpc ./cmd/grpc

.PHONY: build
