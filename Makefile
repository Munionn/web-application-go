# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build  2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build ; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Reset DB container (removes volume and recreates with fresh data)
docker-reset:
	@echo "Stopping and removing database container and volume..."
	@if docker compose down -v 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down -v; \
	fi
	@echo "Starting fresh database container..."
	@if docker compose up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi
	@echo "Database reset complete! Waiting 5 seconds for PostgreSQL to initialize..."
	@sleep 5

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload (like nodemon for Node.js - rebuilds on file changes)
# Uses Air: https://github.com/air-verse/air
watch:
	@echo "Starting live reload (Air)..."
	@echo "Watching for changes. Edit any .go file to trigger rebuild."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	elif go run github.com/air-verse/air@latest 2>/dev/null; then \
		: ; \
	else \
		echo "Installing Air (go install github.com/air-verse/air@latest)..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

.PHONY: all build run test clean watch docker-run docker-down docker-reset itest
