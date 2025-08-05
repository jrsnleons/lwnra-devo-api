FROM golang:1.22-alpine AS builder

# Install build dependencies for CGO (needed for SQLite)
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Remove any test files to avoid conflicts
RUN rm -f test_parser.go test_parser.go.bak

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o main ./cmd/server

# Final stage - use minimal image for production
FROM alpine:latest

# Install runtime dependencies including wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget && \
    adduser -D -s /bin/sh appuser

# Set timezone to Philippine time for scheduler
ENV TZ=Asia/Manila

# Create app directory
WORKDIR /app

# Create data directory for database with proper permissions
RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user for security
RUN chown appuser:appuser /app/main && \
    chmod +x /app/main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8082

# Health check to ensure container is running properly
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application
CMD ["./main"]
