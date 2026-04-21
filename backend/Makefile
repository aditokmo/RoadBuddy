-include .env

setup:
	@echo "Setting up RoadBuddy backend..."
	@echo "Creating .env from .env.example..."
	@cp .env.example .env
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Starting database..."
	@docker compose up db -d --wait
	@echo "Running migrations..."
	@docker compose --profile tools run --rm migrate
	@echo "Setup complete! Run 'make dev' to start development"

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

run:
	@go run cmd/api/main.go

stop:
	-@taskkill /f /im main.exe 2>nul || exit 0
	-@powershell -Command "if (Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue) { Stop-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess -Force }"

# Start DB in Docker, then run backend with Air (dev workflow)
dev:
	@echo "Starting database..."
	@docker compose up db -d --wait
	@echo "Starting backend with live reload..."
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { air } else { Write-Output 'Installing air...'; go install github.com/air-verse/air@latest; air }"

# Full stack in Docker
docker-run:
	@docker compose up --build

docker-down:
	@docker compose down

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

migrate-drop:
	@echo "Dropping all migrations..."
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" drop -f

migrate-version:
	@echo "Current migration version..."
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" version

migrate-create:
	@powershell -Command "$$name = Read-Host 'Enter migration name'; migrate create -ext sql -dir ./migrations -seq $$name"

# Docker-based migration commands
docker-migrate-up:
	@echo "Running migrations in Docker..."
	@docker compose --profile tools run --rm migrate

docker-migrate-down:
	@echo "Rolling back migrations in Docker..."
	@docker compose --profile tools run --rm migrate -path=/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}?sslmode=disable" down

# Database setup/teardown
db-setup: docker-migrate-up
	@echo "Database setup complete"

db-reset: docker-migrate-down docker-migrate-up
	@echo "Database reset complete"

# Complete reset - like a fresh clone
reset:
	@echo "Performing complete reset..."
	@docker compose down -v
	@rm -f .env

clean:
	@echo "Cleaning..."
	@rm -f main.exe

.PHONY: build run dev docker-run docker-down migrate-up migrate-down migrate-drop migrate-version migrate-create docker-migrate-up docker-migrate-down db-setup db-reset reset clean
