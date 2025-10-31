# AI-BPMS Backend Development Setup

## 🎉 Initial Setup Complete!

The AI-BPMS backend development environment has been successfully initialized with the following components:

### 📁 Project Structure
```
ai-bpms-backend/
├── cmd/                          # Application entry points
│   ├── server/main.go           # Main API server
│   └── migrate/main.go          # Database migration tool
├── core-services/               # Performance-critical services
│   ├── process-engine/          # Process + Workflow + Tasks
│   └── user-management/         # Auth + Users + Permissions
├── domain-services/             # Business-specific services
│   ├── analytics-service/       # Reporting & metrics
│   ├── notification-service/    # Real-time notifications
│   ├── ai-service/             # AI integrations
│   └── form-service/           # Dynamic form schemas
├── infrastructure/              # Shared infrastructure
│   └── api-gateway/            # API routing & middleware
├── shared/                     # Shared libraries
│   ├── common/                 # Common utilities
│   └── database/               # Database utilities
├── configs/                    # Configuration files
├── deployment/                 # Deployment configurations
└── scripts/                    # Development scripts
```

### 🛠️ Technologies Stack
- **Language**: Go 1.21+
- **HTTP Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Message Queue**: NATS
- **Cache**: Redis
- **Authentication**: JWT + Keycloak/OIDC support
- **API Documentation**: Swagger/OpenAPI
- **Configuration**: Viper
- **Logging**: Logrus

### 🚀 Quick Start

1. **Prerequisites**
   - Go 1.21+
   - Docker & Docker Compose
   - Make (optional but recommended)

2. **Setup Environment**
   ```bash
   # Windows
   .\scripts\setup.bat
   
   # Linux/macOS
   ./scripts/setup.sh
   ```

3. **Start Development**
   ```bash
   make dev
   ```

4. **Verify Setup**
   - Visit: http://localhost:8081/health
   - API Docs: http://localhost:8081/swagger/index.html

### 📋 Available Commands

```bash
# Development
make dev              # Start full development environment
make run              # Run the application
make test             # Run tests
make fmt              # Format code
make lint             # Run linter

# Database
make migrate          # Run database migrations
make migrate-rollback # Rollback last migration

# Build & Deploy
make build            # Build application
make docker-build     # Build Docker image
make docker-compose-up # Start services

# Documentation
make swagger          # Generate API docs
make help             # Show all commands
```

### 🐳 Development Services

| Service | URL | Credentials |
|---------|-----|-------------|
| API Server | http://localhost:8081 | - |
| Swagger Docs | http://localhost:8081/swagger/index.html | - |
| PostgreSQL | localhost:5432 | bpms_user/bpms_pass |
| NATS | localhost:4222 | - |
| Redis | localhost:6379 | - |
| Keycloak | http://localhost:8080 | admin/admin |
| Prometheus | http://localhost:9090 | - |
| Grafana | http://localhost:3000 | admin/admin |

### ⚙️ Configuration

1. **Environment Variables**
   - Copy `.env.example` to `.env`
   - Update values for your environment

2. **Configuration Files**
   - `configs/config.yaml` - Development config
   - `configs/config.production.yaml` - Production config

### 🔐 Authentication Setup

The project supports multiple authentication providers:

#### Built-in JWT (Default)
```yaml
auth:
  provider: "jwt"
  jwt:
    secret: "your-secret-key"
```

#### Keycloak Integration
```yaml
auth:
  provider: "keycloak"
  keycloak:
    base_url: "http://localhost:8080"
    realm: "ai-bpms"
    client_id: "ai-bpms-backend"
```

#### External OIDC (Auth0, etc.)
```yaml
auth:
  provider: "oidc"
  oidc:
    client_id: "your-client-id"
    domain: "your-domain.auth0.com"
```

### 📊 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/api/v1/auth/login` | POST | User login |
| `/api/v1/processes` | GET | List processes |
| `/api/v1/instances` | POST | Start process instance |
| `/api/v1/tasks` | GET | List user tasks |
| `/swagger/*` | GET | API documentation |

### 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./shared/common/...
```

### 🐛 Development Tools

- **Hot Reload**: `make dev-watch` (requires air)
- **Linting**: `make lint` (golangci-lint)
- **Security Scan**: `make security` (gosec)
- **Code Format**: `make fmt` (gofmt)

### 📝 Next Steps

1. **Implement Core Services**
   - Authentication handlers
   - Process engine logic
   - Task management
   - Business rules engine

2. **Add AI Integration**
   - OpenAI API integration
   - Process generation
   - Rule optimization

3. **Enhance Security**
   - RBAC implementation
   - Rate limiting
   - Audit logging

4. **Performance Optimization**
   - Caching strategies
   - Database optimization
   - Load testing

### 🤝 Development Workflow

1. **Feature Development**
   ```bash
   git checkout -b feature/your-feature
   make dev              # Start environment
   # Make changes
   make test             # Run tests
   make lint             # Check code quality
   git commit -m "feat: your feature"
   ```

2. **Database Changes**
   ```bash
   # Create migration
   make migrate
   # Test rollback
   make migrate-rollback
   ```

3. **API Documentation**
   ```bash
   # Update swagger comments
   make swagger
   # View docs at http://localhost:8081/swagger/index.html
   ```

### 🆘 Troubleshooting

#### Common Issues:

1. **Port conflicts**: Check if ports 8080, 5432, 4222, 6379 are free
2. **Docker issues**: Run `docker-compose down` and `docker-compose up -d`
3. **Database connection**: Verify PostgreSQL is running and accessible
4. **Build errors**: Run `go mod tidy` and `make deps`

#### Logs:
```bash
# Application logs
make logs

# Docker services
docker-compose logs -f

# Check service health
make health-check
```

---

🎯 **The development environment is ready for core AI-BPMS development!**

Next: Start implementing authentication, process engine, and AI integration features.