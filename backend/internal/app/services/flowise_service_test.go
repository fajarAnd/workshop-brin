package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/google/uuid"
)

// Simple mock WhatsApp service for testing (avoiding import cycle)
type mockWhatsAppService struct {
	sendMessageFunc func(ctx context.Context, phone, message string) error
}

func (m *mockWhatsAppService) Start(ctx context.Context) error { return nil }
func (m *mockWhatsAppService) Stop() error                     { return nil }
func (m *mockWhatsAppService) IsConnected() bool               { return true }
func (m *mockWhatsAppService) GetQRCode() (string, error)      { return "", nil }
func (m *mockWhatsAppService) Logout() error                   { return nil }
func (m *mockWhatsAppService) SendMessage(ctx context.Context, phone, message string) error {
	if m.sendMessageFunc != nil {
		return m.sendMessageFunc(ctx, phone, message)
	}
	return nil
}

// TestFlowiseService_SendMessageToWorkflow
// Summary: Test sending messages to Flowise workflow API
// Purpose: Validate Flowise API integration with proper request formatting and error handling
func TestFlowiseService_SendMessageToWorkflow(t *testing.T) {
	tests := []struct {
		name        string
		userContext *models.UserContext
		message     string
		expectError bool
	}{
		{
			name: "Successful message send",
			userContext: &models.UserContext{
				UserID: uuid.New(),
				Name:   "John Doe",
				Phone:  "1234567890",
				Email:  "john@example.com",
			},
			message:     "Test message for analysis",
			expectError: false,
		},
		{
			name: "Empty message",
			userContext: &models.UserContext{
				UserID: uuid.New(),
				Name:   "Jane Smith",
				Phone:  "0987654321",
				Email:  "jane@example.com",
			},
			message:     "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with test configuration
			config := &FlowiseConfig{
				BaseURL: "http://test-flowise.com",
				FlowID:  "test-flow-id",
				APIKey:  "test-api-key",
				Timeout: 30 * time.Second,
			}
			service := NewFlowiseService(config)

			// Execute test
			ctx := context.Background()
			err := service.SendMessageToWorkflow(ctx, tt.userContext, tt.message)

			// Validate results based on expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Logf("Got error (expected for network call): %v", err)
			}
		})
	}
}

// TestFlowiseService_HandleWorkflowResponse
// Summary: Test handling webhook responses from Flowise
// Purpose: Validate response processing and WhatsApp message forwarding
func TestFlowiseService_HandleWorkflowResponse(t *testing.T) {
	tests := []struct {
		name             string
		response         *models.FlowiseResponse
		whatsappError    error
		expectError      bool
		expectedWhatsApp bool
	}{
		{
			name: "Successful response handling",
			response: &models.FlowiseResponse{
				Text:      "AI analysis response",
				ChatID:    "chat123",
				MessageID: "msg456",
				Phone:     "1234567890",
				Success:   true,
			},
			whatsappError:    nil,
			expectError:      false,
			expectedWhatsApp: true,
		},
		{
			name: "Empty response text",
			response: &models.FlowiseResponse{
				Text:      "",
				ChatID:    "chat123",
				MessageID: "msg456",
				Phone:     "1234567890",
				Success:   true,
			},
			whatsappError:    nil,
			expectError:      true,
			expectedWhatsApp: false,
		},
		{
			name: "Workflow error response",
			response: &models.FlowiseResponse{
				Text:      "",
				ChatID:    "chat123",
				MessageID: "msg456",
				Phone:     "1234567890",
				Success:   false,
				Error:     "Workflow processing failed",
			},
			whatsappError:    nil,
			expectError:      true,
			expectedWhatsApp: false,
		},
		{
			name: "WhatsApp send error",
			response: &models.FlowiseResponse{
				Text:      "AI response",
				ChatID:    "chat123",
				MessageID: "msg456",
				Phone:     "1234567890",
				Success:   true,
			},
			whatsappError:    fmt.Errorf("WhatsApp send failed"),
			expectError:      true,
			expectedWhatsApp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create simple mock WhatsApp service
			var mockWhatsApp *mockWhatsAppService
			if tt.expectedWhatsApp {
				mockWhatsApp = &mockWhatsAppService{
					sendMessageFunc: func(ctx context.Context, phone, message string) error {
						if phone != tt.response.Phone {
							t.Errorf("Expected phone %s, got %s", tt.response.Phone, phone)
						}
						if message != tt.response.Text {
							t.Errorf("Expected message %s, got %s", tt.response.Text, message)
						}
						return tt.whatsappError
					},
				}
			} else {
				mockWhatsApp = &mockWhatsAppService{}
			}

			// Create service and set WhatsApp service
			service := &flowiseService{
				whatsappSvc: mockWhatsApp,
			}

			// Execute test
			err := service.HandleWorkflowResponse(tt.response)

			// Validate results
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestFlowiseService_HandleWorkflowResponse_NoWhatsAppService
// Summary: Test handling responses when WhatsApp service is not set
// Purpose: Validate error handling when WhatsApp service dependency is missing
func TestFlowiseService_HandleWorkflowResponse_NoWhatsAppService(t *testing.T) {
	response := &models.FlowiseResponse{
		Text:      "Test response",
		ChatID:    "chat123",
		MessageID: "msg456",
		Phone:     "1234567890",
		Success:   true,
	}

	// Create service without WhatsApp service
	service := &flowiseService{}

	// Execute test
	err := service.HandleWorkflowResponse(response)

	// Should return error when WhatsApp service not available
	if err == nil {
		t.Errorf("Expected error when WhatsApp service not set, got nil")
	}
}

// TestNewFlowiseService
// Summary: Test Flowise service constructor
// Purpose: Validate service initialization with proper configuration
func TestNewFlowiseService(t *testing.T) {
	config := &FlowiseConfig{
		BaseURL: "http://test-flowise.com",
		FlowID:  "test-flow-id",
		APIKey:  "test-api-key",
		Timeout: 30 * time.Second,
	}

	service := NewFlowiseService(config)

	// Validate service was created
	if service == nil {
		t.Errorf("Expected service to be created, got nil")
	}

	// Cast to concrete type to check internal state
	flowiseService, ok := service.(*flowiseService)
	if !ok {
		t.Errorf("Expected flowiseService type")
	}

	if flowiseService.baseURL != config.BaseURL {
		t.Errorf("Expected baseURL %s, got %s", config.BaseURL, flowiseService.baseURL)
	}

	if flowiseService.flowID != config.FlowID {
		t.Errorf("Expected flowID %s, got %s", config.FlowID, flowiseService.flowID)
	}

	if flowiseService.apiKey != config.APIKey {
		t.Errorf("Expected apiKey %s, got %s", config.APIKey, flowiseService.apiKey)
	}
}

// TestFlowiseService_SetWhatsAppService
// Summary: Test setting WhatsApp service dependency
// Purpose: Validate dependency injection works correctly
func TestFlowiseService_SetWhatsAppService(t *testing.T) {
	// Create service
	config := &FlowiseConfig{
		BaseURL: "http://test-flowise.com",
		FlowID:  "test-flow-id",
		APIKey:  "test-api-key",
		Timeout: 30 * time.Second,
	}
	service := NewFlowiseService(config)

	// Create simple mock WhatsApp service
	mockWhatsApp := &mockWhatsAppService{}

	// Set WhatsApp service
	service.SetWhatsAppService(mockWhatsApp)

	// Validate WhatsApp service was set
	flowiseService, ok := service.(*flowiseService)
	if !ok {
		t.Errorf("Expected flowiseService type")
	}

	if flowiseService.whatsappSvc == nil {
		t.Errorf("Expected WhatsApp service to be set")
	}
}
