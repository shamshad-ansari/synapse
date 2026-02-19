# ADR-0002: Multi-tenancy via school_id

## Decision
All user-owned and social data is scoped by `school_id`.

## Rules
- Default: all queries filter by `school_id`
- Global content is explicit opt-in and stored with scope=GLOBAL
- Backend enforces isolation (never rely on frontend)

## Rationale
- Prevent cross-tenant leakage
- Matches product model (school communities + global opt-in)
- Demonstrates privacy-aware system design