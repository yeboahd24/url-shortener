# URL Shortener Makefile

.PHONY: help build run test clean docker-build docker-run docker-stop swagger deps

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the Go application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose (development)"
	@echo "  docker-prod   - Run with Docker Compose (production)"
	@echo "  docker-stop   - Stop Docker Compose services"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  logs          - View application logs"
	@echo "  db-migrate    - Run database migrations"

# Go commands
build:
	@echo "Building application..."
	go build -o url-shortener .

run:
	@echo "Running application..."
	go run .

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f url-shortener
	go clean

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

swagger:
	@echo "Generating Swagger documentation..."
	swag init

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t url-shortener:latest .

docker-run:
	@echo "Starting services with Docker Compose (development)..."
	docker-compose up -d

docker-prod:
	@echo "Starting services with Docker Compose (production)..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

logs:
	@echo "Viewing application logs..."
	docker-compose logs -f app

# Database commands
db-migrate:
	@echo "Running database migrations..."
	docker-compose exec postgres psql -U postgres -d urlshortener -f /docker-entrypoint-initdb.d/init.sql

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

dev-start:
	@echo "Starting development environment..."
	docker-compose up postgres redis -d
	@echo "Database and Redis started. Run 'make run' to start the application."

# Production helpers
prod-deploy:
	@echo "Deploying to production..."
	@if [ ! -f .env ]; then echo "Error: .env file not found. Copy .env.example and configure it."; exit 1; fi
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

health-check:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "Health check failed"

# Utility commands
format:
	@echo "Formatting Go code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run

security-scan:
	@echo "Running security scan..."
	gosec ./...
