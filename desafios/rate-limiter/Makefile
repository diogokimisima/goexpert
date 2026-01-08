.PHONY: help build run test test-coverage docker-up docker-down docker-build clean

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build complete!"

run: ## Run the application locally
	@echo "Running application..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Services started!"
	@echo "Application available at http://localhost:8080"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs: ## Show Docker logs
	@docker-compose logs -f app

docker-rebuild: docker-down docker-build docker-up ## Rebuild and restart Docker containers

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@docker-compose down -v
	@echo "Clean complete!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

integration-test: docker-up ## Run integration tests
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Running integration tests..."
	@bash test.sh
	@$(MAKE) docker-down

.DEFAULT_GOAL := help
