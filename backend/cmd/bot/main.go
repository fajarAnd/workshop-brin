package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fajarAnd/workshop-brin/wa-service/configs"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/handlers"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/repositories"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/server"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/services"
)

func main() {
	log.Println("Starting WhatsApp Bot...")

	// Load configuration
	config := configs.LoadConfig()
	log.Printf("Configuration loaded - Environment: %s, Server: %s",
		config.Server.Environment, config.GetServerAddress())

	// Connect to database
	db, err := configs.ConnectDatabase(&config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer configs.CloseDatabase(db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	workflowConfigRepo := repositories.NewWorkflowConfigRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	workflowConfigService := services.NewWorkflowConfigService(workflowConfigRepo)

	// Initialize N8N service
	n8nConfig := &services.N8NConfig{
		WorkflowURL: config.N8N.WebhookURL,
		APIKey:      config.N8N.APIKey,
		Timeout:     time.Duration(config.N8N.TimeoutSeconds) * time.Second,
	}
	n8nService := services.NewN8NService(n8nConfig)

	// Initialize Flowise service
	flowiseConfig := &services.FlowiseConfig{
		BaseURL: config.Flowise.BaseURL,
		FlowID:  config.Flowise.FlowID,
		APIKey:  config.Flowise.APIKey,
		Timeout: time.Duration(config.Flowise.TimeoutSeconds) * time.Second,
	}
	flowiseService := services.NewFlowiseService(flowiseConfig)

	// Initialize WhatsApp service
	whatsappService := services.NewWhatsAppService(userService, n8nService, flowiseService, workflowConfigService, db)

	// Set circular dependencies - workflow services need WhatsApp service for responses
	n8nService.SetWhatsAppService(whatsappService)
	flowiseService.SetWhatsAppService(whatsappService)

	// Initialize Signal service
	signalService := services.NewSignalService(userService, whatsappService)

	// Initialize handlers
	appHandlers := handlers.NewHandlers(db, userService, n8nService, whatsappService, signalService)

	// Start WhatsApp service
	ctx := context.Background()
	err = whatsappService.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start WhatsApp service: %v", err)
	}

	// Initialize and start HTTP server
	srv := server.NewServer(config, appHandlers)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("WhatsApp Bot started successfully")
	log.Println("Press Ctrl+C to gracefully shutdown...")

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down WhatsApp Bot...")

	// Graceful shutdown
	if err := srv.Stop(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	// Stop WhatsApp service
	if err := whatsappService.Stop(); err != nil {
		log.Printf("Error during WhatsApp service shutdown: %v", err)
	}

	log.Println("WhatsApp Bot stopped")
}
