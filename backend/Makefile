build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

run:
	@go run cmd/api/main.go

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

clean:
	@echo "Cleaning..."
	@rm -f main.exe

.PHONY: build run dev docker-run docker-down clean
