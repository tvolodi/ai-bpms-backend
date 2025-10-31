# AI-BPMS Backend Makefile
# Development and build automation

# Variables
BINARY_NAME=ai-bpms-backend
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server/main.go
MIGRATE_PATH=./cmd/migrate/main.go

# Docker variables
DOCKER_IMAGE=ai-bpms/backend
DOCKER_TAG=latest

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build clean test deps fmt lint run migrate docker-build docker-run docker-compose-up docker-compose-down dev setup

# Default target
all: clean deps fmt lint test build

# Show help
help:
	@echo "$(GREEN)AI-BPMS Backend Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@echo "  help              Show this help message"
	@echo "  setup             Initial project setup (install tools, deps)"
	@echo "  build             Build the application"
	@echo "  clean             Clean build artifacts"
	@echo "  test              Run tests"
	@echo "  test-coverage     Run tests with coverage"
	@echo "  deps              Download dependencies"
	@echo "  deps-update       Update dependencies"
	@echo "  fmt               Format code"
	@echo "  lint              Run linter"
	@echo "  run               Run the application"
	@echo "  migrate           Run database migrations"
	@echo "  migrate-rollback  Rollback last migration"
	@echo "  dev               Start development environment"
	@echo "  docker-build      Build Docker image"
	@echo "  docker-run        Run Docker container"
	@echo "  docker-compose-up Start services with docker-compose"
	@echo "  docker-compose-down Stop docker-compose services"
	@echo "  swagger           Generate Swagger docs"

# Initial project setup
setup:
	@echo "$(GREEN)Setting up development environment...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "$(GREEN)Installing development tools...$(NC)"
	@$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GOGET) -u github.com/swaggo/swag/cmd/swag
	@echo "$(GREEN)Setup complete!$(NC)"

# Build the application
build:
	@echo Building $(BINARY_NAME)...
	@if not exist bin mkdir bin
	@$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo Build complete: $(BINARY_PATH)

# Build migration tool
build-migrate:
	@echo Building migration tool...
	@if not exist bin mkdir bin
	@$(GOBUILD) -o ./bin/migrate $(MIGRATE_PATH)
	@echo Migration tool built: ./bin/migrate

# Clean build artifacts
clean:
	@echo Cleaning...
	@$(GOCLEAN)
	@if exist bin rmdir /s /q bin 2>nul || rm -rf bin/ 2>/dev/null || echo Cleaned
	@if exist logs rmdir /s /q logs 2>nul || rm -rf logs/ 2>/dev/null || echo Logs cleaned

# Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) verify

# Update dependencies
deps-update:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@$(GOMOD) tidy
	@$(GOGET) -u ./...

# Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GOFMT) -w .
	@$(GOCMD) mod tidy

# Run linter
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@$(GOLINT) run

# Run the application
run: build
	@echo "$(GREEN)Starting $(BINARY_NAME)...$(NC)"
	@$(BINARY_PATH)

# Run database migrations
migrate: build-migrate
	@echo "$(GREEN)Running database migrations...$(NC)"
	@./bin/migrate

# Rollback last migration
migrate-rollback: build-migrate
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@./bin/migrate -rollback

# Start development environment
dev:
	@echo Starting development environment...
	docker-compose up -d
	@echo Waiting for services to start...
	@ping -n 11 127.0.0.1 > nul 2>&1 || sleep 10 2>/dev/null || echo Waiting...
	@echo Running migrations...
	@$(MAKE) migrate
	@echo Starting application...
	@$(MAKE) run

# Generate Swagger documentation
swagger:
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@swag init -g cmd/server/main.go -o ./docs
	@echo "$(GREEN)Swagger docs generated in ./docs$(NC)"

# Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

# Run Docker container
docker-run: docker-build
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run -p 8081:8081 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Start services with docker-compose
docker-compose-up:
	@echo "$(GREEN)Starting services with docker-compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services started. Check with: docker-compose ps$(NC)"

# Stop docker-compose services
docker-compose-down:
	@echo "$(YELLOW)Stopping docker-compose services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Services stopped$(NC)"

# Development watch mode (requires air)
dev-watch:
	@echo "$(GREEN)Starting development with hot reload...$(NC)"
	@air

# Install air for hot reloading
install-air:
	@echo "$(GREEN)Installing air for hot reloading...$(NC)"
	@$(GOGET) -u github.com/cosmtrek/air

# Database operations
db-reset: docker-compose-down docker-compose-up
	@echo "$(YELLOW)Resetting database...$(NC)"
	@sleep 5
	@$(MAKE) migrate

# Check services health
health-check:
	@echo Checking services health...
	@powershell -Command "try { Invoke-WebRequest -Uri http://localhost:8081/health -UseBasicParsing | Out-Null; Write-Host 'API service is healthy' } catch { Write-Host 'API service is down' -ForegroundColor Red }" 2>nul || curl -f http://localhost:8081/health 2>/dev/null || echo API service check failed
	@docker-compose ps

# Show logs
logs:
	@echo "$(GREEN)Showing application logs...$(NC)"
	@docker-compose logs -f

# Production build
build-prod:
	@echo "$(GREEN)Building for production...$(NC)"
	@CGO_ENABLED=0 GOOS=linux $(GOBUILD) -ldflags="-w -s" -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Production build complete$(NC)"

# Security scan
security:
	@echo "$(GREEN)Running security scan...$(NC)"
	@$(GOGET) -u github.com/securecodewarrior/gosec/v2/cmd/gosec
	@gosec ./...