# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=squish
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_GO=./cmd/squish/main.go

# Build directory
BUILD_DIR=build

# Version (to be overwritten by semantic-release)
VERSION ?= development

# Ldflags for stripping debug information and setting version
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

# Supported OSs and Architectures
PLATFORMS=darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64

.PHONY: all test clean build build-all

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v $(MAIN_GO)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf $(BUILD_DIR)

build-all: clean
	mkdir -p $(BUILD_DIR)
	$(foreach platform,$(PLATFORMS),\
		$(eval GOOS=$(word 1,$(subst /, ,$(platform))))\
		$(eval GOARCH=$(word 2,$(subst /, ,$(platform))))\
		$(eval EXTENSION=$(if $(filter windows,$(GOOS)),.exe))\
		$(eval BINARY=$(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(EXTENSION))\
		echo "Building for $(GOOS)/$(GOARCH)..." && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(BINARY) -v $(MAIN_GO) \
		&& if [ "$(GOOS)" = "linux" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				cp $(BINARY) $(BUILD_DIR)/$(BINARY_NAME)-linux-x86_64; \
			elif [ "$(GOARCH)" = "arm64" ]; then \
				cp $(BINARY) $(BUILD_DIR)/$(BINARY_NAME)-linux-aarch64; \
			fi; \
		fi; \
	)

# Cross compilation for Unix
build-unix:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v $(MAIN_GO)

# Run the application
run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v $(MAIN_GO)
	./$(BINARY_NAME)

# Dependencies
deps:
	$(GOGET) -v ./...