# Underground Empire (UE) Makefile

# Variables
BINARY_NAME=ued
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building Underground Empire daemon..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ued
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install the binary
.PHONY: install
install: build
	@echo "Installing Underground Empire daemon..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installation complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Check for external blockchain framework references (should be 0)
.PHONY: check
check: fmt lint test
	@echo "All checks passed!"

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation available at http://localhost:6060"

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go mod download
	@echo "Development setup complete"

# Release build
.PHONY: release
release: clean
	@echo "Building release version..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/ued
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/ued
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/ued
	@echo "Release builds complete"

# Help
.PHONY: help
help:
	@echo "Underground Empire (UE) Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  install      - Install the binary"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  check        - Run all checks"
	@echo "  docs         - Generate documentation"
	@echo "  dev-setup    - Setup development environment"
	@echo "  release      - Build release versions"
	@echo "  help         - Show this help" 