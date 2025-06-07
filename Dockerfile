# Build stage for frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production

COPY frontend/ ./
RUN npm run build

# Build stage for Go backend
FROM golang:1.21-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Copy the built frontend from the previous stage
COPY --from=frontend-builder /app/frontend/build ./frontend/build

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o speed-checker .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    speedtest-cli \
    iperf3 \
    tzdata

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary and frontend build
COPY --from=backend-builder /app/speed-checker .
COPY --from=backend-builder /app/frontend/build ./frontend/build

# Create data directory for SQLite database
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV SPEED_CHECKER_DATABASE_DSN="/app/data/speedtest_results.db?_fk=1"
ENV SPEED_CHECKER_SERVER_HOST="0.0.0.0"
ENV SPEED_CHECKER_SERVER_PORT="8080"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/dashboard || exit 1

CMD ["./speed-checker"] 