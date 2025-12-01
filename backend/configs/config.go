package configs

import (
	"log"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	N8N      N8NConfig
	Flowise  FlowiseConfig
	WhatsApp WhatsAppConfig
}

type ServerConfig struct {
	Host            string
	Port            int
	Environment     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type N8NConfig struct {
	WebhookURL      string
	TimeoutSeconds  int
	RetryAttempts   int
	RetryDelay      time.Duration
	APIKey          string
}

type FlowiseConfig struct {
	BaseURL        string
	FlowID         string
	APIKey         string
	TimeoutSeconds int
}

type WhatsAppConfig struct {
	SessionTimeout time.Duration
	QRTimeout      time.Duration
}

// LoadConfig loads application configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading: %v", err)
	}

	config := &Config{
		Server: ServerConfig{
			Host:            getEnvString("SERVER_HOST", "0.0.0.0"),
			Port:            getEnvInt("SERVER_PORT", 8082),
			Environment:     getEnvString("ENVIRONMENT", "development"),
			ReadTimeout:     time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout:    time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
			ShutdownTimeout: time.Duration(getEnvInt("SERVER_SHUTDOWN_TIMEOUT", 30)) * time.Second,
		},
		Database: *LoadDatabaseConfig(),
		N8N: N8NConfig{
			WebhookURL:     getEnvString("N8N_WEBHOOK_URL", ""),
			TimeoutSeconds: getEnvInt("N8N_TIMEOUT_SECONDS", 30),
			RetryAttempts:  getEnvInt("N8N_RETRY_ATTEMPTS", 3),
			RetryDelay:     time.Duration(getEnvInt("N8N_RETRY_DELAY_SECONDS", 2)) * time.Second,
			APIKey:         getEnvString("N8N_API_KEY", ""),
		},
		Flowise: FlowiseConfig{
			BaseURL:        getEnvString("FLOWISE_BASE_URL", ""),
			FlowID:         getEnvString("FLOWISE_FLOW_ID", ""),
			APIKey:         getEnvString("FLOWISE_API_KEY", ""),
			TimeoutSeconds: getEnvInt("FLOWISE_TIMEOUT_SECONDS", 30),
		},
		WhatsApp: WhatsAppConfig{
			SessionTimeout: time.Duration(getEnvInt("WHATSAPP_SESSION_TIMEOUT", 3600)) * time.Second,
			QRTimeout:      time.Duration(getEnvInt("WHATSAPP_QR_TIMEOUT", 120)) * time.Second,
		},
	}

	return config
}

// GetServerAddress returns the server address string
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}
