# API-only Dockerfile for containerized deployments
# Build: docker build -t speed-checker .
# Run API: docker run speed-checker api
# Run Daemon: docker run speed-checker daemon --api-endpoint http://api:8080

FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o speed-checker .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    speedtest-cli \
    iperf3 \
    tzdata \
    wget

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary
COPY --from=builder /app/speed-checker .

# Set ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Expose port (only needed for API server)
EXPOSE 8080

# Set environment variables for PostgreSQL
ENV POSTGRES_HOST="postgres"
ENV POSTGRES_PORT="5432"
ENV POSTGRES_DB="speedchecker"
ENV POSTGRES_USER="speedchecker"
ENV POSTGRES_PASSWORD=""
ENV POSTGRES_SSLMODE="disable"
ENV SPEED_CHECKER_SERVER_HOST="0.0.0.0"
ENV SPEED_CHECKER_SERVER_PORT="8080"

# Health check (works for both API and daemon)
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD pgrep -f speed-checker || exit 1

# Default to API server (can be overridden)
CMD ["./speed-checker", "api"] 