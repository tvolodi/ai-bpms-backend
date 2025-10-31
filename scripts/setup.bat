@echo off
rem Development setup script for AI-BPMS Backend (Windows)

echo Setting up AI-BPMS Backend development environment...

rem Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Go is not installed. Please install Go 1.21+ first.
    exit /b 1
)

echo Go detected

rem Check if Docker is installed
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not installed. Please install Docker first.
    exit /b 1
)

echo Docker detected

rem Check if Docker Compose is installed
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker Compose is not installed. Please install Docker Compose first.
    exit /b 1
)

echo Docker Compose detected

rem Install Go development tools
echo Installing Go development tools...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/cosmtrek/air@latest

rem Download dependencies
echo Downloading Go dependencies...
go mod download
go mod tidy

rem Create .env file if it doesn't exist
if not exist .env (
    echo Creating .env file from template...
    copy .env.example .env
    echo Please update .env file with your configuration
)

rem Create required directories
echo Creating required directories...
if not exist logs mkdir logs
if not exist bin mkdir bin
if not exist tmp mkdir tmp

rem Start development services
echo Starting development services...
docker-compose up -d

rem Wait for services to be ready
echo Waiting for services to start...
timeout /t 10 /nobreak

echo.
echo Development environment setup complete!
echo.
echo Quick start commands:
echo   make dev              # Start development environment
echo   make run              # Run the application
echo   make test             # Run tests
echo   make help             # Show all available commands
echo.
echo Services:
echo   API Server:           http://localhost:8081
echo   Swagger Docs:         http://localhost:8081/swagger/index.html
echo   Health Check:         http://localhost:8081/health
echo   PostgreSQL:           localhost:5432
echo   NATS:                 localhost:4222
echo   Redis:                localhost:6379
echo   Keycloak:             http://localhost:8080 (admin/admin)
echo   Prometheus:           http://localhost:9090
echo   Grafana:              http://localhost:3000 (admin/admin)
echo.
echo Next steps:
echo 1. Update .env file with your configuration
echo 2. Run 'make dev' to start development
echo 3. Visit http://localhost:8081/health to verify everything is working