.PHONY: run build test test-unit test-integration clean docker-up docker-down migrate-up migrate-down

# Run application
run:
	go run cmd/api/main.go

# Build binary
build:
	go build -o bin/api cmd/api/main.go

# Run all tests
test:
	go test -v -cover ./...

# Run unit tests only
test-unit:
	go test -v -cover -short ./internal/features/...

# Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific test
test-specific:
	go test -v -run $(TEST) ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Start all services
docker-up:
	docker-compose up -d

# Stop all services
docker-down:
	docker-compose down -v

# View docker logs
docker-logs:
	docker-compose logs -f

swagger-install:
	@echo "📦 Installing swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Swag installed"

# Generate swagger documentation
swagger:
	@echo "📝 Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "✅ Swagger docs generated at: docs/"
	@echo "📖 View at: http://localhost:8080/swagger/index.html"

# Format swagger comments
swagger-fmt:
	@echo "🎨 Formatting Swagger comments..."
	swag fmt
	@echo "✅ Swagger comments formatted"

# Database migration up
migrate-up:
	go run cmd/migrate/main.go -action=up

# Database migration down
migrate-down:
	go run cmd/migrate/main.go -action=down

# Check migration version
migrate-version:
	go run cmd/migrate/main.go -action=version

# Force migration to specific version
migrate-force:
	go run cmd/migrate/main.go -action=force -version=$(VERSION)

# Tidy dependencies
tidy:
	go mod tidy

# Install dependencies
install:
	go mod download

# Generate mocks
mock:
	@echo "Generating mocks..."
	go generate ./...

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Start development environment
dev: docker-up
	@echo "Waiting for databases to be ready..."
	@sleep 5
	@make migrate-up
	@echo "Development environment ready!"
	@make run

# Run tests in Docker
test-docker:
	docker-compose -f docker-compose.test.yaml up --abort-on-container-exit --exit-code-from app

.DEFAULT_GOAL := run