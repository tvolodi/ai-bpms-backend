# Windows Development Guide

## 🪟 Windows-Specific Setup

Since you're on Windows, you have multiple options for running the development environment:

### Option 1: PowerShell Script (Recommended)
```powershell
# Start full development environment
.\dev.ps1 dev

# Other commands
.\dev.ps1 build      # Build application
.\dev.ps1 run        # Run application  
.\dev.ps1 test       # Run tests
.\dev.ps1 health     # Check services
.\dev.ps1 help       # Show all commands
```

### Option 2: Batch File (Simple)
```cmd
# Start full development environment
dev.bat dev

# Other commands
dev.bat build        # Build application
dev.bat run          # Run application
dev.bat test         # Run tests
dev.bat health       # Check services
```

### Option 3: Manual Commands
If scripts don't work, you can run commands manually:

```powershell
# 1. Start services
docker-compose up -d

# 2. Wait for services (10 seconds)
Start-Sleep 10

# 3. Build and run migrations
go build -o bin/migrate.exe ./cmd/migrate/main.go
./bin/migrate.exe

# 4. Build and run application
go build -o bin/ai-bpms-backend.exe ./cmd/server/main.go
./bin/ai-bpms-backend.exe
```

## 🔧 Windows-Specific Fixes Applied

1. **Fixed Makefile**: Removed Unix-specific commands and color codes
2. **Created PowerShell script**: `dev.ps1` with full functionality
3. **Created Batch file**: `dev.bat` for simple commands
4. **Windows sleep**: Uses `ping` or `timeout` instead of `sleep`
5. **Directory creation**: Uses Windows-compatible commands

## 🚀 Quick Start for Windows

1. **Open PowerShell as Administrator** (recommended)
2. **Navigate to project directory**:
   ```powershell
   cd C:\Users\tvolo\dev\ai-bpms\backend
   ```

3. **Start development environment**:
   ```powershell
   .\dev.ps1 dev
   ```

4. **Verify services are running**:
   ```powershell
   .\dev.ps1 health
   ```

5. **Access the application**:
   - API: http://localhost:8081
   - Health: http://localhost:8081/health
   - Swagger: http://localhost:8081/swagger/index.html

## 🐛 Troubleshooting Windows Issues

### PowerShell Execution Policy
If you get execution policy errors:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Docker Issues
If Docker commands fail:
1. Make sure Docker Desktop is running
2. Check if you're in the correct directory
3. Try running Docker commands manually:
   ```powershell
   docker-compose --version
   docker-compose up -d
   ```

### Port Conflicts
If you get port conflicts:
- Check what's using port 8081: `netstat -ano | findstr :8081`
- Stop the process or change the port in config files

### Build Errors
If Go build fails:
```powershell
go mod tidy
go mod download
```

## 🎯 Available Commands

### PowerShell Script (`.\dev.ps1`)
| Command | Description |
|---------|-------------|
| `dev` | Start full development environment |
| `build` | Build the application |
| `run` | Build and run the application |
| `test` | Run all tests |
| `fmt` | Format Go code |
| `clean` | Clean build artifacts |
| `services-up` | Start Docker services only |
| `services-down` | Stop Docker services |
| `migrate` | Run database migrations |
| `health` | Check service health |
| `help` | Show help message |

### Batch File (`dev.bat`)
Same commands as PowerShell script but with simpler implementation.

### Traditional Make (if you have it)
If you have `make` installed (via Git Bash, MSYS2, or WSL):
```bash
make dev
make build
make test
```

## 🔄 Development Workflow

1. **Start environment**: `.\dev.ps1 dev`
2. **Make code changes**
3. **Test changes**: `.\dev.ps1 test`
4. **Run application**: `.\dev.ps1 run`
5. **Check health**: `.\dev.ps1 health`

The environment is now Windows-compatible and should work smoothly on your system!