package config

import "os"

// Config holds application configuration
type Config struct {
	Port            string
	DatabasePath    string
	FacebookToken   string
	Environment     string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabasePath:  getEnv("DB_PATH", "devotionals.db"),
		FacebookToken: getEnv("FB_ACCESS_TOKEN", ""),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
