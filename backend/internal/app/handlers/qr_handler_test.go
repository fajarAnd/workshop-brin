package handlers

import (
	"fmt"
	"github.com/fajarAnd/workshop-brin/wa-service/testutils/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGetQRCode
// Summary: Test QR code endpoint functionality
// Purpose: Validate QR code retrieval via HTTP endpoint
func TestGetQRCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		mockConnected   bool
		mockQRCode      string
		mockQRError     error
		expectedStatus  int
		expectedSuccess bool
	}{
		{
			name:            "QR code available",
			mockConnected:   false,
			mockQRCode:      "mock-qr-code-data",
			mockQRError:     nil,
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
		},
		{
			name:            "QR code not available",
			mockConnected:   false,
			mockQRCode:      "",
			mockQRError:     fmt.Errorf("QR code not available"),
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
		},
		{
			name:            "Already connected",
			mockConnected:   true,
			mockQRCode:      "",
			mockQRError:     fmt.Errorf("already connected"),
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWhatsAppService := mocks.NewMockWhatsAppService(t)
			mockWhatsAppService.On("GetQRCode").Return(tt.mockQRCode, tt.mockQRError)

			handler := NewQRHandler(mockWhatsAppService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/api/v1/qr/", nil)
			c.Request = req

			handler.GetQRCode(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestGetConnectionStatus
// Summary: Test connection status endpoint
// Purpose: Validate WhatsApp connection status reporting
func TestGetConnectionStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockConnected  bool
		expectedStatus int
	}{
		{
			name:           "Connected",
			mockConnected:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Disconnected",
			mockConnected:  false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWhatsAppService := mocks.NewMockWhatsAppService(t)
			mockWhatsAppService.On("IsConnected").Return(tt.mockConnected)

			handler := NewQRHandler(mockWhatsAppService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/api/v1/qr/status", nil)
			c.Request = req

			handler.GetConnectionStatus(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestShowQRPage
// Summary: Test QR page HTML rendering
// Purpose: Validate HTML page rendering for different connection states
func TestShowQRPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		mockConnected    bool
		mockQRCode       string
		mockQRError      error
		expectedStatus   int
		expectedContains string
	}{
		{
			name:             "Connected - shows success page",
			mockConnected:    true,
			expectedStatus:   http.StatusOK,
			expectedContains: "WhatsApp Bot Connected!",
		},
		{
			name:             "QR code available - shows QR page",
			mockConnected:    false,
			mockQRCode:       "mock-qr-code",
			expectedStatus:   http.StatusOK,
			expectedContains: "WhatsApp Bot Setup",
		},
		{
			name:             "QR code not available - shows error page",
			mockConnected:    false,
			mockQRError:      fmt.Errorf("QR not ready"),
			expectedStatus:   http.StatusOK,
			expectedContains: "QR Code Not Available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWhatsAppService := mocks.NewMockWhatsAppService(t)
			mockWhatsAppService.On("IsConnected").Return(tt.mockConnected)
			if !tt.mockConnected {
				mockWhatsAppService.On("GetQRCode").Return(tt.mockQRCode, tt.mockQRError)
			}

			handler := NewQRHandler(mockWhatsAppService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/api/v1/qr/page", nil)
			c.Request = req

			handler.ShowQRPage(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedContains != "" {
				body := w.Body.String()
				if !strings.Contains(body, tt.expectedContains) {
					t.Errorf("Expected response to contain '%s'", tt.expectedContains)
				}
			}
		})
	}
}
