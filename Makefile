.PHONY: help up down logs ps migrate-up migrate-down rebuild fmt lint test check-db-url web seed

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
	@echo "  make rebuild      - rebuild containers and restart (use after Go changes)"
	@echo "  make logs         - tail docker logs"
	@echo "  make ps           - show docker containers"
	@echo ""
	@echo "DB:"
	@echo "  make migrate-up   - apply db migrations"
	@echo "  make migrate-down - rollback the last db migration"
	@echo "  make seed         - load local SQL (Alex + cohort feed/tutoring; needs postgres up)"
	@echo ""
	@echo "Frontend:"
	@echo "  make web          - run the Angular dev server (localhost:4200)"
	@echo ""
	@echo "Quality:"
	@echo "  make fmt          - format code"
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

# Full local data: Alex prereqs, learning seed, cohort feed/tutoring, planner seed.
# Requires: docker compose up -d postgres (and migrate-up). Uses psql inside the postgres container.
seed:
	@docker compose ps postgres --status running --quiet 2>/dev/null | grep -q . || (echo "ERROR: postgres not running. Run: make up"; exit 1)
	cat infra/db/bootstrap_alex_prereqs.sql infra/db/populate_alex_data.sql infra/db/populate_alex_metrics.sql infra/db/demo_cohort_feed_tutoring.sql infra/seed/dev_seed.sql | docker compose exec -T postgres psql -U synapse -d synapse -v ON_ERROR_STOP=1

web:
	@echo "Starting Angular dev server..."
	cd apps/web && npm install && npm start

fmt:
	@echo "TODO: add gofmt + frontend formatting"

lint:
	@echo "TODO: add golangci-lint + frontend lint"

test:
	cd services/api-gateway && go test ./...

# Requires running API + migrations; set TOKEN or LOGIN_EMAIL, LOGIN_PASSWORD, SCHOOL_DOMAIN
test-feed-api:
	cd services/api-gateway && ./scripts/test_feed_api.sh
