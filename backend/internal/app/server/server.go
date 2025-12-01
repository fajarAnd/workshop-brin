package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fajarAnd/workshop-brin/wa-service/configs"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	gin      *gin.Engine
	config   *configs.Config
	server   *http.Server
	handlers *handlers.Handlers
}

func NewServer(config *configs.Config, handlers *handlers.Handlers) *Server {
	// Set gin mode based on environment
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return &Server{
		config:   config,
		handlers: handlers,
	}
}

func (s *Server) setupRoutes() {
	s.gin = gin.New()

	// Add middleware
	s.gin.Use(LoggingMiddleware())
	s.gin.Use(ErrorHandlingMiddleware())
	s.gin.Use(SecurityMiddleware())
	s.gin.Use(RateLimitingMiddleware())
	s.gin.Use(CORSMiddleware())
	s.gin.Use(RequestResponseLoggingMiddleware())

	// Setup routes
	SetupRoutes(s.gin, s.handlers)
}

func (s *Server) Start() error {
	s.setupRoutes()

	s.server = &http.Server{
		Addr:         s.config.GetServerAddress(),
		Handler:      s.gin,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	log.Printf("[Server] Starting HTTP server on %s", s.config.GetServerAddress())
	log.Printf("[Server] Environment: %s", s.config.Server.Environment)

	// Start server in goroutine so it doesn't block
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] Failed to start server: %v", err)
		}
	}()

	log.Printf("[Server] HTTP server started successfully")
	return nil
}

func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	log.Printf("[Server] Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("[Server] Error during server shutdown: %v", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Printf("[Server] HTTP server stopped gracefully")
	return nil
}

func (s *Server) GetAddress() string {
	return s.config.GetServerAddress()
}
