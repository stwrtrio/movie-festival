# Variables
APP_CMD := cmd/app/main.go
APP_NAME := movie-festival
TEST_CMD := go test -v ./tests/...
LINT_CMD := go vet ./... && go fmt ./...

# Default target
run: lint test start

# Target to check syntax (lint)
lint:
	@echo "Checking code syntax..."
	@$(LINT_CMD)
	@echo "Syntax check passed."

# Target to run tests
test:
	@echo "Running tests..."
	@$(TEST_CMD)
	@echo "All tests passed."

# Target to run the application
start:
	@echo "Starting the application..."
	@go run $(APP_CMD)

# Target to clean test cache (Optional)
clean:
	@echo "Starting clean test cache"
	@go clean -cache
	@echo "Cache is cleaned"

# Help target to display usage
help:
	@echo "Usage:"
	@echo "  make run   - Run the application with syntax check and tests"
	@echo "  make lint  - Check code syntax"
	@echo "  make test  - Run all tests"
	@echo "  make start - Run the application directly"
