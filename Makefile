.PHONY: help up down logs ps migrate-up migrate-down fmt lint test check-db-url

.DEFAULT_GOAL := help

# Optionally load .env if present (won't error if missing)
-include .env
export

DB_URL ?= $(DATABASE_URL)

help:
	@echo "Synapse Dev Commands"
	@echo ""
	@echo "Infra:"
	@echo "  make up           - start local infra (docker compose)"
	@echo "  make down         - stop local infra"
	@echo "  make logs         - tail docker logs"
	@echo "  make ps           - show docker containers"
	@echo ""
	@echo "DB:"
	@echo "  make migrate-up   - apply db migrations"
	@echo "  make migrate-down - rollback db migrations"
	@echo ""
	@echo "Quality:"
	@echo "  make fmt          - format code (go fmt + frontend fmt)"
	@echo "  make lint         - lint code"
	@echo "  make test         - run tests"

rebuild:
	docker compose up -d --build

up:
	docker compose up -d 

down:
	docker compose down

logs:
	docker compose logs -f --tail=200

ps:
	docker compose ps

check-db-url:
	@if [ -z "$(DB_URL)" ]; then \
		echo "ERROR: DB_URL is empty. Set DATABASE_URL (or create .env with DATABASE_URL=...)"; \
		exit 1; \
	fi

migrate-up: check-db-url
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" up

migrate-down: check-db-url
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" down 1

fmt:
	@echo "TODO: add gofmt + frontend formatting in Phase 0 Step 3/4"

lint:
	@echo "TODO: add golangci-lint + frontend lint in Phase 0 Step 3/4"

test:
	@echo "TODO: add go test + frontend tests in Phase 0 Step 3/4"