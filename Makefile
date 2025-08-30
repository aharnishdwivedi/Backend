.PHONY: build test test-coverage run clean migrate-up migrate-down

# Build the application
build:
	go build -o bin/incident-triage-assistant cmd/main.go

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run the application
run:
	go run cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations
migrate-up:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/incident_triage" up

migrate-down:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/incident_triage" down

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Setup development environment
setup: deps migrate-up

# Development workflow
dev: setup run

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Docker development workflow
docker-dev: docker-build docker-up
