# AI-BPMS Backend

## Overview
Core backend service for the AI-powered Business Process Management System. Built with Go-lang for high performance, scalability, and AI code generation simplicity.

## Technology Stack
- **Language**: Go 1.21+
- **HTTP Framework**: Gin (lightweight, fast routing)
- **Database**: PostgreSQL with GORM ORM
- **Message Queue**: NATS for process coordination and events
- **Rule Engine**: expr library for business rule evaluation
- **Authentication**: JWT tokens with golang-jwt/jwt
- **AI Integration**: Standard HTTP client for AI service APIs
- **WebSocket**: Gorilla WebSocket for real-time updates
- **Configuration**: Viper for configuration management
- **Logging**: Logrus or Zap for structured logging
- **API Documentation**: Swag for Swagger/OpenAPI generation
- **Testing**: Testify for unit/integration tests
- **Monitoring**: Prometheus metrics + health checks

## Core Features
- Process definition and management
- Workflow execution engine
- Business rule evaluation and AI-enhanced rules
- User and role management
- Process monitoring and analytics
- REST API for client applications
- Real-time notifications (WebSocket)

## Performance-Optimized Architecture

**Decision**: Hybrid microservices with smart service boundaries to balance AI-friendly patterns with enterprise performance requirements.

```
ai-bpms-backend/
├── core-services/         # Performance-critical, shared data domains
│   ├── process-engine/    # Combined: Process + Workflow + Tasks
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── handlers/  # HTTP handlers
│   │   │   ├── services/  # Business logic
│   │   │   ├── repositories/ # Data access with optimized queries
│   │   │   ├── models/    # Domain models
│   │   │   └── cache/     # Multi-level caching
│   │   ├── pkg/api/       # Public interfaces
│   │   └── configs/
│   └── user-management/   # Combined: Auth + Users + Permissions
│       ├── cmd/main.go
│       ├── internal/
│       │   ├── handlers/
│       │   ├── services/
│       │   ├── repositories/
│       │   ├── auth/      # Authentication providers
│       │   │   ├── jwt/   # Built-in JWT
│       │   │   ├── oidc/  # External OIDC (Keycloak/Auth0)
│       │   │   └── providers/ # Auth provider integrations
│       │   └── rbac/      # Role-based access control
│       └── configs/
├── domain-services/       # Business-specific, can be separate
│   ├── analytics-service/ # Reporting & metrics (read-optimized)
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── aggregators/ # Data aggregation
│   │   │   ├── collectors/  # Metrics collection
│   │   │   └── exporters/   # Export to BI tools
│   │   └── configs/
│   ├── notification-service/ # Real-time notifications
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── channels/    # Notification channels
│   │   │   ├── templates/   # Message templates
│   │   │   └── delivery/    # Delivery mechanisms
│   │   └── configs/
│   ├── ai-service/        # AI integrations & processing
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── clients/     # AI service clients
│   │   │   ├── processors/  # AI response processing
│   │   │   ├── embeddings/  # Vector embeddings
│   │   │   └── training/    # Model training data
│   │   └── configs/
│   └── form-service/      # Dynamic form schemas
│       ├── cmd/main.go
│       ├── internal/
│       │   ├── generators/  # Schema generators
│       │   ├── validators/  # Form validation
│       │   └── renderers/   # Form rendering
│       └── configs/
├── infrastructure/        # Shared infrastructure
│   ├── api-gateway/       # Single entry point + routing
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── router/      # Request routing
│   │   │   ├── middleware/  # Cross-cutting concerns
│   │   │   ├── cache/       # Response caching
│   │   │   ├── ratelimit/   # Rate limiting
│   │   │   └── monitoring/  # Request monitoring
│   │   └── configs/
│   ├── event-bus/         # NATS messaging system
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── publishers/  # Event publishers
│   │   │   ├── subscribers/ # Event subscribers
│   │   │   └── routing/     # Message routing
│   │   └── configs/
│   └── monitoring/        # Observability stack
│       ├── metrics/       # Prometheus metrics
│       ├── logging/       # Centralized logging
│       ├── tracing/       # Distributed tracing
│       └── alerting/      # Alert management
├── shared/               # Shared across all services
│   ├── database/         # Shared database for core entities
│   │   ├── migrations/   # Database migrations
│   │   ├── views/        # Optimized database views
│   │   ├── indexes/      # Performance indexes
│   │   └── pools/        # Connection pool management
│   ├── cache/            # Distributed caching (Redis)
│   │   ├── strategies/   # Caching strategies
│   │   ├── invalidation/ # Cache invalidation
│   │   └── serializers/  # Data serialization
│   ├── common/           # Shared libraries
│   │   ├── types/        # Common types
│   │   ├── utils/        # Utility functions
│   │   ├── errors/       # Error handling
│   │   ├── validation/   # Data validation
│   │   └── middleware/   # Shared middleware
│   └── security/         # Security utilities
│       ├── encryption/   # Data encryption
│       ├── signing/      # Digital signatures
│       └── audit/        # Audit logging
├── performance/          # Performance optimization
│   ├── load-testing/     # Load test scripts
│   ├── benchmarks/       # Performance benchmarks
│   ├── profiling/        # Performance profiling
│   └── optimization/     # Query optimization tools
└── deployment/           # Deployment configurations
    ├── docker/           # Docker configurations
    ├── kubernetes/       # K8s manifests
    ├── helm/            # Helm charts
    └── terraform/       # Infrastructure as code
```

### **Key Architectural Decisions**:

1. **Core Services**: Combine tightly coupled domains (Process+Workflow+Tasks, Auth+Users+Permissions) to eliminate inter-service latency
2. **Shared Database**: Use shared database for core entities to enable complex queries and transactions
3. **Smart Caching**: Multi-level caching (local + Redis + database query cache)
4. **Async Processing**: Event-driven architecture for non-critical operations
5. **Performance Monitoring**: Built-in observability and automated optimization

## Key Components
- **Process Engine**: Orchestrates workflow execution with goroutine-based parallel processing
- **Rule Engine**: Evaluates business rules and AI-enhanced logic using expr library
- **API Gateway**: RESTful API with Gin framework and middleware pipeline
- **Event System**: NATS-based event streaming for process coordination
- **AI Coordinator**: HTTP client integration with external AI services (OpenAI, etc.)
- **WebSocket Hub**: Real-time notifications and process updates
- **Authentication**: JWT-based auth with RBAC (Role-Based Access Control)
- **Task Scheduler**: Background job processing and task execution
- **Form Schema Generator**: Dynamic JSON schema generation for RJSF integration

## Authentication & Authorization Strategy

### **Keycloak Integration (Primary Authentication)**

**Why Keycloak for AI-BPMS:**
- ✅ **Enterprise Security**: Production-ready OAuth2/OIDC implementation
- ✅ **Zero Licensing Cost**: Completely free, open-source
- ✅ **Self-Hosted Control**: Full data sovereignty and customization
- ✅ **AI-Friendly Patterns**: Standard OAuth2/OIDC that AI models know well
- ✅ **BPMS Features**: Perfect for role-based process management
- ✅ **Scalable**: Handles thousands of users efficiently

### **Architecture Overview:**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React Client  │───▶│   Go Backend    │───▶│    Keycloak     │
│                 │    │                 │    │                 │
│ • Login Flow    │    │ • Token Verify  │    │ • User Store    │
│ • Token Storage │    │ • RBAC Check    │    │ • Role Mgmt     │
│ • Auto Refresh  │    │ • API Gateway   │    │ • SSO/LDAP      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Go Backend Integration:**
```go
// Keycloak OIDC configuration
type KeycloakConfig struct {
    BaseURL      string `yaml:"base_url"`      // http://localhost:8080
    Realm        string `yaml:"realm"`         // ai-bpms
    ClientID     string `yaml:"client_id"`     // ai-bpms-backend
    ClientSecret string `yaml:"client_secret"` // your-client-secret
}

// OIDC token validation
type KeycloakService struct {
    config    *KeycloakConfig
    verifier  *oidc.IDTokenVerifier
    provider  *oidc.Provider
    httpClient *http.Client
}

func NewKeycloakService(config *KeycloakConfig) (*KeycloakService, error) {
    ctx := context.Background()
    
    // Initialize OIDC provider
    providerURL := fmt.Sprintf("%s/realms/%s", config.BaseURL, config.Realm)
    provider, err := oidc.NewProvider(ctx, providerURL)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider: %v", err)
    }
    
    // Configure token verifier
    verifier := provider.Verifier(&oidc.Config{
        ClientID: config.ClientID,
    })
    
    return &KeycloakService{
        config:   config,
        provider: provider,
        verifier: verifier,
    }, nil
}

// Validate JWT token from Keycloak
func (k *KeycloakService) ValidateToken(tokenString string) (*User, error) {
    ctx := context.Background()
    
    // Verify JWT signature and claims
    idToken, err := k.verifier.Verify(ctx, tokenString)
    if err != nil {
        return nil, fmt.Errorf("failed to verify token: %v", err)
    }
    
    // Extract user claims
    var claims struct {
        Sub           string   `json:"sub"`
        Email         string   `json:"email"`
        EmailVerified bool     `json:"email_verified"`
        Name          string   `json:"name"`
        FamilyName    string   `json:"family_name"`
        GivenName     string   `json:"given_name"`
        RealmRoles    []string `json:"realm_access.roles"`
        ResourceRoles map[string]struct {
            Roles []string `json:"roles"`
        } `json:"resource_access"`
    }
    
    if err := idToken.Claims(&claims); err != nil {
        return nil, fmt.Errorf("failed to parse claims: %v", err)
    }
    
    // Map to internal user structure
    user := &User{
        ID:        claims.Sub,
        Email:     claims.Email,
        FirstName: claims.GivenName,
        LastName:  claims.FamilyName,
        IsActive:  claims.EmailVerified,
        Roles:     extractBPMSRoles(claims.RealmRoles, claims.ResourceRoles),
    }
    
    return user, nil
}
```

### **Authentication Middleware:**
```go
// Keycloak authentication middleware
func KeycloakAuthMiddleware(kcService *KeycloakService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Extract Bearer token
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            c.JSON(401, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }
        
        // Validate with Keycloak
        user, err := kcService.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
            c.Abort()
            return
        }
        
        // Store user in context
        c.Set("user", user)
        c.Set("auth_method", "keycloak")
        c.Next()
    }
}

// Role-based authorization middleware
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(401, gin.H{"error": "User not authenticated"})
            c.Abort()
            return
        }
        
        userObj := user.(*User)
        if !userObj.HasAnyRole(roles...) {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```
### **BPMS Role Mapping:**
```go
// Map Keycloak roles to BPMS permissions
func extractBPMSRoles(realmRoles []string, resourceRoles map[string]struct{ Roles []string }) []string {
    var bpmsRoles []string
    
    // Map realm roles
    for _, role := range realmRoles {
        switch role {
        case "bpms-admin":
            bpmsRoles = append(bpmsRoles, "admin")
        case "bpms-manager":
            bpmsRoles = append(bpmsRoles, "manager")
        case "bpms-user":
            bpmsRoles = append(bpmsRoles, "user")
        }
    }
    
    // Map client-specific roles
    if clientRoles, exists := resourceRoles["ai-bpms-backend"]; exists {
        for _, role := range clientRoles.Roles {
            switch role {
            case "process-designer":
                bpmsRoles = append(bpmsRoles, "process-designer")
            case "task-assignee":
                bpmsRoles = append(bpmsRoles, "task-assignee")
            case "analytics-viewer":
                bpmsRoles = append(bpmsRoles, "analytics-viewer")
            }
        }
    }
    
    return bpmsRoles
}

// BPMS Permission system
type BPMSPermission string

const (
    // Process permissions
    ProcessCreate BPMSPermission = "process:create"
    ProcessRead   BPMSPermission = "process:read"
    ProcessUpdate BPMSPermission = "process:update"
    ProcessDelete BPMSPermission = "process:delete"
    ProcessStart  BPMSPermission = "process:start"
    
    // Task permissions
    TaskAssign   BPMSPermission = "task:assign"
    TaskComplete BPMSPermission = "task:complete"
    TaskView     BPMSPermission = "task:view"
    TaskDelegate BPMSPermission = "task:delegate"
    
    // Analytics permissions
    AnalyticsView   BPMSPermission = "analytics:view"
    AnalyticsExport BPMSPermission = "analytics:export"
    
    // Admin permissions
    UserManage     BPMSPermission = "user:manage"
    SystemConfig   BPMSPermission = "system:config"
    RoleManage     BPMSPermission = "role:manage"
)

// Role to permission mapping
var RolePermissions = map[string][]BPMSPermission{
    "admin": {
        ProcessCreate, ProcessRead, ProcessUpdate, ProcessDelete, ProcessStart,
        TaskAssign, TaskComplete, TaskView, TaskDelegate,
        AnalyticsView, AnalyticsExport,
        UserManage, SystemConfig, RoleManage,
    },
    "manager": {
        ProcessCreate, ProcessRead, ProcessUpdate, ProcessStart,
        TaskAssign, TaskComplete, TaskView, TaskDelegate,
        AnalyticsView, AnalyticsExport,
    },
    "process-designer": {
        ProcessCreate, ProcessRead, ProcessUpdate,
        AnalyticsView,
    },
    "user": {
        ProcessRead, ProcessStart,
        TaskComplete, TaskView,
    },
    "analytics-viewer": {
        ProcessRead, AnalyticsView,
    },
}
```

### **Keycloak Setup Configuration:**
```yaml
# Keycloak realm configuration for AI-BPMS
keycloak:
  base_url: "http://localhost:8080"
  realm: "ai-bpms"
  
  # Backend service client
  backend_client:
    client_id: "ai-bpms-backend"
    client_secret: "your-backend-secret"
    
  # Frontend client (public)
  frontend_client:
    client_id: "ai-bpms-frontend"
    public: true
    redirect_uris:
      - "http://localhost:3000/*"
      - "https://app.ai-bpms.com/*"
    
  # Realm roles
  realm_roles:
    - name: "bpms-admin"
      description: "Full BPMS administration access"
    - name: "bpms-manager" 
      description: "Department-level process management"
    - name: "bpms-user"
      description: "Basic process participation"
      
  # Client roles for fine-grained permissions
  client_roles:
    ai-bpms-backend:
      - name: "process-designer"
        description: "Can create and modify process definitions"
      - name: "task-assignee"
        description: "Can assign tasks to other users"
      - name: "analytics-viewer"
        description: "Can view process analytics and reports"
```

### **Development Setup:**
```bash
# Start Keycloak with Docker
docker run -d \
  --name keycloak \
  -p 8080:8080 \
  -e KEYCLOAK_ADMIN=admin \
  -e KEYCLOAK_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:latest \
  start-dev

# Access Keycloak admin console
# http://localhost:8080/admin
# Username: admin, Password: admin

# Go dependencies for OIDC
go get github.com/coreos/go-oidc/v3/oidc
go get golang.org/x/oauth2
```

### **React Client Integration:**
```typescript
// Keycloak configuration for React client
import Keycloak from 'keycloak-js';

const keycloakConfig = {
  url: 'http://localhost:8080/',
  realm: 'ai-bpms',
  clientId: 'ai-bpms-frontend'
};

const keycloak = new Keycloak(keycloakConfig);

// Initialize Keycloak
export const initKeycloak = () => {
  return keycloak.init({
    onLoad: 'login-required',
    checkLoginIframe: false,
    pkceMethod: 'S256'
  });
};

// Get token for API calls
export const getToken = () => {
  return keycloak.token;
};

// Auto-refresh token
export const updateToken = () => {
  return keycloak.updateToken(30);
};
```
### **API Endpoints with Keycloak:**
```go
// Protected API routes with role-based access
func setupRoutes(router *gin.Engine, kcService *KeycloakService) {
    api := router.Group("/api/v1")
    api.Use(KeycloakAuthMiddleware(kcService))
    
    // Process management (requires manager or admin role)
    processes := api.Group("/processes")
    processes.Use(RequireRole("manager", "admin", "process-designer"))
    {
        processes.GET("", listProcesses)
        processes.POST("", RequireRole("manager", "admin"), createProcess)
        processes.PUT("/:id", RequireRole("manager", "admin"), updateProcess)
        processes.DELETE("/:id", RequireRole("admin"), deleteProcess)
    }
    
    // Task management (all authenticated users)
    tasks := api.Group("/tasks")
    {
        tasks.GET("", listUserTasks)           // All users can view their tasks
        tasks.POST("/:id/complete", completeTask) // Complete assigned tasks
        tasks.POST("/:id/assign", RequireRole("manager", "admin", "task-assignee"), assignTask)
    }
    
    // Analytics (requires analytics role)
    analytics := api.Group("/analytics")
    analytics.Use(RequireRole("analytics-viewer", "manager", "admin"))
    {
        analytics.GET("/dashboard", getDashboard)
        analytics.GET("/reports", getReports)
        analytics.GET("/export", RequireRole("manager", "admin"), exportData)
    }
    
    // Admin endpoints (admin only)
    admin := api.Group("/admin")
    admin.Use(RequireRole("admin"))
    {
        admin.GET("/users", listUsers)
        admin.POST("/users", createUser)
        admin.PUT("/users/:id/roles", updateUserRoles)
    }
}
```

### **Configuration:**
```yaml
# config.yaml
server:
  port: 8080
  host: "localhost"

# Keycloak integration
auth:
  provider: "keycloak"
  keycloak:
    base_url: "http://localhost:8080"
    realm: "ai-bpms"
    client_id: "ai-bpms-backend"
    client_secret: "your-client-secret"
    
  # Token settings
  token:
    verification_cache_ttl: "5m"    # Cache verified tokens
    jwks_cache_ttl: "1h"           # Cache Keycloak public keys
    
  # Role mapping
  role_mapping:
    bpms-admin: ["admin"]
    bpms-manager: ["manager"]
    bpms-user: ["user"]

database:
  host: "localhost"
  port: 5432
  user: "bpms_user"
  password: "bpms_pass"
  dbname: "ai_bpms"
```

### **Keycloak Benefits for BPMS:**
- ✅ **Enterprise Ready**: Production-grade security out of the box
- ✅ **Role Management**: Perfect for BPMS hierarchical permissions
- ✅ **SSO Integration**: LDAP, SAML, OAuth providers
- ✅ **User Self-Service**: Password reset, profile management
- ✅ **Audit Trails**: Complete authentication logging
- ✅ **Admin UI**: Easy user and role management
- ✅ **API Friendly**: RESTful admin APIs for automation
- ✅ **Clustering**: High availability support
- ✅ **Themes**: Customizable login pages
- ✅ **Free & Open Source**: No licensing costs
    
    OIDC struct {
        Provider     string `yaml:"provider"`     // auth0, keycloak, etc.
        ClientID     string `yaml:"client_id"`
        ClientSecret string `yaml:"client_secret"`
        Domain       string `yaml:"domain"`
        Audience     string `yaml:"audience"`
    } `yaml:"oidc"`
    
    BuiltIn struct {
        Enabled          bool `yaml:"enabled"`
        RequireOTP       bool `yaml:"require_otp"`
        PasswordlessOnly bool `yaml:"passwordless_only"`
    } `yaml:"built_in"`
}
```

### **Why External Provider First:**
- ✅ **Enterprise Security**: Proven security implementation
- ✅ **Compliance**: SOC2, ISO27001, GDPR ready
- ✅ **Zero Password Risk**: OAuth2/OIDC eliminates password storage
- ✅ **Advanced Features**: SSO, MFA, breach detection built-in
- ✅ **Maintenance**: Security updates handled by provider
- ✅ **AI-Friendly**: Standard OAuth2 patterns well-known to AI

### **Phase 2: Enhanced Built-in (Development/Self-Hosted)**
For development or self-hosted scenarios, enhanced built-in auth:

```go
// Security-hardened built-in authentication
type SecureAuthService struct {
    passwordService  *PasswordService
    otpService      *OTPService
    magicLinkService *MagicLinkService
    breachChecker   *BreachChecker
    rateLimiter     *RateLimiter
}

// Multiple authentication flows
type AuthFlow struct {
    Type        AuthFlowType `json:"type"`
    Identifier  string       `json:"identifier"` // email
    Challenge   string       `json:"challenge"`  // password, otp, magic_link
    MFARequired bool         `json:"mfa_required"`
}

type AuthFlowType string
const (
    FlowPassword    AuthFlowType = "password"
    FlowPasswordless AuthFlowType = "passwordless"
    FlowMagicLink   AuthFlowType = "magic_link"
    FlowOTP         AuthFlowType = "otp"
)
```

### **Security Features (Built-in Enhancement):**
```go
// Password security
type PasswordPolicy struct {
    MinLength        int    `yaml:"min_length"`        // 12+ characters
    RequireEntropy   int    `yaml:"require_entropy"`   // Entropy-based strength
    CheckBreaches    bool   `yaml:"check_breaches"`    // HaveIBeenPwned API
    BlockCommon      bool   `yaml:"block_common"`      // Block common passwords
    RequireMFA       bool   `yaml:"require_mfa"`       // Mandatory 2FA
}

// Multi-factor authentication
type MFAService struct {
    TOTPProvider     TOTPProvider     // Time-based OTP
    SMSProvider      SMSProvider      // SMS codes
    EmailProvider    EmailProvider    // Email codes
    BackupCodes      BackupCodeService // Recovery codes
}

// Passwordless authentication
type PasswordlessService struct {
    MagicLinkTTL    time.Duration // 15 minutes
    OTPCodeLength   int           // 6 digits
    OTPCodeTTL      time.Duration // 5 minutes
    MaxAttempts     int           // 3 attempts
}
```
### **Security Architecture:**
```go
// Unified authentication middleware supporting multiple providers
func UnifiedAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        token := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Try external provider first (more secure)
        if user, err := validateOIDCToken(token); err == nil {
            c.Set("user", user)
            c.Set("auth_method", "oidc")
            c.Next()
            return
        }
        
        // Fallback to built-in JWT (if enabled)
        if config.Auth.BuiltIn.Enabled {
            if user, err := validateJWTToken(token); err == nil {
                c.Set("user", user)
                c.Set("auth_method", "jwt")
                c.Next()
                return
            }
        }
        
        c.JSON(401, gin.H{"error": "Invalid token"})
        c.Abort()
    }
}
```

### **Production Deployment Recommendations:**

#### **For Enterprise/Production:**
1. **Auth0**: $23/month per user - excellent DX, AI-friendly
2. **Keycloak**: Free, self-hosted - enterprise features, more complex
3. **Azure AD B2C**: Microsoft ecosystem integration
4. **Okta**: Enterprise-grade, higher cost

#### **For Development/Self-Hosted:**
1. **Enhanced Built-in**: Passwordless + MFA + breach detection
2. **Keycloak**: Free alternative for full control
3. **Supabase Auth**: Postgres-integrated auth service

### **Risk Mitigation Strategy:**
```yaml
# Production configuration example
auth:
  primary: "auth0"  # or "keycloak"
  fallback_enabled: false  # Disable built-in for production
  
  auth0:
    domain: "your-tenant.auth0.com"
    client_id: "your-client-id"
    client_secret: "your-client-secret"
    audience: "https://api.ai-bpms.com"
    
  security:
    require_mfa: true          # Mandatory for admin/manager
    session_timeout: "4h"      # Auto-logout
    max_concurrent_sessions: 3  # Limit active sessions
    ip_allowlist: []           # Optional IP restrictions
    
  audit:
    log_all_attempts: true     # Complete audit trail
    alert_failed_logins: 5     # Alert after 5 failures
    retention_days: 365        # Keep logs for compliance
```
### **Security Assessment:**

#### **External Provider (Auth0/Keycloak) - RECOMMENDED:**
- ✅ **Zero Password Storage**: No password breach risk
- ✅ **Professional Security**: Dedicated security teams
- ✅ **Compliance Ready**: SOC2, ISO27001, GDPR certified
- ✅ **Advanced Threats**: Bot detection, anomaly detection
- ✅ **Breach Monitoring**: Automatic credential monitoring
- ✅ **AI-Friendly**: Standard OAuth2/OIDC patterns
- ⚠️ **Dependency**: External service dependency
- ⚠️ **Cost**: Per-user pricing for SaaS options

#### **Enhanced Built-in - DEVELOPMENT/SELF-HOSTED:**
- ✅ **Full Control**: Complete customization
- ✅ **No External Deps**: Self-contained
- ✅ **AI-Optimized**: Simple patterns for code generation
- ✅ **Passwordless Options**: Magic links, OTP codes
- ⚠️ **Security Responsibility**: You own all security aspects
- ⚠️ **Compliance Burden**: Self-certification required
- ⚠️ **Attack Surface**: More code to secure

### **Final Recommendation:**

**For Production/Enterprise**: Start with **Auth0** or **Keycloak**
**For Development**: Enhanced built-in with passwordless options
**For AI Code Generation**: Both approaches use standard, well-documented patterns

The security benefits of external providers significantly outweigh the additional complexity, especially for business-critical BPMS systems handling sensitive process data.
    Email         string    `json:"email" gorm:"uniqueIndex"`
    PasswordHash  string    `json:"-"` // Never serialize
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    Role          string    `json:"role"` // admin, manager, user
    IsActive      bool      `json:"is_active"`
    LastLogin     time.Time `json:"last_login"`
    
    // BPMS-specific permissions
    ProcessGroups []string  `json:"process_groups" gorm:"type:text[]"`
    TaskFilters   []string  `json:"task_filters" gorm:"type:text[]"`
    Permissions   []string  `json:"permissions" gorm:"type:text[]"`
}

type AuthService interface {
    Login(email, password string) (*TokenPair, error)
    RefreshToken(refreshToken string) (*TokenPair, error)
    ValidateToken(token string) (*User, error)
    Logout(token string) error
    Register(user CreateUserRequest) (*User, error)
}
```

### **Built-in Features:**
- **JWT Tokens**: Access + Refresh token rotation
- **Password Security**: bcrypt hashing with salt
- **Role-Based Access**: Admin, Manager, User roles
- **Process Permissions**: Fine-grained process access control
- **Session Management**: Token blacklisting and expiration
- **Audit Logging**: Complete authentication audit trail

### **Phase 2: External Provider Integration (Enterprise)**
```go
// OIDC integration for enterprise customers
type OIDCProvider struct {
    Name         string `json:"name"`
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"-"`
    DiscoveryURL string `json:"discovery_url"`
    RedirectURL  string `json:"redirect_url"`
}

// Unified authentication middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        
        if strings.HasPrefix(token, "eyJ") {
            // Built-in JWT validation
            user, err := validateJWTToken(token)
            if err != nil {
                c.JSON(401, gin.H{"error": "Invalid token"})
                c.Abort()
                return
            }
            c.Set("user", user)
        } else {
            // External provider validation (OIDC/OAuth2)
            user, err := validateExternalToken(token)
            if err != nil {
                c.JSON(401, gin.H{"error": "Invalid external token"})
                c.Abort()
                return
            }
            c.Set("user", user)
        }
        
        c.Next()
    }
}
```

### **External Provider Support:**
- **OpenID Connect**: Google, Microsoft, Okta, Auth0
- **SAML 2.0**: Enterprise SSO integration
- **LDAP/Active Directory**: Corporate directory integration
- **OAuth2**: GitHub, GitLab, custom providers

### **RBAC (Role-Based Access Control):**
```go
// Permission system
type Permission string

const (
    // Process permissions
    ProcessCreate Permission = "process:create"
    ProcessRead   Permission = "process:read"
    ProcessUpdate Permission = "process:update"
    ProcessDelete Permission = "process:delete"
    ProcessStart  Permission = "process:start"
    
    // Task permissions
    TaskAssign   Permission = "task:assign"
    TaskComplete Permission = "task:complete"
    TaskView     Permission = "task:view"
    
    // Admin permissions
    UserManage   Permission = "user:manage"
    SystemConfig Permission = "system:config"
)

// Role definitions
var RolePermissions = map[string][]Permission{
    "admin": {
        ProcessCreate, ProcessRead, ProcessUpdate, ProcessDelete,
        TaskAssign, TaskComplete, TaskView,
        UserManage, SystemConfig,
    },
    "manager": {
        ProcessCreate, ProcessRead, ProcessUpdate, ProcessStart,
        TaskAssign, TaskComplete, TaskView,
    },
    "user": {
        ProcessRead, ProcessStart, TaskComplete, TaskView,
    },
}
```

### **Security Features:**
- **Token Security**: Short-lived access tokens (15min) + long-lived refresh (7 days)
- **Rate Limiting**: Login attempt throttling and DDoS protection
- **Password Policy**: Configurable complexity requirements
- **2FA Support**: TOTP-based two-factor authentication (optional)
- **Account Lockout**: Automatic lockout after failed attempts
- **Audit Trail**: Complete authentication and authorization logging

## Configuration
```yaml
# Authentication configuration
auth:
  jwt:
    secret: "your-256-bit-secret"
    access_token_duration: "15m"
    refresh_token_duration: "168h" # 7 days
  
  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_symbols: true
  
  external_providers:
    oidc:
      enabled: false
      providers:
        - name: "google"
          client_id: "your-google-client-id"
          client_secret: "your-google-client-secret"
          discovery_url: "https://accounts.google.com/.well-known/openid_configuration"
        - name: "microsoft"
          client_id: "your-azure-client-id"
          client_secret: "your-azure-client-secret"
          discovery_url: "https://login.microsoftonline.com/common/v2.0/.well-known/openid_configuration"
```
```bash
# Prerequisites
go version  # Requires Go 1.21+
docker --version  # For PostgreSQL and NATS
docker-compose --version

# Clone and setup
git clone <repository>
cd ai-bpms/backend

# Install dependencies
go mod init github.com/your-org/ai-bpms-backend
go mod tidy

# Core dependencies
go get github.com/gin-gonic/gin
go get gorm.io/gorm gorm.io/driver/postgres
go get github.com/nats-io/nats.go
go get github.com/expr-lang/expr
go get github.com/golang-jwt/jwt/v5
go get github.com/gorilla/websocket
go get github.com/spf13/viper
go get github.com/sirupsen/logrus
go get github.com/swaggo/swag/cmd/swag
go get github.com/stretchr/testify

# Development tools
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Start development environment
docker-compose up -d  # PostgreSQL + NATS + Redis
go run cmd/migrate/main.go  # Run database migrations
go run cmd/server/main.go   # Start API server on :8080
```

## Configuration
```yaml
# configs/config.yaml
server:
  port: 8080
  host: "localhost"
  timeout: 30s

database:
  host: "localhost"
  port: 5432
  user: "bpms_user"
  password: "bpms_pass"
  dbname: "ai_bpms"
  sslmode: "disable"

nats:
  url: "nats://localhost:4222"
  cluster_id: "bpms-cluster"

jwt:
  secret: "your-jwt-secret"
  expiration: "24h"

ai:
  openai:
    api_key: "your-openai-key"
    base_url: "https://api.openai.com/v1"
  custom:
    endpoint: "http://localhost:5000"
```

## API Endpoints
```
# Authentication
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
GET    /api/v1/auth/profile

# Process Management
GET    /api/v1/processes          # List process definitions
POST   /api/v1/processes          # Create process definition
GET    /api/v1/processes/:id      # Get process definition
PUT    /api/v1/processes/:id      # Update process definition
DELETE /api/v1/processes/:id      # Delete process definition

# Process Instances
GET    /api/v1/instances          # List process instances
POST   /api/v1/instances          # Start process instance
GET    /api/v1/instances/:id      # Get instance details
PUT    /api/v1/instances/:id      # Update instance
DELETE /api/v1/instances/:id      # Cancel instance

# Tasks
GET    /api/v1/tasks              # List user tasks
GET    /api/v1/tasks/:id          # Get task details
POST   /api/v1/tasks/:id/complete # Complete task
POST   /api/v1/tasks/:id/assign   # Assign task

# Forms & Schemas
GET    /api/v1/forms/schema/:id   # Get JSON schema for task
POST   /api/v1/forms/validate     # Validate form data

# Rules
GET    /api/v1/rules              # List business rules
POST   /api/v1/rules              # Create rule
PUT    /api/v1/rules/:id          # Update rule
POST   /api/v1/rules/evaluate     # Evaluate rule expression

# AI Integration
POST   /api/v1/ai/process         # AI-assisted process generation
POST   /api/v1/ai/rules           # AI-generated business rules
POST   /api/v1/ai/optimize        # Process optimization suggestions

# Analytics
GET    /api/v1/analytics/dashboard     # Dashboard metrics
GET    /api/v1/analytics/processes     # Process performance
GET    /api/v1/analytics/instances     # Instance statistics

# WebSocket
WS     /ws/notifications          # Real-time process updates
WS     /ws/tasks                  # Task notifications
```

## Performance Goals & Optimizations

### **Target Performance Metrics**
- **API Latency**: P95 < 200ms, P99 < 500ms
- **Throughput**: > 100 processes/second, > 1000 concurrent users
- **Database**: Query time avg < 10ms, connection pool < 80% utilization
- **Cache Hit Rate**: > 90% for frequently accessed data
- **Process Engine**: Process start < 100ms, task assignment < 50ms

### **Performance Optimization Strategies**

#### **1. Smart Service Boundaries**
- **Core Services**: Combine tightly coupled domains to eliminate network hops
- **Process Engine**: Process + Workflow + Tasks in single service
- **User Management**: Auth + Users + Permissions in single service
- **Shared Database**: Enable complex queries and ACID transactions

#### **2. Multi-Level Caching**
```go
// L1: Local cache (100μs access)
// L2: Redis distributed cache (1-2ms access)  
// L3: Database with query cache (5-10ms access)

type CacheManager struct {
    localCache  *bigcache.BigCache
    redisCache  *redis.Client
    dbCache     *gorm.DB
}
```

#### **3. Database Optimizations**
- **Optimized Views**: Pre-computed complex queries
- **Batch Operations**: Reduce database round trips
- **Connection Pooling**: Shared pools between related services
- **Query Optimization**: Dynamic query optimization based on context

#### **4. Async Processing**
- **Event-Driven**: Non-critical operations moved to background
- **NATS Streaming**: High-performance message processing
- **Background Workers**: Separate processing for analytics, notifications

#### **5. Performance Monitoring**
- **Real-time Metrics**: Prometheus + Grafana dashboards
- **Automated Alerting**: Performance degradation detection
- **Load Testing**: Continuous performance validation
- **Query Profiling**: Automatic slow query detection and optimization

### **Load Testing Strategy**
- **Normal Load**: 100 concurrent users, 500 RPS
- **Peak Load**: 500 concurrent users, 2000 RPS  
- **Stress Test**: 1000 concurrent users, 5000 RPS
- **Acceptable Latency**: 95% of requests < 200ms under peak load

See `backend_performance_strategy.md` for detailed performance architecture and optimization techniques.

## AI Integration Features
- **Process Generation**: AI-assisted workflow creation from natural language
- **Rule Optimization**: AI-powered business rule suggestions
- **Anomaly Detection**: AI monitoring for process bottlenecks
- **Dynamic Forms**: AI-generated form schemas based on context
- **Decision Support**: AI recommendations for process decisions
- **Natural Language**: Query processes using natural language

## Security & Compliance
- **Authentication**: JWT with refresh token rotation
- **Authorization**: Role-based access control (RBAC)
- **Data Encryption**: TLS 1.3 for transport, AES-256 for data at rest
- **Audit Logging**: Complete audit trail for all operations
- **Rate Limiting**: Request rate limiting and DDoS protection
- **Input Validation**: Comprehensive input sanitization
- **GDPR Compliance**: Data privacy and right to deletion

## Status
🚧 **Planning Phase** - Ready for Go project initialization and core development

## Getting Started
1. **Setup Environment**: Install Go 1.21+, Docker, and PostgreSQL
2. **Initialize Project**: Run `go mod init` and install dependencies
3. **Start Services**: Use `docker-compose up -d` for external services
4. **Run Migrations**: Execute database schema setup
5. **Start Server**: Launch API server with `go run cmd/server/main.go`
6. **API Documentation**: Access Swagger UI at `http://localhost:8080/swagger`

## Development Workflow
- **Clean Architecture**: Domain-driven design with clear separation
- **Test-Driven**: Unit tests with testify framework
- **API-First**: OpenAPI specification for client generation
- **Event-Driven**: NATS messaging for decoupled components
- **AI-Ready**: Extensible AI integration framework
- **Monitoring**: Built-in health checks and metrics