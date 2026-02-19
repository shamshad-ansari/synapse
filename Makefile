.PHONY: help up down logs ps migrate-up migrate-down fmt lint test

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

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f --tail=200

ps:
	docker compose ps

# NOTE: These will be wired in Phase 0 Step 2 once the DB container and migrate tool are set up.
migrate-up:
	@echo "TODO: wire migrations in Phase 0 Step 2"

migrate-down:
	@echo "TODO: wire migrations in Phase 0 Step 2"

fmt:
	@echo "TODO: add gofmt + frontend formatting in Phase 0 Step 3/4"

lint:
	@echo "TODO: add golangci-lint + frontend lint in Phase 0 Step 3/4"

test:
	@echo "TODO: add go test + frontend tests in Phase 0 Step 3/4"
