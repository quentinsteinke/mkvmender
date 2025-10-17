.PHONY: all build clean test install run-server run-cli help

# Variables
BINARY_CLI=bin/mkvmender
BINARY_SERVER=bin/mkvmender-server
GO=go
GOFLAGS=-v

# Default target
all: clean build

# Build both CLI and server
build:
	@echo "Building CLI..."
	@mkdir -p bin
	@$(GO) build $(GOFLAGS) -o $(BINARY_CLI) ./cmd/cli
	@echo "Building server..."
	@$(GO) build $(GOFLAGS) -o $(BINARY_SERVER) ./cmd/server
	@echo "Build complete!"

# Build CLI only
build-cli:
	@echo "Building CLI..."
	@mkdir -p bin
	@$(GO) build $(GOFLAGS) -o $(BINARY_CLI) ./cmd/cli

# Build server only
build-server:
	@echo "Building server..."
	@mkdir -p bin
	@$(GO) build $(GOFLAGS) -o $(BINARY_SERVER) ./cmd/server

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test ./... -v

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GO) mod tidy
	@$(GO) mod download

# Install CLI to system
install: build-cli
	@echo "Installing mkvmender to /usr/local/bin..."
	@cp $(BINARY_CLI) /usr/local/bin/mkvmender
	@echo "Installation complete!"

# Run server
run-server: build-server
	@echo "Starting server..."
	@$(BINARY_SERVER)

# Run CLI (with arguments: make run-cli ARGS="hash file.mkv")
run-cli: build-cli
	@$(BINARY_CLI) $(ARGS)

# Format code
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build everything (default)"
	@echo "  build        - Build both CLI and server"
	@echo "  build-cli    - Build CLI only"
	@echo "  build-server - Build server only"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  install      - Install CLI to /usr/local/bin"
	@echo "  run-server   - Build and run server"
	@echo "  run-cli      - Build and run CLI (use ARGS='command' to pass arguments)"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  help         - Show this help message"
