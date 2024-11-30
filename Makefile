# Variables
APP_CMD = cmd/app/main.go
APP_NAME = movie-festival-service
TEST_CMD = go test -v ./tests/...
LINT_CMD = go vet ./... && go fmt ./...
MOCK_SOURCE_DIR = ./internal/repositories/*.go
MOCK_DEST_DIR = ./tests/mocks

# Tasks
.PHONY: all run lint test start clean-test build mocks

# Default target
run: lint mocks test start

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
clean-test:
	@echo "Starting clean test cache"
	@go clean -cache
	@echo "Test cache is cleaned"

# Build the application
build:
	@echo "Building the application..."
	go build -o $(APP_NAME) $(APP_CMD)

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@for f in $(MOCK_SOURCE_DIR); do \
		basename=$$(basename $$f .go); \
		echo "Generating mock for: $$basename"; \
		mockgen -source=$$f -destination=$(MOCK_DEST_DIR)/$$basename"_mock.go" -package=mocks; \
	done
	@echo "Mocks generated successfully!"

# Help target to display usage
help:
	@echo "Usage:"
	@echo "  make run   - Run the application with syntax check and tests"
	@echo "  make lint  - Check code syntax"
	@echo "  make test  - Run all tests"
	@echo "  make start - Run the application directly"
	@echo "  make mocks - Generate repository mocks"
