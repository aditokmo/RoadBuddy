# RoadBuddy Backend

Go backend API for the RoadBuddy ride-sharing platform.

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Air (automatic install on first `make dev`)

## Quick Start

### 1. Clone & Setup

```bash
git clone https://github.com/aditokmo/RoadBuddy.git
cd backend
```

If you don’t have `make`, first install it, then use this command.

```bash
make setup
```

The `make setup` command will:
- Copy `.env.example` to `.env`
- Download Go dependencies
- Start the PostgreSQL database in Docker
- Run database migrations

### 2. Start Development

```bash
make dev
```

This starts the backend with live reload using Air. The database will be running in the background.

## Stopping Services

```bash
make docker-down
```

## Available Commands

- `make setup` - First-time setup (dependencies, database, migrations)
- `make dev` - Start development server with hot reload
- `make docker-run` - Run all services in Docker
- `make docker-down` - Stop all Docker services
- `make docker-migrate-up` - Run pending migrations
- `make docker-migrate-down` - Rollback last migration
- `make db-reset` - Reset database
- `make clean` - Clean build
