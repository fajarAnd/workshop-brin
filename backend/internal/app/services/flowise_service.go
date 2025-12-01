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

type FlowiseService interface {
	SendMessageToWorkflow(ctx context.Context, userContext *models.UserContext, message string) error
	HandleWorkflowResponse(response *models.FlowiseResponse) error
	SetWhatsAppService(whatsappSvc WhatsAppService)
}

type flowiseService struct {
	httpClient  *http.Client
	baseURL     string
	flowID      string
	apiKey      string
	whatsappSvc WhatsAppService
}

type FlowiseConfig struct {
	BaseURL string
	FlowID  string
	APIKey  string
	Timeout time.Duration
}

func NewFlowiseService(config *FlowiseConfig) FlowiseService {
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &flowiseService{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		flowID:     config.FlowID,
		apiKey:     config.APIKey,
	}
}

func (s *flowiseService) SetWhatsAppService(whatsappSvc WhatsAppService) {
	s.whatsappSvc = whatsappSvc
}

func (s *flowiseService) SendMessageToWorkflow(ctx context.Context, userContext *models.UserContext, message string) error {
	log.Printf("[FlowiseService] Sending message to workflow for user %s: %s", userContext.Name, message)

	messageID := uuid.New().String()

	request := &models.FlowiseRequest{
		Question: message,
		OverrideConfig: &models.FlowiseOverrideConfig{
			SessionID: messageID,
			Vars: map[string]interface{}{
				"userContext": map[string]interface{}{
					"user_id": userContext.UserID,
					"name":    userContext.Name,
					"phone":   userContext.Phone,
					"email":   userContext.Email,
				},
				"timestamp": time.Now(),
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("[FlowiseService] Failed to marshal request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/prediction/%s", s.baseURL, s.flowID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[FlowiseService] Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[FlowiseService] Failed to send HTTP request: %v", err)
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[FlowiseService] Flowise API returned error status: %d", resp.StatusCode)
		return fmt.Errorf("Flowise API returned error status: %d", resp.StatusCode)
	}

	log.Printf("[FlowiseService] Message sent to Flowise workflow successfully (MessageID: %s)", messageID)
	return nil
}

func (s *flowiseService) HandleWorkflowResponse(response *models.FlowiseResponse) error {
	log.Printf("[FlowiseService] Handling workflow response for phone %s (MessageID: %s)", response.Phone, response.MessageID)

	if !response.Success {
		log.Printf("[FlowiseService] Flowise workflow returned error: %s", response.Error)
		return fmt.Errorf("Flowise workflow error: %s", response.Error)
	}

	if response.Text == "" {
		log.Printf("[FlowiseService] Empty response from Flowise workflow (MessageID: %s)", response.MessageID)
		return fmt.Errorf("empty response from Flowise workflow")
	}

	if s.whatsappSvc == nil {
		log.Printf("[FlowiseService] WhatsApp service not set, cannot send response")
		return fmt.Errorf("WhatsApp service not available")
	}

	ctx := context.Background()
	err := s.whatsappSvc.SendMessage(ctx, response.Phone, response.Text)
	if err != nil {
		log.Printf("[FlowiseService] Failed to send response to WhatsApp user %s: %v", response.Phone, err)
		return fmt.Errorf("failed to send response to WhatsApp: %w", err)
	}

	log.Printf("[FlowiseService] Response sent to WhatsApp user %s successfully (MessageID: %s)", response.Phone, response.MessageID)
	return nil
}
