unit_test:
	go test ./...

race_test:
	go test ./... --race

fmt:
	gofumpt -l -w .

install-tools:
	@echo "Installing devtools"
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1

install-packages:
	go mod tidy

check:
	golangci-lint run \
        --build-tags "${BUILD_TAG}" \
        --timeout=20m0s \
        --enable=gofmt \
        --enable=unconvert \
        --enable=unparam \
        --enable=asciicheck \
        --enable=misspell \
        --enable=revive \
        --enable=decorder \
        --enable=reassign \
        --enable=usestdlibvars \
        --enable=nilerr \
        --enable=gosec \
        --enable=exportloopref \
        --enable=whitespace \
        --enable=goimports \
        --enable=gocyclo \
        --enable=lll


build-bot:
	go build -o build/main cmd/app.go