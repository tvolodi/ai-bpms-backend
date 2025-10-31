@echo off
rem AI-BPMS Development Commands for Windows
rem Simple batch file alternative to make commands

setlocal enabledelayedexpansion

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="dev" goto dev
if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="test" goto test
if "%1"=="clean" goto clean
if "%1"=="services-up" goto services_up
if "%1"=="services-down" goto services_down
if "%1"=="migrate" goto migrate
if "%1"=="health" goto health

echo Unknown command: %1
echo Use 'dev.bat help' to see available commands
exit /b 1

:help
echo AI-BPMS Backend Development Commands
echo ===================================
echo.
echo Usage: dev.bat [command]
echo.
echo Available commands:
echo   help          Show this help message
echo   dev           Start full development environment
echo   build         Build the application
echo   run           Build and run the application
echo   test          Run tests
echo   clean         Clean build artifacts
echo   services-up   Start Docker services
echo   services-down Stop Docker services
echo   migrate       Run database migrations
echo   health        Check service health
echo.
echo Examples:
echo   dev.bat dev         # Start development environment
echo   dev.bat run         # Run the application
echo   dev.bat test        # Run tests
goto end

:dev
echo Starting development environment...
call :services_up
if !errorlevel! neq 0 exit /b !errorlevel!
echo Waiting for services to start...
timeout /t 10 /nobreak >nul
call :migrate
if !errorlevel! neq 0 exit /b !errorlevel!
call :health
echo.
echo Development environment ready!
echo API Server: http://localhost:8081
echo Health Check: http://localhost:8081/health
echo Swagger Docs: http://localhost:8081/swagger/index.html
echo.
echo To start the application: dev.bat run
goto end

:build
echo Building AI-BPMS Backend...
if not exist bin mkdir bin
go build -o bin\ai-bpms-backend.exe .\cmd\server\main.go
if !errorlevel! neq 0 (
    echo Build failed!
    exit /b 1
)
echo Build completed successfully!
goto end

:run
call :build
if !errorlevel! neq 0 exit /b !errorlevel!
echo Starting AI-BPMS Backend...
.\bin\ai-bpms-backend.exe
goto end

:test
echo Running tests...
go test -v .\...
if !errorlevel! neq 0 (
    echo Some tests failed!
    exit /b 1
)
echo All tests passed!
goto end

:clean
echo Cleaning build artifacts...
go clean
if exist bin rmdir /s /q bin 2>nul
if exist logs rmdir /s /q logs 2>nul
echo Clean completed!
goto end

:services_up
echo Starting Docker services...
docker-compose up -d
if !errorlevel! neq 0 (
    echo Failed to start services!
    exit /b 1
)
echo Services started successfully!
goto end

:services_down
echo Stopping Docker services...
docker-compose down
if !errorlevel! neq 0 (
    echo Failed to stop services!
    exit /b 1
)
echo Services stopped successfully!
goto end

:migrate
echo Running database migrations...
if not exist bin mkdir bin
go build -o bin\migrate.exe .\cmd\migrate\main.go
if !errorlevel! neq 0 (
    echo Migration tool build failed!
    exit /b 1
)
.\bin\migrate.exe
if !errorlevel! neq 0 (
    echo Migrations failed!
    exit /b 1
)
echo Migrations completed successfully!
goto end

:health
echo Checking service health...
powershell -Command "try { $response = Invoke-WebRequest -Uri 'http://localhost:8081/health' -UseBasicParsing -TimeoutSec 5; if ($response.StatusCode -eq 200) { Write-Host 'API Server is healthy' -ForegroundColor Green } else { Write-Host 'API Server returned status:' $response.StatusCode -ForegroundColor Red } } catch { Write-Host 'API Server is not responding' -ForegroundColor Red }"
echo.
echo Docker services status:
docker-compose ps
goto end

:end
endlocal