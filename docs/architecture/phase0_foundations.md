# Phase 0 — Foundations

Goal: establish a production-style monorepo with local infra, migrations, a Go API gateway skeleton, and an Angular app shell.

## Steps
1. Repo bootstrap (structure + basic tooling)
2. Local infra (Postgres + pgvector, Redis, migrations)
3. Go API gateway skeleton (health + auth scaffold + tenancy context)
4. Angular shell (routing + auth scaffold + guarded routes)

## Phase 0 Definition of Done
- `make up` runs infra and services
- migrations apply cleanly
- `/healthz` and `/readyz` return OK
- `/v1/me` returns current user with `school_id`
- Angular app loads and calls `/v1/me`