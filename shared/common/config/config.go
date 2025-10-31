package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	NATS     NATSConfig     `mapstructure:"nats"`
	Redis    RedisConfig    `mapstructure:"redis"`
	AI       AIConfig       `mapstructure:"ai"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Security SecurityConfig `mapstructure:"security"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	TLS          TLSConfig     `mapstructure:"tls"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"dbname"`
	SSLMode      string        `mapstructure:"sslmode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Provider string             `mapstructure:"provider"` // jwt, keycloak, auth0
	JWT      JWTConfig          `mapstructure:"jwt"`
	Keycloak KeycloakConfig     `mapstructure:"keycloak"`
	OIDC     OIDCConfig         `mapstructure:"oidc"`
	BuiltIn  BuiltInConfig      `mapstructure:"built_in"`
	Security AuthSecurityConfig `mapstructure:"security"`
}

// JWTConfig contains JWT configuration
type JWTConfig struct {
	Secret               string        `mapstructure:"secret"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	Issuer               string        `mapstructure:"issuer"`
	Audience             string        `mapstructure:"audience"`
}

// KeycloakConfig contains Keycloak configuration
type KeycloakConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// OIDCConfig contains OIDC provider configuration
type OIDCConfig struct {
	Provider     string `mapstructure:"provider"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Domain       string `mapstructure:"domain"`
	Audience     string `mapstructure:"audience"`
}

// BuiltInConfig contains built-in auth configuration
type BuiltInConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	RequireOTP       bool `mapstructure:"require_otp"`
	PasswordlessOnly bool `mapstructure:"passwordless_only"`
}

// AuthSecurityConfig contains auth security settings
type AuthSecurityConfig struct {
	RequireMFA            bool          `mapstructure:"require_mfa"`
	SessionTimeout        time.Duration `mapstructure:"session_timeout"`
	MaxConcurrentSessions int           `mapstructure:"max_concurrent_sessions"`
	IPAllowlist           []string      `mapstructure:"ip_allowlist"`
}

// NATSConfig contains NATS configuration
type NATSConfig struct {
	URL       string `mapstructure:"url"`
	ClusterID string `mapstructure:"cluster_id"`
	ClientID  string `mapstructure:"client_id"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// AIConfig contains AI integration configuration
type AIConfig struct {
	OpenAI OpenAIConfig   `mapstructure:"openai"`
	Custom CustomAIConfig `mapstructure:"custom"`
}

// OpenAIConfig contains OpenAI configuration
type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

// CustomAIConfig contains custom AI service configuration
type CustomAIConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	APIKey   string `mapstructure:"api_key"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json, text
	Output string `mapstructure:"output"` // stdout, file
	File   string `mapstructure:"file"`
}

// CacheConfig contains caching configuration
type CacheConfig struct {
	TTL             time.Duration `mapstructure:"ttl"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	MaxSize         int           `mapstructure:"max_size"`
	Strategy        string        `mapstructure:"strategy"` // lru, lfu
}

// SecurityConfig contains security configuration
type SecurityConfig struct {
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// EncryptionConfig contains encryption configuration
type EncryptionConfig struct {
	Key       string `mapstructure:"key"`
	Algorithm string `mapstructure:"algorithm"`
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

// Load loads configuration from files and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Environment variable configuration
	viper.SetEnvPrefix("BPMS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is okay, we'll use defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8081)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("server.tls.enabled", false)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "bpms_user")
	viper.SetDefault("database.password", "bpms_pass")
	viper.SetDefault("database.dbname", "ai_bpms")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_lifetime", "5m")

	// Auth defaults
	viper.SetDefault("auth.provider", "jwt")
	viper.SetDefault("auth.jwt.secret", "change-me-in-production")
	viper.SetDefault("auth.jwt.access_token_duration", "15m")
	viper.SetDefault("auth.jwt.refresh_token_duration", "168h") // 7 days
	viper.SetDefault("auth.jwt.issuer", "ai-bpms")
	viper.SetDefault("auth.jwt.audience", "ai-bpms-users")
	viper.SetDefault("auth.built_in.enabled", true)
	viper.SetDefault("auth.security.session_timeout", "4h")
	viper.SetDefault("auth.security.max_concurrent_sessions", 3)

	// NATS defaults
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("nats.cluster_id", "bpms-cluster")
	viper.SetDefault("nats.client_id", "bpms-client")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// AI defaults
	viper.SetDefault("ai.openai.base_url", "https://api.openai.com/v1")
	viper.SetDefault("ai.openai.model", "gpt-4")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// Cache defaults
	viper.SetDefault("cache.ttl", "1h")
	viper.SetDefault("cache.cleanup_interval", "10m")
	viper.SetDefault("cache.max_size", 1000)
	viper.SetDefault("cache.strategy", "lru")

	// Security defaults
	viper.SetDefault("security.rate_limit.enabled", true)
	viper.SetDefault("security.rate_limit.rps", 100)
	viper.SetDefault("security.rate_limit.burst", 200)
	viper.SetDefault("security.cors.allowed_origins", []string{"*"})
	viper.SetDefault("security.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("security.cors.allowed_headers", []string{"*"})
	viper.SetDefault("security.encryption.algorithm", "AES-256-GCM")

	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.port", 9090)
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr returns Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
