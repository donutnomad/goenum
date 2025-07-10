# Makefile for goenum project

# Binary name
BINARY_NAME=goenum

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
BUILD_FLAGS=-v

# Source files
SOURCES=main.go plurals.go

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(SOURCES)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary to $GOPATH/bin
.PHONY: install
install:
	$(GOBUILD) $(BUILD_FLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(SOURCES)

# Run the generator on test files
.PHONY: generate
generate: build
	./$(BINARY_NAME) test_complex.go

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build     - Build the binary"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  install   - Install binary to GOPATH/bin"
	@echo "  generate  - Run generator on test files"
	@echo "  fmt       - Format code"
	@echo "  lint      - Run linter"
	@echo "  help      - Show this help message"