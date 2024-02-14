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

### mock

mock:
	mockgen -source=./client/interface.go      -destination=./client/mock.go      -package=client
	mockgen -source=./wallet/interface.go      -destination=./wallet/mock.go      -package=wallet
	mockgen -source=./store/interface.go       -destination=./store/mock.go       -package=store
	mockgen -source=./twitter_api/interface.go -destination=./twitter_api/mock.go -package=twitter_api
	mockgen -source=./turboswap/interface.go   -destination=./turboswap/mock.go   -package=turboswap

### Formatting, linting, and vetting
fmt:
	gofumpt -l -w .
	go mod tidy
	
check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

### building
build:
	go build -o build/robopac-discord ./cmd/discord
	go build -o build/robopac-cmd     ./cmd/cmd

build-cmd:
	go build -o build/robopac-cmd     ./cmd/cmd

build-dc:
	go build -o build/robopac-discord ./cmd/discord

.PHONY: build
