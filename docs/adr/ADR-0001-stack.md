# ADR-0001: Tech Stack

## Decision
- Frontend: Angular (TypeScript, standalone components)
- Backend: Go
- DB: PostgreSQL (+ pgvector)
- Cache: Redis
- Async jobs: start simple (in-service worker/queue), scale later

## Rationale
- Strong typing end-to-end
- Supports multi-tenant partitioning by `school_id`
