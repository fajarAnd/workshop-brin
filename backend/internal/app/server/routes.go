package server

import (
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handlers *handlers.Handlers) {
	// API version group
	api := r.Group("/api/v1")

	// Health endpoints
	health := api.Group("/health")
	{
		health.GET("/", handlers.Health.HealthCheck)
		health.GET("/status", handlers.Health.Status)
	}

	// Webhook endpoints
	webhook := api.Group("/webhook")
	{
		webhook.POST("/n8n/response", handlers.Webhook.HandleN8NResponse)
		webhook.POST("/n8n/signal", handlers.Webhook.HandleN8NSignal)
	}

	// QR Code endpoints for WhatsApp bot setup
	qr := api.Group("/qr")
	{
		qr.GET("/", handlers.QR.GetQRCode)                 // JSON response
		qr.GET("/status", handlers.QR.GetConnectionStatus) // JSON response
		qr.GET("/page", handlers.QR.ShowQRPage)            // HTML page
		qr.GET("/image", handlers.QR.GetQRImage)           // PNG image
	}

	// WhatsApp management endpoints
	whatsapp := api.Group("/whatsapp")
	{
		whatsapp.POST("/logout", handlers.WhatsApp.Logout)
		whatsapp.GET("/status", handlers.WhatsApp.GetConnectionStatus)
	}

	// Root health check (for load balancers)
	r.GET("/health", handlers.Health.HealthCheck)
}
