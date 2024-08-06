# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=squish
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/squish

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/squish
	./$(BINARY_NAME)

deps:
	$(GOGET) ./...
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/squish

# Linting
lint:
	golangci-lint run

# Format code
fmt:
	gofmt -s -w .

# Check if code is formatted
fmt-check:
	test -z $$(gofmt -l .)

# Generate mocks for testing
mocks:
	mockgen -source=pkg/esbuild/plugin.go -destination=pkg/esbuild/mocks/mock_plugin.go

# Self-bundle (assuming squish can bundle itself)
self-bundle: build
	./$(BINARY_NAME) --src ./cmd/squish --dist ./dist

# Install golangci-lint
install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1

.PHONY: all build test clean run deps build-linux lint fmt fmt-check mocks self-bundle install-linter