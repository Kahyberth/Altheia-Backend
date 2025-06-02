# Makefile for Altheia Backend

.PHONY: test test-verbose test-coverage test-auth test-clinical test-utils benchmark clean build help

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run all tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-auth     - Run authentication module tests"
	@echo "  test-clinical - Run clinical module tests"
	@echo "  test-utils    - Run utility function tests"
	@echo "  benchmark     - Run benchmark tests"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean test cache and build artifacts"
	@echo "  deps          - Download dependencies"

# Run all tests
test:
	@echo "Running all tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	@echo "Generating detailed coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run auth module tests
test-auth:
	@echo "Running authentication module tests..."
	go test -v ./internal/auth/...

# Run clinical module tests
test-clinical:
	@echo "Running clinical module tests..."
	go test -v ./internal/clinical/...

# Run utils tests
test-utils:
	@echo "Running utility function tests..."
	go test -v ./pkg/utils/...

# Run benchmark tests
benchmark:
	@echo "Running benchmark tests..."
	go test -bench=. ./...

# Run benchmark with memory stats
benchmark-mem:
	@echo "Running benchmark tests with memory statistics..."
	go test -bench=. -benchmem ./...

# Run specific benchmark
benchmark-auth:
	@echo "Running authentication benchmark tests..."
	go test -bench=. ./internal/auth/...

benchmark-utils:
	@echo "Running utils benchmark tests..."
	go test -bench=. ./pkg/utils/...

# Build the application
build:
	@echo "Building application..."
	go build -o bin/altheia ./cmd/main.go

# Clean test cache and build artifacts
clean:
	@echo "Cleaning test cache and build artifacts..."
	go clean -testcache
	go clean -cache
	rm -f coverage.out coverage.html
	rm -rf bin/

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Run tests in race condition detection mode
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Run tests with short flag (skip long running tests)
test-short:
	@echo "Running short tests..."
	go test -short ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run all quality checks
check: fmt vet test
	@echo "All quality checks completed!"

# Set up test environment
setup-test:
	@echo "Setting up test environment..."
	export JWT_SECRET=test-secret-key-for-testing-purposes

# Run tests with environment setup
test-env: setup-test test 