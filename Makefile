.PHONY: build run dev test clean docker-build docker-run

# Application name
APP_NAME := lwnra-devo-api

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) ./cmd/server

# Run the application
run: build
	@echo "Starting $(APP_NAME)..."
	@./bin/$(APP_NAME)

# Run in development mode with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server..."
	@air -c .air.toml

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f devotionals.db

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Start server with environment variables
start: build
	@echo "Starting server with environment variables..."
	@PORT=8080 DB_PATH=devotionals.db ENVIRONMENT=development ./bin/$(APP_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode with auto-reload"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  install-deps   - Install development dependencies"
	@echo "  start          - Start server with environment variables"
	@echo "  help           - Show this help message"
