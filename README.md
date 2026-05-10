# Synapse

Synapse is a multi-tenant learning intelligence platform for university students. It connects to Canvas LMS, models per-topic mastery, powers spaced-repetition review, plans study sessions, and supports school-partitioned peer learning through feed and tutoring workflows.

## What Is In This Repo

This is a monorepo with an Angular frontend, Go backend services, local Docker infrastructure, SQL migrations, and seed data.

- `apps/web` - Angular app for the landing page and authenticated student dashboard.
- `services/api-gateway` - Go REST API for auth, profile, learning, review, planner, feed, tutoring, and AI-assisted flashcard generation.
- `services/lms-service` - Go service for Canvas OAuth, LMS connection status, course sync, and LMS data import.
- `services/mock-canvas` - Local Canvas-compatible OAuth/API mock for end-to-end development.
- `infra/migrations` - Postgres migrations managed with `golang-migrate`.
- `infra/db` and `infra/seed` - Local bootstrap and demo data.
- `docs` - Architecture notes, ADRs, and mock Canvas documentation.

## Tech Stack

- Frontend: Angular, standalone components, Angular Signals, Tailwind utilities, Lucide icons.
- Backend: Go, Chi, pgx, Redis, Zap, JWT auth.
- Data: Postgres 16 with pgvector, Redis 7.
- Local orchestration: Docker Compose and Make.

## Prerequisites

- Docker and Docker Compose
- Go 1.23+
- Node.js and npm
- `make`

## Local Setup

Create a local environment file:

```sh
cp .env.example .env
```

For local Docker development, keep `DATABASE_URL` pointed at the Compose service hostname:

```env
DATABASE_URL=postgres://synapse:synapse@postgres:5432/synapse?sslmode=disable
REDIS_URL=redis://redis:6379
JWT_SECRET=dev-secret-change-in-prod
ENCRYPTION_KEY=<base64-encoded-32-byte-key>
```

Generate an LMS encryption key with:

```sh
openssl rand -base64 32
```

AI provider keys are optional unless you are testing RAG flashcard generation.

## Quickstart

Start the local infrastructure and backend services:

```sh
make up
```

Apply database migrations:

```sh
make migrate-up
```

Start the Angular dev server:

```sh
make web
```

Open the app at `http://localhost:4200`.

Service URLs:

- Web app: `http://localhost:4200`
- API gateway: `http://localhost:8080`
- LMS service: `http://localhost:8081`
- Mock Canvas: `http://localhost:8082`
- Postgres: `localhost:5432`
- Redis: `localhost:6379`

## Common Commands

```sh
make help          # show available commands
make up            # start Docker services
make down          # stop Docker services
make rebuild       # rebuild containers and restart
make logs          # tail Docker logs
make ps            # list Docker services
make migrate-up    # apply migrations
make migrate-down  # roll back the latest migration
make seed          # load local demo data
make bhatta-seed   # load Fisk demo data for abhatta01@my.fisk.edu
make web           # install frontend deps and run Angular dev server
make test          # run API gateway tests
```

## Running Services Manually

The recommended local flow is Docker Compose plus `make web`, but services can also be run from their directories during focused development:

```sh
cd services/api-gateway
go test ./...
```

```sh
cd apps/web
npm install
npm start
```

The Angular app proxies `/v1` requests to `http://localhost:8080` in development.

## API Overview

Health checks:

- `GET /healthz`
- `GET /readyz`

Public API gateway routes:

- `POST /v1/auth/register`
- `POST /v1/auth/login`
- `POST /v1/auth/logout`

Authenticated API gateway areas:

- `GET /v1/me`
- `GET /v1/profile/summary`
- `/v1/courses`
- `/v1/notes`
- `/v1/flashcards`
- `/v1/review`
- `/v1/insights`
- `/v1/planner`
- `/v1/tutoring`
- `/v1/feed`

LMS service routes:

- `GET /v1/lms/connect/canvas`
- `GET /v1/lms/callback/canvas`
- `POST /v1/lms/connect/token`
- `GET /v1/lms/status`
- `GET /v1/lms/courses`
- `POST /v1/lms/sync`
- `DELETE /v1/lms/disconnect`

Authenticated routes require an `Authorization: Bearer <token>` header. Tenant isolation is based on the `school_id` JWT claim.

## Mock Canvas

`services/mock-canvas` provides a local Canvas-like OAuth flow and API for development. It runs at `http://localhost:8082` and supports:

- OAuth authorize, approve, and token exchange endpoints.
- Canvas-like course, assignment, submission, and user endpoints.
- A static dev PAT token: `mock-canvas-pat-dev`.

See `docs/mock-canvas.md` for the full local OAuth contract.

## Database

Migrations live in `infra/migrations` and are applied with:

```sh
make migrate-up
```

The schema includes schools, users, Canvas/LMS connection data, learning objects, RAG support, planner data, feed data, and tutoring data. User-owned tables are scoped by `school_id`.

## Development Notes

- Backend services follow a clean architecture layout under `internal/domain`, `internal/service`, `internal/repository`, and `internal/transport`.
- Go database access uses raw pgx queries.
- The frontend uses standalone Angular components, Signals, reactive forms, and guarded authenticated routes.
- Root `.env` values are consumed by Docker Compose.
- Do not run migrations from application startup; use `make migrate-up`.

## Documentation

- `docs/architecture/phase0_foundations.md` - foundation architecture notes.
- `docs/adr/ADR-0001-stack.md` - stack decisions.
- `docs/adr/ADR-0002-multitenancy-schoolid.md` - school tenancy model.
- `docs/mock-canvas.md` - mock Canvas API and OAuth behavior.
