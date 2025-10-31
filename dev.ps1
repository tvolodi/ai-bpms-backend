# AI-BPMS Development Script for Windows
param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Build-Application {
    Write-ColorOutput "Building AI-BPMS Backend..." "Green"
    
    if (!(Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    
    go build -o "bin/ai-bpms-backend.exe" "./cmd/server/main.go"
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput "Build completed successfully!" "Green"
    } else {
        Write-ColorOutput "Build failed!" "Red"
        exit 1
    }
}

function Build-MigrationTool {
    Write-ColorOutput "Building migration tool..." "Green"
    
    if (!(Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    
    go build -o "bin/migrate.exe" "./cmd/migrate/main.go"
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput "Migration tool built successfully!" "Green"
    } else {
        Write-ColorOutput "Migration tool build failed!" "Red"
        exit 1
    }
}

function Start-Services {
    Write-ColorOutput "Starting development services..." "Green"
    
    docker-compose up -d
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput "Services started successfully!" "Green"
        Write-ColorOutput "Waiting for services to initialize..." "Yellow"
        Start-Sleep -Seconds 10
    } else {
        Write-ColorOutput "Failed to start services!" "Red"
        exit 1
    }
}

function Stop-Services {
    Write-ColorOutput "Stopping development services..." "Yellow"
    
    docker-compose down
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput "Services stopped successfully!" "Green"
    } else {
        Write-ColorOutput "Failed to stop services!" "Red"
    }
}

function Run-Migrations {
    Write-ColorOutput "Running database migrations..." "Green"
    
    Build-MigrationTool
    
    & "./bin/migrate.exe"
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput "Migrations completed successfully!" "Green"
    } else {
        Write-ColorOutput "Migrations failed!" "Red"
        exit 1
    }
}

function Start-Application {
    Write-ColorOutput "Starting AI-BPMS Backend..." "Green"
    
    Build-Application
    
    & "./bin/ai-bpms-backend.exe"
}

function Test-Services {
    Write-ColorOutput "Testing service health..." "Green"
    
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8081/health" -UseBasicParsing -TimeoutSec 5
        if ($response.StatusCode -eq 200) {
            Write-ColorOutput "API Server is healthy" "Green"
        } else {
            Write-ColorOutput "API Server returned status: $($response.StatusCode)" "Red"
        }
    } catch {
        Write-ColorOutput "API Server is not responding" "Red"
    }
    
    Write-ColorOutput "Docker services status:" "Yellow"
    docker-compose ps
}

function Start-Development {
    Write-ColorOutput "=== Starting AI-BPMS Development Environment ===" "Green"
    
    Start-Services
    Run-Migrations
    Test-Services
    
    Write-ColorOutput "=== Development Environment Ready ===" "Green"
    Write-ColorOutput "API Server: http://localhost:8081" "Green"
    Write-ColorOutput "Health Check: http://localhost:8081/health" "Green"
    Write-ColorOutput "Swagger Docs: http://localhost:8081/swagger/index.html" "Green"
    Write-ColorOutput "Keycloak: http://localhost:8080 (admin/admin)" "Green"
    Write-ColorOutput "To start the application: .\dev.ps1 run" "Yellow"
}

function Run-Tests {
    Write-ColorOutput "Running tests..." "Green"
    
    go test -v ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput "All tests passed!" "Green"
    } else {
        Write-ColorOutput "Some tests failed!" "Red"
        exit 1
    }
}

function Format-Code {
    Write-ColorOutput "Formatting Go code..." "Green"
    
    go fmt ./...
    go mod tidy
    
    Write-ColorOutput "Code formatting completed!" "Green"
}

function Clean-Build {
    Write-ColorOutput "Cleaning build artifacts..." "Yellow"
    
    go clean
    
    if (Test-Path "bin") {
        Remove-Item -Path "bin" -Recurse -Force
        Write-ColorOutput "Removed bin directory" "Green"
    }
    
    if (Test-Path "logs") {
        Remove-Item -Path "logs" -Recurse -Force
        Write-ColorOutput "Removed logs directory" "Green"
    }
    
    Write-ColorOutput "Clean completed!" "Green"
}

function Show-Help {
    Write-ColorOutput "AI-BPMS Backend Development Script" "Green"
    Write-ColorOutput "=================================" "Green"
    Write-Host ""
    Write-ColorOutput "Usage: .\dev.ps1 [command]" "Yellow"
    Write-Host ""
    Write-ColorOutput "Available commands:" "Green"
    Write-Host "  help          Show this help message"
    Write-Host "  dev           Start full development environment"
    Write-Host "  build         Build the application"
    Write-Host "  run           Build and run the application"
    Write-Host "  test          Run tests"
    Write-Host "  fmt           Format code"
    Write-Host "  clean         Clean build artifacts"
    Write-Host "  services-up   Start Docker services"
    Write-Host "  services-down Stop Docker services"
    Write-Host "  migrate       Run database migrations"
    Write-Host "  health        Check service health"
    Write-Host ""
    Write-ColorOutput "Examples:" "Yellow"
    Write-Host "  .\dev.ps1 dev         # Start development environment"
    Write-Host "  .\dev.ps1 run         # Run the application"
    Write-Host "  .\dev.ps1 test        # Run tests"
}

switch ($Command.ToLower()) {
    "help" { Show-Help }
    "dev" { Start-Development }
    "build" { Build-Application }
    "run" { Start-Application }
    "test" { Run-Tests }
    "fmt" { Format-Code }
    "clean" { Clean-Build }
    "services-up" { Start-Services }
    "services-down" { Stop-Services }
    "migrate" { Run-Migrations }
    "health" { Test-Services }
    default {
        Write-ColorOutput "Unknown command: $Command" "Red"
        Write-ColorOutput "Use '.\dev.ps1 help' to see available commands" "Yellow"
        exit 1
    }
}