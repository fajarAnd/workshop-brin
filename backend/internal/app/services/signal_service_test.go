package services

import (
	"testing"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/stretchr/testify/assert"
)

// TestFormatNumber
// Summary: Test number formatting for Indonesian currency
// Purpose: Validate that numbers are properly formatted with comma separators
func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "Single digit",
			input:    5,
			expected: "5",
		},
		{
			name:     "Three digits",
			input:    123,
			expected: "123",
		},
		{
			name:     "Four digits",
			input:    1234,
			expected: "1,234",
		},
		{
			name:     "Five digits",
			input:    12345,
			expected: "12,345",
		},
		{
			name:     "Seven digits",
			input:    1234567,
			expected: "1,234,567",
		},
		{
			name:     "Large number (stock price)",
			input:    9420,
			expected: "9,420",
		},
		{
			name:     "Zero value",
			input:    0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatCurrency
// Summary: Test currency formatting for Indonesian Rupiah
// Purpose: Validate that float prices are properly formatted with comma separators
func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "Small price",
			input:    123.0,
			expected: "123",
		},
		{
			name:     "Four digit price",
			input:    1234.5,
			expected: "1,234",
		},
		{
			name:     "Stock price with decimals",
			input:    9180.75,
			expected: "9,180",
		},
		{
			name:     "Large price",
			input:    123456.99,
			expected: "123,456",
		},
		{
			name:     "Zero value",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "Fractional value rounds down",
			input:    999.99,
			expected: "999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCurrency(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatSignalMessage
// Summary: Test signal message formatting includes entry price
// Purpose: Validate that formatted signal messages contain entry price information
func TestFormatSignalMessage(t *testing.T) {
	tests := []struct {
		name             string
		signal           *models.Signal
		expectedContains []string
	}{
		{
			name: "Bullish signal with entry price",
			signal: &models.Signal{
				Ticker:           "BBCA",
				LastDate:         "2025-08-10",
				LastClose:        9420,
				EntryPrice:       9655,
				EntryGapPercent:  2.5,
				Stop:             9180.0,
				Target:           9720.0,
				RiskReward:       1.85,
				BacktestWinRate:  68.5,
				TotalTrades:      147,
				ConfluenceScore:  8.2,
				ConfluenceHits:   "Strong bullish signals",
				OverallSentiment: "bullish",
				ConfidenceScore:  85.3,
				SentimentScore:   72.8,
				AnalysisSummary:  "Strong technical confluence with bullish sentiment",
			},
			expectedContains: []string{
				"🚀 *SIGNAL ALERT: BBCA*",
				"• Last Price: Rp 9,420",
				"• *Entry Price*: Rp 9,655",
				"• Entry Gap: +2.5%",
				"• *Stop Loss*: Rp 9,180",
				"• *Target*: Rp 9,720",
				"• Risk/Reward: 1:1.85",
			},
		},
		{
			name: "Bearish signal with entry price",
			signal: &models.Signal{
				Ticker:           "TLKM",
				LastDate:         "2025-08-10",
				LastClose:        3500,
				EntryPrice:       3437,
				EntryGapPercent:  -1.8,
				Stop:             3600.0,
				Target:           3300.0,
				RiskReward:       1.67,
				BacktestWinRate:  72.1,
				TotalTrades:      89,
				ConfluenceScore:  7.8,
				ConfluenceHits:   "Bearish divergence signals",
				OverallSentiment: "bearish",
				ConfidenceScore:  78.9,
				SentimentScore:   45.2,
				AnalysisSummary:  "Bearish technical indicators with high volume",
			},
			expectedContains: []string{
				"🚀 *SIGNAL ALERT: TLKM*",
				"• Last Price: Rp 3,500",
				"• *Entry Price*: Rp 3,437",
				"• Entry Gap: -1.8%",
				"• *Stop Loss*: Rp 3,600",
				"• *Target*: Rp 3,300",
			},
		},
	}

	service := &signalService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FormatSignalMessage(tt.signal)

			for _, expected := range tt.expectedContains {
				assert.Contains(t, result, expected, "Message should contain: %s", expected)
			}

			assert.Contains(t, result, "📊 *TRADING DETAILS*", "Should contain trading details section")
			assert.Contains(t, result, "🎯 *TRADING LEVELS*", "Should contain trading levels section")
			assert.Contains(t, result, "📈 *ANALYSIS METRICS*", "Should contain analysis metrics section")
		})
	}
}
