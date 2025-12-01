package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"
)

type SignalService interface {
	ProcessSignal(ctx context.Context, signal *models.Signal) (*models.SignalResponse, error)
	FormatSignalMessage(signal *models.Signal) string
}

type signalService struct {
	userService     UserService
	whatsappService WhatsAppService
}

func NewSignalService(userService UserService, whatsappService WhatsAppService) SignalService {
	return &signalService{
		userService:     userService,
		whatsappService: whatsappService,
	}
}

func (s *signalService) ProcessSignal(ctx context.Context, signal *models.Signal) (*models.SignalResponse, error) {
	log.Printf("[SignalService] Processing signal for ticker: %s", signal.Ticker)
	startTime := time.Now()

	users, err := s.userService.GetEligibleUsers(ctx)
	if err != nil {
		log.Printf("[SignalService] Failed to get eligible users: %v", err)
		return nil, fmt.Errorf("failed to get eligible users: %w", err)
	}

	if len(users) == 0 {
		log.Printf("[SignalService] No eligible users found")
		return &models.SignalResponse{
			Ticker:           signal.Ticker,
			UsersNotified:    0,
			Timestamp:        time.Now(),
			ProcessingTimeMs: time.Since(startTime).Milliseconds(),
		}, nil
	}

	message := s.FormatSignalMessage(signal)
	successCount := 0

	for _, user := range users {
		err := s.whatsappService.SendMessage(ctx, user.Phone, message)
		if err != nil {
			log.Printf("[SignalService] Failed to send signal to user %s (%s): %v", user.Name, user.Phone, err)
			continue
		}
		successCount++
		log.Printf("[SignalService] Signal sent successfully to user %s (%s)", user.Name, user.Phone)
	}

	processingTime := time.Since(startTime).Milliseconds()
	log.Printf("[SignalService] Signal processing completed for %s: %d/%d users notified in %dms",
		signal.Ticker, successCount, len(users), processingTime)

	return &models.SignalResponse{
		Ticker:           signal.Ticker,
		UsersNotified:    successCount,
		Timestamp:        time.Now(),
		ProcessingTimeMs: processingTime,
	}, nil
}

func (s *signalService) FormatSignalMessage(signal *models.Signal) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("🚀 *SIGNAL ALERT: %s*\n\n", signal.Ticker))

	builder.WriteString("📊 *TRADING DETAILS*\n")
	builder.WriteString(fmt.Sprintf("• Last Price: Rp %s\n", formatNumber(signal.LastClose)))
	builder.WriteString(fmt.Sprintf("• *Entry Price*: Rp %s\n", formatNumber(signal.EntryPrice)))
	builder.WriteString(fmt.Sprintf("• *Date*: %s\n", signal.LastDate))
	builder.WriteString(fmt.Sprintf("• Entry Gap: %+.1f%%\n\n", signal.EntryGapPercent))

	builder.WriteString("🎯 *TRADING LEVELS*\n")
	builder.WriteString(fmt.Sprintf("• *Stop Loss*: Rp %s\n", formatCurrency(signal.Stop)))
	builder.WriteString(fmt.Sprintf("• *Target*: Rp %s\n", formatCurrency(signal.Target)))
	builder.WriteString(fmt.Sprintf("• Risk/Reward: 1:%.2f\n\n", signal.RiskReward))

	builder.WriteString("📈 *ANALYSIS METRICS*\n")
	builder.WriteString(fmt.Sprintf("• *Confluence Score*: %.1f/10\n", signal.ConfluenceScore))
	if signal.ConfluenceHits != "" {
		builder.WriteString("• *Confluence Details*:\n")
		confluenceItems := strings.Split(signal.ConfluenceHits, "|")
		for _, item := range confluenceItems {
			item = strings.TrimSpace(item)
			if item != "" {
				builder.WriteString(fmt.Sprintf("  - ✓ %s\n", item))
			}
		}
	}
	builder.WriteString(fmt.Sprintf("• Backtest Win Rate: %.1f%%\n", signal.BacktestWinRate))
	builder.WriteString(fmt.Sprintf("• Total Trades: %d\n", signal.TotalTrades))
	builder.WriteString(fmt.Sprintf("• Confidence: %.1f%%\n\n", signal.ConfidenceScore))

	builder.WriteString("💭 *SENTIMENT ANALYSIS*\n")
	builder.WriteString(fmt.Sprintf("• Overall: %s\n", strings.Title(signal.OverallSentiment)))
	builder.WriteString(fmt.Sprintf("• Sentiment Score: %.1f%%\n\n", signal.SentimentScore))

	builder.WriteString("📋 *SUMMARY*\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", signal.AnalysisSummary))

	builder.WriteString("━━━━━━━━━━━━━━━━━━━━\n")
	builder.WriteString("❓ *Punya pertanyaan tentang signal ini?*\n")
	builder.WriteString("• Metode perhitungan\n")
	builder.WriteString("• Analisis teknikal\n")
	builder.WriteString("• Analisis sentimen\n")
	builder.WriteString("• Strategi trading\n\n")
	builder.WriteString("Silahkan tanyakan langsung! 💬")

	return builder.String()
}

func formatNumber(price int) string {
	priceStr := fmt.Sprintf("%d", price)
	result := ""

	for i, digit := range priceStr {
		if i > 0 && (len(priceStr)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	return result
}

func formatCurrency(price float64) string {
	priceInt := int(price)
	priceStr := fmt.Sprintf("%d", priceInt)
	result := ""

	for i, digit := range priceStr {
		if i > 0 && (len(priceStr)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	return result
}
