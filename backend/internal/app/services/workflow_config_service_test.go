package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/fajarAnd/workshop-brin/wa-service/testutils/mocks"
)

// TestWorkflowConfigService_GetActiveWorkflowType
// Summary: Test getting active workflow type from repository
// Purpose: Validate service correctly retrieves and logs workflow type from repository
func TestWorkflowConfigService_GetActiveWorkflowType(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockError      error
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Success - N8N workflow",
			mockResponse:   "n8n",
			mockError:      nil,
			expectedResult: "n8n",
			expectError:    false,
		},
		{
			name:           "Success - Flowise workflow",
			mockResponse:   "flowise",
			mockError:      nil,
			expectedResult: "flowise",
			expectError:    false,
		},
		{
			name:           "Repository error",
			mockResponse:   "",
			mockError:      fmt.Errorf("database connection failed"),
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository using mockery
			mockRepo := mocks.NewMockWorkflowConfigRepository(t)
			mockRepo.EXPECT().GetActiveWorkflowType(context.Background()).Return(tt.mockResponse, tt.mockError)

			// Create service with mock repository
			service := NewWorkflowConfigService(mockRepo)

			// Execute test
			ctx := context.Background()
			result, err := service.GetActiveWorkflowType(ctx)

			// Validate results
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tt.expectedResult != result {
				t.Errorf("Expected result %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

// TestWorkflowConfigService_GetActiveWorkflowType_Integration
// Summary: Test service behavior with realistic scenarios
// Purpose: Validate service handles various workflow configurations correctly
func TestWorkflowConfigService_GetActiveWorkflowType_Integration(t *testing.T) {
	tests := []struct {
		name           string
		workflowType   string
		expectedResult string
	}{
		{
			name:           "Default N8N configuration",
			workflowType:   "n8n",
			expectedResult: "n8n",
		},
		{
			name:           "Flowise configuration",
			workflowType:   "flowise",
			expectedResult: "flowise",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository with specific workflow type using mockery
			mockRepo := mocks.NewMockWorkflowConfigRepository(t)
			mockRepo.EXPECT().GetActiveWorkflowType(context.Background()).Return(tt.workflowType, nil)

			// Create service
			service := NewWorkflowConfigService(mockRepo)

			// Execute test
			ctx := context.Background()
			result, err := service.GetActiveWorkflowType(ctx)

			// Validate results
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expectedResult {
				t.Errorf("Expected workflow type %s, got %s", tt.expectedResult, result)
			}
		})
	}
}
