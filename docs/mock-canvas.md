# Mock Canvas API Contract

This document defines how `services/mock-canvas` behaves in local development.
It is intentionally close to Canvas OAuth/API behavior, while still simple for
repeatable local testing.

## Base URL

- Browser-facing URL: `http://localhost:8082`
- Docker-internal URL: `http://mock-canvas:8082`

Use `localhost` in frontend/browser flows and `mock-canvas` for service-to-service
calls inside Docker.

## OAuth 2.0 (Authorization Code)

### 1) Authorization endpoint

- `GET /login/oauth2/auth`
- Required query params:
  - `client_id`
  - `redirect_uri`
  - `response_type=code`
  - `state` (recommended)

Behavior:
- Returns a consent page.
- If required params are missing/invalid, returns `400`.

### 2) User approval endpoint

- `GET /login/oauth2/approve`
- Required query params:
  - `client_id`
  - `redirect_uri`
  - `state` (optional but expected)

Behavior:
- Generates a one-time authorization code (TTL: 5 minutes).
- Redirects to:
  - `{redirect_uri}?code=<code>&state=<state>`

### 3) Token exchange endpoint

- `POST /login/oauth2/token`
- Content-Type: `application/x-www-form-urlencoded`
- Required form fields:
  - `grant_type=authorization_code`
  - `client_id`
  - `client_secret`
  - `redirect_uri`
  - `code`

Behavior:
- Validates code exists, is unexpired, and matches `client_id` + `redirect_uri`.
- Authorization code is single-use.
- On success (`200`):
  - `access_token`
  - `refresh_token`
  - `token_type=Bearer`
  - `expires_in=36000`
- On failure (`400`), returns OAuth-style errors:
  - `invalid_request`
  - `unsupported_grant_type`
  - `invalid_grant`

## Access Tokens (Mock behavior)

Protected `GET /api/v1/*` endpoints require:
- `Authorization: Bearer <token>`

Accepted tokens:
- OAuth access tokens minted by `/login/oauth2/token`
- Static dev PAT token: `mock-canvas-pat-dev`

Rejected tokens:
- Any arbitrary/non-issued token (returns `401`)

## Implemented Canvas-like endpoints

All endpoints below require a valid Bearer token:

- `GET /api/v1/users/self`
- `GET /api/v1/courses`
- `GET /api/v1/courses/{courseID}`
- `GET /api/v1/courses/{courseID}/assignments`
- `GET /api/v1/courses/{courseID}/students/submissions`

Notes:
- `GET /api/v1/courses/{courseID}?include[]=syllabus_body` returns `syllabus_body`.
- `GET /api/v1/courses` omits `syllabus_body` by default.
- Error shape for token failures is Canvas-like:
  - `{"errors":[{"message":"Invalid access token."}],"status":"unauthenticated"}`

## Local end-to-end checklist

1. Start services:
   - `docker compose up -d api-gateway lms-service mock-canvas`
2. Login/register in web app.
3. Open Canvas connect page and use:
   - `http://localhost:8082`
4. Authorize on mock consent screen.
5. Confirm redirect to:
   - `/canvas/connected?status=success`
6. Trigger sync and verify data appears via LMS endpoints.

## Why this differs from production Canvas

- No tenant/account-level app registration checks.
- Consent UI is simplified.
- Token lifecycle is in-memory and resets on container restart.
- Scope handling is not enforced.

These differences are intentional for local velocity, while preserving realistic
OAuth semantics (required params, one-time codes, code matching, and token validity).
