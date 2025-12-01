package handlers

import (
	"net/http"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/services"

	"github.com/gin-gonic/gin"
)

type WhatsAppHandler interface {
	Logout(c *gin.Context)
	GetConnectionStatus(c *gin.Context)
}

type whatsappHandler struct {
	whatsappService services.WhatsAppService
}

func NewWhatsAppHandler(whatsappService services.WhatsAppService) WhatsAppHandler {
	return &whatsappHandler{
		whatsappService: whatsappService,
	}
}

func (h *whatsappHandler) Logout(c *gin.Context) {
	err := h.whatsappService.Logout()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to logout from WhatsApp",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully logged out from WhatsApp. Device has been removed.",
	})
}

func (h *whatsappHandler) GetConnectionStatus(c *gin.Context) {
	isConnected := h.whatsappService.IsConnected()

	status := "disconnected"
	message := "WhatsApp bot is not connected. Please scan QR code to connect."

	if isConnected {
		status = "connected"
		message = "WhatsApp bot is connected and ready to receive messages."
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"connected": isConnected,
		"status":    status,
		"message":   message,
	})
}
