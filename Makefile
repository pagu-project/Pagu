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

packages:
	go mod tidy


### Formatting, linting, and vetting
fmt:
	gofumpt -l -w .
	go mod tidy
	godot -w .

check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

### building
discord-bot:
	go build -o build/main cmd/app.go
        