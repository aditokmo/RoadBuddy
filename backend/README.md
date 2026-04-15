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
cp .env.example .env
```

### 2. Configure Environment

Edit `.env` and set (or use defaults for local dev):

### 3. Start Development

```bash
make docker-run
```

## Stopping Services

```bash
make docker-down
```
