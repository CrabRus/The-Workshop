package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret            string
	ExpirationHours   int
	RefreshExpiryDays int
}

type LoggerConfig struct {
	Level string
	JSON  bool
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "ecommerce_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			ExpirationHours:   getEnvInt("JWT_EXPIRY_HOURS", 24),
			RefreshExpiryDays: getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			JSON:  getEnvBool("LOG_JSON", false),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password +
		"@" + c.Database.Host + ":" + c.Database.Port +
		"/" + c.Database.DBName + "?sslmode=" + c.Database.SSLMode
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

func (c *Config) GetJWTExpiration() time.Duration {
	return time.Duration(c.JWT.ExpirationHours) * time.Hour
}
