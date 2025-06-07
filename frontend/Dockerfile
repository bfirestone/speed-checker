# Frontend Dockerfile for SvelteKit application
# Build: docker build -t speed-checker-frontend ./frontend
# Run: docker run -p 3000:3000 speed-checker-frontend

FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:22-alpine

WORKDIR /app

# Copy built application from Node.js adapter output
COPY --from=builder /app/build ./
COPY --from=builder /app/package*.json ./

# Install only production dependencies
RUN npm ci --omit=dev --legacy-peer-deps

# Create non-root user
RUN addgroup -g 1001 -S sveltegroup && \
    adduser -u 1001 -S svelteuser -G sveltegroup

# Change ownership
RUN chown -R svelteuser:sveltegroup /app
USER svelteuser

# Expose port
EXPOSE 3000

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

# Start the application (SvelteKit server)
CMD ["node", "index.js"] 