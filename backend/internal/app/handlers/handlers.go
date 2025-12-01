package handlers

import (
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Health   HealthHandler
	Webhook  WebhookHandler
	QR       QRHandler
	WhatsApp WhatsAppHandler
}

func NewHandlers(db *pgxpool.Pool, userService services.UserService, n8nService services.N8NService, whatsappService services.WhatsAppService, signalService services.SignalService) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(db),
		Webhook:  NewWebhookHandler(n8nService, signalService),
		QR:       NewQRHandler(whatsappService),
		WhatsApp: NewWhatsAppHandler(whatsappService),
	}
}
