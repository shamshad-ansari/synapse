# Synapse — Network-Aware Learning Intelligence Platform

Synapse is a multi-tenant learning intelligence system that models individual mastery and collective confusion to optimize what students should study and who they should learn from.

## Monorepo layout

- `apps/web` — Angular frontend
- `services/api-gateway` — Go REST gateway (auth, tenancy context, routing)
- `infra` — Postgres/Redis, migrations, local dev resources
- `docs` — architecture notes + ADRs
- `scripts` — local dev scripts

## Phase 0 (Foundations)

Step 1: Repo bootstrap (this step)  
Step 2: Infra wiring (docker-compose, Postgres, Redis, migrations)  
Step 3: Go API gateway skeleton (`/healthz`, `/readyz`, `/v1/me`)  
Step 4: Angular app shell (routing, auth scaffold, guarded routes)

## Quickstart (will work after Phase 0 Step 2)
- `make up` — start infra + services
- `make migrate-up` — apply DB migrations
- `make down` — stop everything
