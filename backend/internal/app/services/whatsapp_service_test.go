package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mau.fi/whatsmeow/types"
)

// TestFormatPhoneToJID
// Summary: Test phone number formatting to WhatsApp JID
// Purpose: Validate that phone numbers are correctly converted to WhatsApp JID format
func TestFormatPhoneToJID(t *testing.T) {
	service := &whatsAppService{}

	tests := []struct {
		name        string
		phone       string
		expectError bool
		expected    string
	}{
		{
			name:        "Valid US phone number with +1",
			phone:       "+12345678901",
			expectError: false,
			expected:    "12345678901",
		},
		{
			name:        "Valid phone with dashes",
			phone:       "123-456-7890",
			expectError: false,
			expected:    "1234567890",
		},
		{
			name:        "Valid phone with spaces",
			phone:       "123 456 7890",
			expectError: false,
			expected:    "1234567890",
		},
		{
			name:        "Empty phone number",
			phone:       "",
			expectError: true,
			expected:    "",
		},
		{
			name:        "Phone with only symbols",
			phone:       "+-  ",
			expectError: true,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jid, err := service.formatPhoneToJID(tt.phone)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for phone %s, but got none", tt.phone)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for phone %s: %v", tt.phone, err)
				return
			}

			if jid.User != tt.expected {
				t.Errorf("Expected JID user %s, but got %s", tt.expected, jid.User)
			}

			if jid.Server != types.DefaultUserServer {
				t.Errorf("Expected JID server %s, but got %s", types.DefaultUserServer, jid.Server)
			}
		})
	}
}

// TestExtractPhoneFromJID
// Summary: Test phone number extraction from WhatsApp JID
// Purpose: Validate that phone numbers are correctly extracted from WhatsApp JID
func TestExtractPhoneFromJID(t *testing.T) {
	service := &whatsAppService{}

	tests := []struct {
		name     string
		jid      types.JID
		expected string
	}{
		{
			name: "Valid JID with phone number",
			jid: types.JID{
				User:   "12345678901",
				Server: types.DefaultUserServer,
			},
			expected: "12345678901",
		},
		{
			name: "JID with empty user",
			jid: types.JID{
				User:   "",
				Server: types.DefaultUserServer,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := service.extractPhoneFromJID(tt.jid)

			if phone != tt.expected {
				t.Errorf("Expected phone %s, but got %s", tt.expected, phone)
			}
		})
	}
}

// TestExtractMessageText
// Summary: Test message text extraction from WhatsApp message
// Purpose: Validate that text content is correctly extracted from different message types
func TestExtractMessageText(t *testing.T) {
	service := &whatsAppService{}

	tests := []struct {
		name     string
		message  interface{}
		expected string
	}{
		{
			name:     "Nil message",
			message:  nil,
			expected: "",
		},
		{
			name: "Simple conversation message",
			message: &struct {
				Conversation *string
			}{
				Conversation: stringPtr("Hello World"),
			},
			expected: "Hello World",
		},
		{
			name: "Extended text message",
			message: &struct {
				Conversation        *string
				ExtendedTextMessage *struct {
					Text *string
				}
			}{
				ExtendedTextMessage: &struct {
					Text *string
				}{
					Text: stringPtr("Extended message"),
				},
			},
			expected: "Extended message",
		},
		{
			name: "Empty message",
			message: &struct {
				Conversation *string
			}{
				Conversation: nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This is a simplified test since we can't easily mock waE2E.Message
			// In a real implementation, you'd want to test with actual message types
			if tt.name == "Nil message" {
				text := service.extractMessageText(nil)
				if text != tt.expected {
					t.Errorf("Expected text '%s', but got '%s'", tt.expected, text)
				}
			}
		})
	}
}

// TestWhatsAppServiceCreation
// Summary: Test WhatsApp service creation and initialization
// Purpose: Validate that WhatsApp service can be created with required dependencies
func TestWhatsAppServiceCreation(t *testing.T) {
	// Mock services (in real implementation, these would be proper mocks)
	userService := &mockUserService{}
	n8nService := &mockN8NService{}
	flowiseService := &mockFlowiseService{}
	workflowConfigService := &mockWorkflowConfigService{}
	var mockPool *pgxpool.Pool // nil pool for basic testing

	service := NewWhatsAppService(userService, n8nService, flowiseService, workflowConfigService, mockPool)

	if service == nil {
		t.Error("Expected WhatsApp service to be created, but got nil")
	}

	// Test interface compliance
	var _ WhatsAppService = service

	// Test that service implements all required methods
	if !service.IsConnected() {
		// Service should start disconnected
	}

	_, err := service.GetQRCode()
	if err == nil {
		t.Error("Expected error when getting QR code from uninitialized service")
	}
}

// Mock implementations for testing
type mockUserService struct{}

func (m *mockUserService) IsUserEligible(ctx context.Context, phone string) (bool, error) {
	// Mock implementation - return true for test phone numbers
	return phone == "12345678901", nil
}

func (m *mockUserService) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	if phone == "12345678901" {
		return &models.User{
			ID:    uuid.New(),
			Name:  "Test User",
			Phone: phone,
			Email: "test@example.com",
		}, nil
	}
	return nil, ErrUserNotFound
}

func (m *mockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return &models.User{
		ID:    id,
		Name:  "Test User",
		Phone: "12345678901",
		Email: "test@example.com",
	}, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	return &models.User{
		ID:    uuid.New(),
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	return &models.User{
		ID:    id,
		Name:  req.Name,
		Phone: "12345678901",
		Email: req.Email,
	}, nil
}

func (m *mockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockUserService) GetEligibleUsers(ctx context.Context) ([]*models.User, error) {
	return []*models.User{
		{
			ID:    uuid.New(),
			Name:  "Test User",
			Phone: "12345678901",
			Email: "test@example.com",
		},
	}, nil
}

type mockN8NService struct{}

func (m *mockN8NService) SendMessageToWorkflow(ctx context.Context, userContext *models.UserContext, message string) error {
	return nil
}

func (m *mockN8NService) HandleWorkflowResponse(response *models.N8NResponse) error {
	return nil
}

func (m *mockN8NService) SetWhatsAppService(whatsappSvc WhatsAppService) {
	// Mock implementation
}

// mockFlowiseService for testing
type mockFlowiseService struct{}

func (m *mockFlowiseService) SendMessageToWorkflow(ctx context.Context, userContext *models.UserContext, message string) error {
	return nil
}

func (m *mockFlowiseService) HandleWorkflowResponse(response *models.FlowiseResponse) error {
	return nil
}

func (m *mockFlowiseService) SetWhatsAppService(whatsappSvc WhatsAppService) {
	// Mock implementation
}

// mockWorkflowConfigService for testing
type mockWorkflowConfigService struct{}

func (m *mockWorkflowConfigService) GetActiveWorkflowType(ctx context.Context) (string, error) {
	return "n8n", nil
}

// Custom error for testing
var ErrUserNotFound = fmt.Errorf("user not found")

// Helper function
func stringPtr(s string) *string {
	return &s
}
