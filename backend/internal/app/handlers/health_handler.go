package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler interface {
	HealthCheck(c *gin.Context)
	Status(c *gin.Context)
}

type healthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) HealthHandler {
	return &healthHandler{db: db}
}

func (h *healthHandler) HealthCheck(c *gin.Context) {
	// Simple health check
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "WhatsApp Bot is healthy",
		Data: models.HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Server:    "running",
			Database:  "not_checked",
		},
	})
}

func (h *healthHandler) Status(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Server:    "running",
		Database:  "unknown",
	}

	// Check database connectivity
	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			status.Status = "unhealthy"
			status.Database = "disconnected"
		} else {
			status.Database = "connected"
		}
	}

	// Set HTTP status based on overall health
	httpStatus := http.StatusOK
	if status.Status == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, models.APIResponse{
		Success: status.Status == "healthy",
		Message: "System status",
		Data:    status,
	})
}
