.PHONY: help install run build format lint test coverage clean wire-gen

APP_NAME = github.com/haidang666/go-app
CMD_PATH = ./cmd/server
BIN_PATH = ./bin
BINARY_NAME = go-app

# Default target
help:
	@echo "Go App - Available targets:"
	@echo "  make install       - Install dependencies"
	@echo "  make run           - Run the server"
	@echo "  make build         - Build the binary"
	@echo "  make format        - Format code with go fmt"
	@echo "  make lint          - Lint code with go vet"
	@echo "  make test          - Run all tests"
	@echo "  make coverage      - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make wire-gen      - Generate wire dependency injection"

install:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

run:
	@echo "Running server..."
	go run $(CMD_PATH)/main.go

build: clean
	@echo "Building binary..."
	mkdir -p $(BIN_PATH)
	go build -o $(BIN_PATH)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Binary built: $(BIN_PATH)/$(BINARY_NAME)"

format:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	go vet ./...

test:
	@echo "Running tests..."
	go test -v ./...

coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_PATH)
	rm -f coverage.out coverage.html

wire-gen:
	@echo "Generating wire dependencies..."
	go generate ./internal/app