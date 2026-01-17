.PHONY: help build test clean fmt fmt-check deps lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all packages
	@echo "Building coil..."
	go build ./...

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

clean: ## Clean build artifacts and cache
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	go clean -cache -testcache

fmt: ## Format code to 80-column limit
	@echo "Formatting code with golines..."
	@$(shell go env GOPATH)/bin/golines -w -m 80 --base-formatter=gofmt \
		--shorten-comments \
		--ignore-generated \
		.
	@echo "✓ Code formatted"

fmt-check: ## Check if code adheres to 80-column limit
	@echo "Checking code formatting..."
	@$(shell go env GOPATH)/bin/golines -m 80 --base-formatter=gofmt \
		--shorten-comments \
		--ignore-generated \
		--dry-run \
		. || (echo "✗ Formatting issues. Run 'make fmt'" && exit 1)
	@echo "✓ Code formatting OK"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Skipping..."; \
		echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies updated"

setup: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/segmentio/golines@latest
	@echo "✓ Tools installed"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ Vet passed"

check: fmt-check vet test ## Run all checks (format, vet, tests)
