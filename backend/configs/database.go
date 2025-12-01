package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	Database          string
	SSLMode           string
	MaxConnections    int
	MinConnections    int
	MaxLifetime       time.Duration
	MaxIdleTime       time.Duration
	HealthCheckPeriod time.Duration
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	config := &DatabaseConfig{
		Host:              getEnvString("DB_HOST", "localhost"),
		Port:              getEnvInt("DB_PORT", 5432),
		Username:          getEnvString("DB_USER", "postgres"),
		Password:          getEnvString("DB_PASSWORD", ""),
		Database:          getEnvString("DB_NAME", "wa-service"),
		SSLMode:           getEnvString("DB_SSL_MODE", "disable"),
		MaxConnections:    getEnvInt("DB_MAX_CONNECTIONS", 25),
		MinConnections:    getEnvInt("DB_MIN_CONNECTIONS", 5),
		MaxLifetime:       time.Duration(getEnvInt("DB_MAX_LIFETIME_MINUTES", 60)) * time.Minute,
		MaxIdleTime:       time.Duration(getEnvInt("DB_MAX_IDLE_MINUTES", 30)) * time.Minute,
		HealthCheckPeriod: time.Duration(getEnvInt("DB_HEALTH_CHECK_MINUTES", 1)) * time.Minute,
	}

	return config
}

// GetConnectionString returns PostgreSQL connection string
func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

// ConnectDatabase creates a new database connection pool
func ConnectDatabase(config *DatabaseConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Build connection config
	pgxConfig, err := pgxpool.ParseConfig(config.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	pgxConfig.MaxConns = int32(config.MaxConnections)
	pgxConfig.MinConns = int32(config.MinConnections)
	pgxConfig.MaxConnLifetime = config.MaxLifetime
	pgxConfig.MaxConnIdleTime = config.MaxIdleTime
	pgxConfig.HealthCheckPeriod = config.HealthCheckPeriod

	// Disable prepared statement cache to avoid conflicts with other drivers
	// This prevents "prepared statement already exists" errors when using multiple database drivers
	// Use simple protocol to avoid prepared statement caching
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully: %s:%d/%s", config.Host, config.Port, config.Database)
	return pool, nil
}

// CloseDatabase gracefully closes the database connection pool
func CloseDatabase(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("Database connection pool closed")
	}
}

// Helper functions to get environment variables with defaults
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
