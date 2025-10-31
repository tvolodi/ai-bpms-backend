#!/bin/bash

# Development setup script for AI-BPMS Backend

set -e

echo "🚀 Setting up AI-BPMS Backend development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ "$GO_VERSION" < "1.21" ]]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go 1.21+."
    exit 1
fi

echo "✅ Go $GO_VERSION detected"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

echo "✅ Docker detected"

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker Compose detected"

# Install Go development tools
echo "📦 Installing Go development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/cosmtrek/air@latest

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  Please update .env file with your configuration"
fi

# Create required directories
echo "📁 Creating required directories..."
mkdir -p logs
mkdir -p bin
mkdir -p tmp

# Start development services
echo "🐳 Starting development services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Run database migrations
echo "🗄️  Running database migrations..."
make migrate

# Generate Swagger documentation
echo "📚 Generating API documentation..."
make swagger

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "Quick start commands:"
echo "  make dev              # Start development environment"
echo "  make run              # Run the application"
echo "  make test             # Run tests"
echo "  make help             # Show all available commands"
echo ""
echo "Services:"
echo "  API Server:           http://localhost:8081"
echo "  Swagger Docs:         http://localhost:8081/swagger/index.html"
echo "  Health Check:         http://localhost:8081/health"
echo "  PostgreSQL:           localhost:5432"
echo "  NATS:                 localhost:4222"
echo "  Redis:                localhost:6379"
echo "  Keycloak:             http://localhost:8080 (admin/admin)"
echo "  Prometheus:           http://localhost:9090"
echo "  Grafana:              http://localhost:3000 (admin/admin)"
echo ""
echo "Next steps:"
echo "1. Update .env file with your configuration"
echo "2. Run 'make dev' to start development"
echo "3. Visit http://localhost:8081/health to verify everything is working"