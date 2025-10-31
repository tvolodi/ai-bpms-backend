package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/tvolodi/ai-bpms-backend/shared/common/config"
	"github.com/tvolodi/ai-bpms-backend/shared/common/middleware"
	"github.com/tvolodi/ai-bpms-backend/shared/database"
)

// @title AI-BPMS Backend API
// @version 1.0
// @description AI-powered Business Process Management System Backend
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.ai-bpms.com/support
// @contact.email support@ai-bpms.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	setupLogging(cfg.Logging)

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Setup middleware
	setupMiddleware(router, cfg)

	// Setup routes
	setupRoutes(router, cfg, db)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server starting on %s", srv.Addr)

		if cfg.Server.TLS.Enabled {
			if err := srv.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Server shutting down...")

	// Give server 30 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

func setupLogging(cfg config.LoggingConfig) {
	// Set log level
	switch cfg.Level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if cfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set log output
	if cfg.Output == "file" && cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		logrus.SetOutput(file)
	}
}

func setupMiddleware(router *gin.Engine, cfg *config.Config) {
	// Recovery middleware
	router.Use(gin.Recovery())

	// Request logging middleware
	router.Use(middleware.RequestLogger())

	// CORS middleware
	router.Use(middleware.CORS(cfg.Security.CORS))

	// Rate limiting middleware
	if cfg.Security.RateLimit.Enabled {
		router.Use(middleware.RateLimit(cfg.Security.RateLimit))
	}

	// Security headers middleware
	router.Use(middleware.SecurityHeaders())

	// Request ID middleware
	router.Use(middleware.RequestID())
}

func setupRoutes(router *gin.Engine, cfg *config.Config, db interface{}) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Metrics endpoint (if enabled)
	if cfg.Metrics.Enabled {
		router.GET(cfg.Metrics.Path, metricsHandler)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", loginHandler)
			auth.POST("/logout", logoutHandler)
			auth.POST("/refresh", refreshHandler)
			auth.GET("/profile", profileHandler)
		}

		// Process routes
		processes := v1.Group("/processes")
		// TODO: Add authentication middleware
		{
			processes.GET("", listProcesses)
			processes.POST("", createProcess)
			processes.GET("/:id", getProcess)
			processes.PUT("/:id", updateProcess)
			processes.DELETE("/:id", deleteProcess)
		}

		// Process instance routes
		instances := v1.Group("/instances")
		// TODO: Add authentication middleware
		{
			instances.GET("", listInstances)
			instances.POST("", startInstance)
			instances.GET("/:id", getInstance)
			instances.PUT("/:id", updateInstance)
			instances.DELETE("/:id", cancelInstance)
		}

		// Task routes
		tasks := v1.Group("/tasks")
		// TODO: Add authentication middleware
		{
			tasks.GET("", listTasks)
			tasks.GET("/:id", getTask)
			tasks.POST("/:id/complete", completeTask)
			tasks.POST("/:id/assign", assignTask)
		}

		// Form schema routes
		forms := v1.Group("/forms")
		// TODO: Add authentication middleware
		{
			forms.GET("/schema/:id", getFormSchema)
			forms.POST("/validate", validateForm)
		}

		// Business rules routes
		rules := v1.Group("/rules")
		// TODO: Add authentication middleware
		{
			rules.GET("", listRules)
			rules.POST("", createRule)
			rules.PUT("/:id", updateRule)
			rules.POST("/evaluate", evaluateRule)
		}

		// AI integration routes
		ai := v1.Group("/ai")
		// TODO: Add authentication middleware
		{
			ai.POST("/process", aiGenerateProcess)
			ai.POST("/rules", aiGenerateRules)
			ai.POST("/optimize", aiOptimizeProcess)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		// TODO: Add authentication middleware
		{
			analytics.GET("/dashboard", getDashboard)
			analytics.GET("/processes", getProcessAnalytics)
			analytics.GET("/instances", getInstanceAnalytics)
		}

		// Admin routes
		admin := v1.Group("/admin")
		// TODO: Add authentication and admin role middleware
		{
			admin.GET("/users", listUsers)
			admin.POST("/users", createUser)
			admin.PUT("/users/:id", updateUser)
			admin.DELETE("/users/:id", deleteUser)
			admin.PUT("/users/:id/roles", updateUserRoles)
		}
	}

	// WebSocket routes
	router.GET("/ws/notifications", wsNotifications)
	router.GET("/ws/tasks", wsTasks)
}

// Health check handler
// @Summary Health check
// @Description Get the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "ai-bpms-backend",
		"version":   "1.0.0",
	})
}

// Metrics handler
func metricsHandler(c *gin.Context) {
	// TODO: Implement Prometheus metrics
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics endpoint - TODO: Implement Prometheus metrics",
	})
}

// Placeholder handlers - TODO: Implement actual logic

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - TODO: Implement"})
}

func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout endpoint - TODO: Implement"})
}

func refreshHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Refresh endpoint - TODO: Implement"})
}

func profileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Profile endpoint - TODO: Implement"})
}

func listProcesses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List processes - TODO: Implement"})
}

func createProcess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create process - TODO: Implement"})
}

func getProcess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get process - TODO: Implement"})
}

func updateProcess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update process - TODO: Implement"})
}

func deleteProcess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete process - TODO: Implement"})
}

func listInstances(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List instances - TODO: Implement"})
}

func startInstance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Start instance - TODO: Implement"})
}

func getInstance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get instance - TODO: Implement"})
}

func updateInstance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update instance - TODO: Implement"})
}

func cancelInstance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Cancel instance - TODO: Implement"})
}

func listTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List tasks - TODO: Implement"})
}

func getTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get task - TODO: Implement"})
}

func completeTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Complete task - TODO: Implement"})
}

func assignTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Assign task - TODO: Implement"})
}

func getFormSchema(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get form schema - TODO: Implement"})
}

func validateForm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Validate form - TODO: Implement"})
}

func listRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List rules - TODO: Implement"})
}

func createRule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create rule - TODO: Implement"})
}

func updateRule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update rule - TODO: Implement"})
}

func evaluateRule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Evaluate rule - TODO: Implement"})
}

func aiGenerateProcess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "AI generate process - TODO: Implement"})
}

func aiGenerateRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "AI generate rules - TODO: Implement"})
}

func aiOptimizeProcess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "AI optimize process - TODO: Implement"})
}

func getDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get dashboard - TODO: Implement"})
}

func getProcessAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get process analytics - TODO: Implement"})
}

func getInstanceAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get instance analytics - TODO: Implement"})
}

func listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users - TODO: Implement"})
}

func createUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create user - TODO: Implement"})
}

func updateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user - TODO: Implement"})
}

func deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete user - TODO: Implement"})
}

func updateUserRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user roles - TODO: Implement"})
}

func wsNotifications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket notifications - TODO: Implement"})
}

func wsTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket tasks - TODO: Implement"})
}
