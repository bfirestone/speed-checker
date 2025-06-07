# Speed Checker Makefile
.PHONY: help generate build clean test docker

# Default target
help:
	@echo "Speed Checker Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  generate     - Generate code from OpenAPI spec"
	@echo "  build        - Build the speed-checker binary"
	@echo "  clean        - Clean generated files and build artifacts"
	@echo "  test         - Run tests"
	@echo "  docker       - Build Docker image"
	@echo "  dev          - Run in development mode"
	@echo ""

# Generate code from OpenAPI spec
generate:
	@echo "🔨 Generating server code..."
	mkdir -p internal/api
	go generate -x ./internal/api
	
	@echo "🔨 Generating client SDK..."
	mkdir -p internal/client
	go generate -x ./internal/client
	
	@echo "🔨 Generating shared types..."
	mkdir -p internal/types
	go generate -x ./internal/types
	
	@echo "✅ Code generation complete!"

# Build the application
build: generate
	@echo "🏗️ Building speed-checker..."
	go build -o speed-checker .
	@echo "✅ Build complete!"

# Clean generated files and build artifacts
clean:
	@echo "🧹 Cleaning generated files..."
	rm -f internal/api/server.gen.go
	rm -f internal/client/client.gen.go
	rm -f internal/types/api.gen.go
	rm -f speed-checker
	@echo "✅ Clean complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	docker build -t speed-checker .

# Development mode
dev: generate
	@echo "🚀 Starting development server..."
	go run ./cmd/speed-checker all 