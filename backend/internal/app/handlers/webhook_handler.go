package handlers

import (
	"log"
	"net/http"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/services"

	"github.com/gin-gonic/gin"
)

type WebhookHandler interface {
	HandleN8NResponse(c *gin.Context)
	HandleN8NSignal(c *gin.Context)
}

type webhookHandler struct {
	n8nService    services.N8NService
	signalService services.SignalService
}

func NewWebhookHandler(n8nService services.N8NService, signalService services.SignalService) WebhookHandler {
	return &webhookHandler{
		n8nService:    n8nService,
		signalService: signalService,
	}
}

func (h *webhookHandler) HandleN8NResponse(c *gin.Context) {
	log.Printf("[WebhookHandler] Received N8N response from %s", c.ClientIP())

	var response models.N8NResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		log.Printf("[WebhookHandler] Invalid JSON payload: %v", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid JSON payload",
		})
		return
	}

	// Log the received response
	log.Printf("[WebhookHandler] N8N Response - MessageID: %s, Phone: %s, Success: %t",
		response.MessageID, response.Phone, response.Success)

	// Handle the response using N8N service
	err := h.n8nService.HandleWorkflowResponse(&response)
	if err != nil {
		log.Printf("[WebhookHandler] Failed to handle workflow response: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to process response",
		})
		return
	}

	// Send acknowledgment
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Response processed successfully",
		Data: gin.H{
			"message_id": response.MessageID,
			"processed":  true,
		},
	})

	log.Printf("[WebhookHandler] N8N response processed successfully")
}

func (h *webhookHandler) HandleN8NSignal(c *gin.Context) {
	log.Printf("[WebhookHandler] Received N8N signal from %s", c.ClientIP())

	var signal models.Signal
	if err := c.ShouldBindJSON(&signal); err != nil {
		log.Printf("[WebhookHandler] Invalid signal JSON payload: %v", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid JSON payload",
		})
		return
	}

	log.Printf("[WebhookHandler] Processing signal for ticker: %s, entry price: %d", signal.Ticker, signal.EntryPrice)

	response, err := h.signalService.ProcessSignal(c.Request.Context(), &signal)
	if err != nil {
		log.Printf("[WebhookHandler] Failed to process signal: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to process signal",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Signal processed successfully",
		Data:    response,
	})

	log.Printf("[WebhookHandler] Signal processed successfully for %s, notified %d users",
		signal.Ticker, response.UsersNotified)
}
