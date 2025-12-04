# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Set GOTOOLCHAIN to auto to allow downloading required Go version
ENV GOTOOLCHAIN=auto

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, GOOS=linux for linux target
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags '-extldflags "-static" -s -w' \
    -o bot ./cmd/bot

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and tools for entrypoint
RUN apk --no-cache add ca-certificates tzdata wget postgresql-client

# Install golang-migrate
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin migrate

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bot ./api

# Copy migrations directory
COPY --from=builder /build/migrations ./migrations

# Copy entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    chown appuser:appuser /usr/local/bin/docker-entrypoint.sh

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set entrypoint to run migrations before starting app
ENTRYPOINT ["docker-entrypoint.sh"]

# Run the application
CMD ["./api"]