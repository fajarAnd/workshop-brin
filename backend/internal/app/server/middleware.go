package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[Server] %v | %3d | %13v | %15s | %-7s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	})
}

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("[Server] Panic recovered: %v", recovered)

		c.JSON(500, models.APIResponse{
			Success: false,
			Error:   "Internal server error",
		})
		c.Abort()
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestResponseLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		log.Printf("[Server] Request | %s %s | IP: %s | Headers: %v | Body: %s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP(),
			c.Request.Header, string(requestBody))

		c.Next()

		duration := time.Since(startTime)
		log.Printf("[Server] Response | %s %s | Status: %d | Duration: %v | Body: %s",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(),
			duration, w.body.String())
	}
}

var (
	rateLimiters = make(map[string]*rate.Limiter)
	mu           sync.RWMutex
)

func getRateLimiter(clientIP string) *rate.Limiter {
	mu.RLock()
	limiter, exists := rateLimiters[clientIP]
	mu.RUnlock()

	if !exists {
		mu.Lock()
		defer mu.Unlock()
		limiter = rate.NewLimiter(rate.Every(time.Minute), 100) // 100 requests per minute
		rateLimiters[clientIP] = limiter
	}

	return limiter
}

func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getRateLimiter(clientIP)

		if !limiter.Allow() {
			log.Printf("[Server] Rate limit exceeded for IP: %s", clientIP)
			c.JSON(http.StatusTooManyRequests, models.APIResponse{
				Success: false,
				Error:   "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "" && contentType != "application/json" {
				log.Printf("[Server] Security warning: Unusual content type from IP %s: %s", c.ClientIP(), contentType)
			}
		}

		c.Next()
	}
}
