package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/google/uuid"
)

type N8NService interface {
	SendMessageToWorkflow(ctx context.Context, userContext *models.UserContext, message string) error
	HandleWorkflowResponse(response *models.N8NResponse) error
	SetWhatsAppService(whatsappSvc WhatsAppService)
}

type n8nService struct {
	httpClient  *http.Client
	workflowURL string
	apiKey      string
	whatsappSvc WhatsAppService
}

type N8NConfig struct {
	WorkflowURL string
	APIKey      string
	Timeout     time.Duration
}

func NewN8NService(config *N8NConfig) N8NService {
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &n8nService{
		httpClient:  httpClient,
		workflowURL: config.WorkflowURL,
		apiKey:      config.APIKey,
	}
}

func (s *n8nService) SetWhatsAppService(whatsappSvc WhatsAppService) {
	s.whatsappSvc = whatsappSvc
}

func (s *n8nService) SendMessageToWorkflow(ctx context.Context, userContext *models.UserContext, message string) error {
	log.Printf("[N8NService] Sending message to workflow for user %s: %s", userContext.Name, message)

	// Generate message ID for correlation
	messageID := uuid.New().String()

	// Create N8N request payload
	request := &models.N8NRequest{
		UserContext: userContext,
		Message:     message,
		MessageID:   messageID,
		Timestamp:   time.Now(),
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("[N8NService] Failed to marshal request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.workflowURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[N8NService] Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[N8NService] Failed to send HTTP request: %v", err)
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("[N8NService] N8N workflow returned error status: %d", resp.StatusCode)
		return fmt.Errorf("N8N workflow returned error status: %d", resp.StatusCode)
	}

	log.Printf("[N8NService] Message sent to N8N workflow successfully (MessageID: %s)", messageID)
	return nil
}

func (s *n8nService) HandleWorkflowResponse(response *models.N8NResponse) error {
	log.Printf("[N8NService] Handling workflow response for phone %s (MessageID: %s)", response.Phone, response.MessageID)

	if !response.Success {
		log.Printf("[N8NService] N8N workflow returned error: %s", response.Error)
		return fmt.Errorf("N8N workflow error: %s", response.Error)
	}

	if response.Response == "" {
		log.Printf("[N8NService] Empty response from N8N workflow (MessageID: %s)", response.MessageID)
		return fmt.Errorf("empty response from N8N workflow")
	}

	// Send response back to WhatsApp user
	if s.whatsappSvc == nil {
		log.Printf("[N8NService] WhatsApp service not set, cannot send response")
		return fmt.Errorf("WhatsApp service not available")
	}

	ctx := context.Background()
	err := s.whatsappSvc.SendMessage(ctx, response.Phone, response.Response)
	if err != nil {
		log.Printf("[N8NService] Failed to send response to WhatsApp user %s: %v", response.Phone, err)
		return fmt.Errorf("failed to send response to WhatsApp: %w", err)
	}

	log.Printf("[N8NService] Response sent to WhatsApp user %s successfully (MessageID: %s)", response.Phone, response.MessageID)
	return nil
}
